// Package api defines all HTTP request routers and handler functions.
// It creates an echo router that implements the RequestHandler struct
package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"syscall"

	"github.com/gorilla/sessions"
	"github.com/logto-io/go/v2/client"

	"github.com/charm-113c/project-zero/api/handlers"
	"github.com/charm-113c/project-zero/config"
	"github.com/charm-113c/project-zero/database"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// RequestHandler holds interfaces, and these interfaces define the methods
// that the Router will need to handle
type RequestHandler struct {
	AccountReqs AccountRequests
	SocialReqs  SocialRequests
	EventReqs   EventRequests
	MapReqs     MapRequests
}

// NewRequestHandler instantiates a RequestHandler
func NewRequestHandler(db database.Storage, logger *zap.Logger) *RequestHandler {
	return &RequestHandler{
		handlers.NewAccountHandler(db.Conns.AccTableOps, logger),
		handlers.NewSocialHandler(db.Conns.SocialTableOps, logger),
		handlers.NewEventHandler(db.Conns.EvTableOps, logger),
		handlers.NewMapHandler(db.Conns.MapTableOps, logger),
	}
}

// InitRouter is the function responsible for instantiating a Router based on the configuration.
func InitRouter(ctx context.Context, db database.Storage, cfg *config.Config, logger *zap.Logger) (*echo.Echo, error) {
	// Create Echo router that will handle the requests
	e := echo.New()

	rh := NewRequestHandler(db, logger)

	// Change ulimit to maximize number of connections
	// TODO: find optimal number of file descriiptor for given hardware
	logger.Info(`Setting max n° of file descriptors to hardware limit / 2`)
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return e, fmt.Errorf("error getting number of open files: %w", err)
	}
	logger.Sugar().Infof("Setting max number of open files to %d", rLimit.Max/2)
	rLimit.Cur = rLimit.Max / 2

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return e, fmt.Errorf("error changing max number of open files: %w", err)
	}

	logger.Info("Max number of open files has been updated")

	if err = initSessionStore(e, cfg.Server.DevMode); err != nil {
		return e, err
	}

	logtoCfg := initLogtoCfg(cfg)

	if err := setUpRoutes(e, rh, logtoCfg, logger); err != nil {
		err = fmt.Errorf("router failed to set up routes: %v", err)
		return e, err
	}

	return e, nil
}

// initSessionStore primes the router with a store middleware,
// which will be needed to store user sessions
func initSessionStore(e *echo.Echo, devMode bool) error {
	if devMode {
		store := sessions.NewFilesystemStore("/", []byte("1234567890"))
		store.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400, // 1 day
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}

		e.Use(session.Middleware(store))
		return nil
	}
	// TODO: else use database package's cache
	return errors.New("trying to initialise session store in prod mode but Cache hasn't been implemented yet")
}

func initLogtoCfg(cfg *config.Config) *client.LogtoConfig {
	// if !cfg.Server.DevMode {
	// TODO: in this case, create the config here
	// and save it in cache!
	// We need the config at every request in order
	// to use Logto's services
	// }
	return &client.LogtoConfig{
		Endpoint:  cfg.Logto.Endpoint,
		AppId:     cfg.Logto.AppID,
		AppSecret: cfg.Logto.AppSecret,
	}
}
