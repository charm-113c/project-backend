package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/charm-113c/project-zero/api"
	"github.com/charm-113c/project-zero/config"
	"github.com/charm-113c/project-zero/database"
	"github.com/charm-113c/project-zero/util"
	"go.uber.org/zap"
)

// Server is the struct containing vital server data
type Server struct {
	db  database.Storage
	cfg config.Config
}

func main() {
	if err := run(); err != nil {
		log.Println("FATAL: server has run into an error: ", err)
		os.Exit(1)
		// Equivalent to log.Fatal, but more explicit
	}
	log.Println("Server has shutdown properly")
}

// TODO: use (implement?) profiling tools to find out bottlenecks

// run is an auxiliary function that initializes and effectively starts the server
// and connects to all necessary services
func run() error {
	log.Println("Server starting")
	var srv Server

	if err := config.LoadConfig(&srv.cfg); err != nil {
		return fmt.Errorf("couldn't load configuration: %w", err)
	}

	if err := srv.cfg.Validate(); err != nil {
		return fmt.Errorf("the current configuration is not valid: %w", err)
	}

	if srv.cfg.Server.DevMode {
		log.Println("Starting server in development mode")
	} else {
		log.Println("Starting server in production mode")
	}

	// Create context that listens to interrupt signals (e.g. CTRL+C) to pass down
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	// Start logger
	logger, logfile, err := startLogger(&srv.cfg)
	if err != nil {
		log.Printf("FATAL: couldn't start logger: %v", err)
		return err
	}
	defer func() {
		// Dump remaining logger buffer into log file when shutting down
		err := logger.Sync()
		// BUG: Sync fails because of underlying system things (has to do with fsync)
		// This doesn't seem to have any noticeable impact so far
		if err != nil {
			log.Printf("WARNING: logger failed to Sync(): %v", err)
		}
		// Close file
		err = logfile.Close()
		if err != nil {
			log.Printf("WARNING: error while closing log file: %v", err)
		}
	}()

	// The logger is now active and its output will be visible in the log file
	logger.Info("Logger constructed successfully")
	// The sugarred logger allows for more flexibility, check package docs for details

	logger.Info("Server is running with the following configuration:",
		zap.String(".env config file path", srv.cfg.CfgPaths.EnvPath),
		zap.String("YAML config file path", srv.cfg.CfgPaths.CfgPath),
		zap.String("Server Host", srv.cfg.Server.Host),
		zap.Uint16("Server Port", srv.cfg.Server.Port),
		zap.Bool("Server Development Mode", srv.cfg.Server.DevMode),
		zap.String("Log file path", srv.cfg.Server.LogFile),
		zap.String("DB Host", srv.cfg.Database.Host),
		zap.Uint16("DB Port", srv.cfg.Database.Port),
	)

	logger.Info("Initializing storage")
	storage := new(database.Storage)
	// Begin storage initialization operation with retry mechanism
	err = util.RetryOperation(ctx, logger, func() error {
		return database.StartStorage(ctx, srv.cfg, storage, logger)
	}, 1, 10*time.Millisecond)
	if err != nil {
		logger.Error("Failed to initialize storage", zap.String("error", err.Error()))
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	logger.Info("Initializing router")
	echoRouter, err := api.InitRouter(ctx, *storage, &srv.cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize router", zap.String("error", err.Error()))
		return fmt.Errorf("failed to initialize router: %w", err)
	}

	// Buffered channel to listen to server errors
	serverErrors := make(chan error, 1)

	listenAddr := srv.cfg.Server.Host + ":" + strconv.Itoa(int(srv.cfg.Server.Port))
	// Listen to requests
	if srv.cfg.Server.DevMode {
		// Start HTTP server
		go func() {
			logger.Warn("Starting HTTP (non-secure) server")
			// TODO: Now that we're hardcoding Echo, update this part
			// WARN: Check what error is returned, net/http returns a ErrServerClosed which should be handled
			// as a proper shutdown (so it isn't actually an error). For each router, check IN THE IMPLEMENTATION
			// of Router.Listen if similar errors are thrown
			serverErrors <- echoRouter.Start(listenAddr)
			logger.Info("Shutting down server")
		}()
	} else {
		// START HTTPS Server
		go func() {
			logger.Info("Starting HTTPS server")
			serverErrors <- echoRouter.StartAutoTLS(listenAddr) // (listenAddr, srv.cfg.Server.CertFile, srv.cfg.Server.KeyFile)
			// AutoTLS should handle all the TLS for us? Wow, alright
			// WARN: Same check as above
			logger.Info("Shutting down server")
		}()
	}

	// TODO: initialize all other necessary functionalities (e.g. Websocket, msg queue)

	// Give some time to complete shut down
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Listen to shutdown signals
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			logger.Sugar().Error("Fatal server error listening to API requests:", err)
			// Log at console and return
			// What do we need to do if a server error occurs?
			// A graceful program exit is probably the best we can do.
		}
		// Else call graceful shutdown
		if err := echoRouter.Shutdown(shutdownCtx); err != nil {
			logger.Sugar().Error("Error calling server Shutdown(): ", err)
			logger.Sugar().Warn("Could not gracefully shutdown, forcefully closing server")
		}
	case <-ctx.Done():
		logger.Sugar().Info("Received shutdown signal. Initiating graceful shutdown")
		// Shutdown server
		if err := echoRouter.Shutdown(shutdownCtx); err != nil {
			logger.Sugar().Error("Error calling server Shutdown(): ", err)
			logger.Sugar().Warn("Could not gracefully shutdown, forcefully closing server")
			if err = echoRouter.Close(); err != nil {
				logger.Sugar().Error("Error during forecul shutdown: ", err)
			}
		}

		// TODO: Close Websocket conns once they're set up
	}

	// Close storage conns
	// NOTE: Since conns close automatically even without calling the methods,
	// errors when closing them should not interrupt the program's flow
	if err := storage.Conns.Close.CloseConns(); err != nil {
		logger.Sugar().Error("Error closing storage connections:", err)
	}
	if err := storage.Cache.Close(); err != nil {
		logger.Sugar().Error("Error closing cache connections:", err)
	}

	return nil
}

// startLogger creates a logger using the zap package and opens/creates logfiles.
// Running the server in development mode
// creates a zap.Development logger instead of a Production one, allowing the use of DPanic.
func startLogger(cfg *config.Config) (*zap.Logger, *os.File, error) {
	config := zap.NewProductionConfig()
	if cfg.Server.DevMode {
		// Set dev logger and write log file in current directory for easy visibility
		config = zap.NewDevelopmentConfig()
		cfg.Server.LogFile = "dev.log"
	}
	// TODO: implement log rotation to keep log file size in check
	// Additionally, log errors in dedicated files on top of current

	// Open the file to write on, and close it in the calling func
	logFile, err := os.OpenFile(cfg.Server.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		err = fmt.Errorf("error opening log file: %w", err)
		return nil, nil, err
	}

	config.OutputPaths = []string{cfg.Server.LogFile, "stdout"}

	logger, err := config.Build()
	if err != nil {
		err = fmt.Errorf("error building logger: %w", err)
		return nil, nil, err
	}

	return logger, logFile, nil
}
