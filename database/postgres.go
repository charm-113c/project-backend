package database

import (
	"context"
	"fmt"

	"github.com/charm-113c/project-zero/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// startPostgres establishes a connection pool with the DB designated through the
// config, and then populates the Conns field of the Storage struct, essentially
// implementing the Storage.Conns StorageHandler interfaces
func startPostgres(ctx context.Context, cfg config.Config, stg *Storage) error {
	// Construct DB URL for connection
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)
	poolCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("error parsing connection string: %w", err)
	}

	if cfg.Database.ConnPoolSize == 0 {
		stg.logger.Warn("DB connection max pool size not set, setting max to default value of 20")
		poolCfg.MaxConns = int32(20)
	} else {
		poolCfg.MaxConns = int32(cfg.Database.ConnPoolSize)
	}

	// Create connection pool
	connPool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return fmt.Errorf("error creating DB connection pool: %w", err)
	}

	stg.logger.Info("DB connection pool created")
	if err = connPool.Ping(ctx); err != nil {
		return fmt.Errorf("connection to DB not established, ping query to DB failed: %w", err)
	}
	stg.logger.Info("Connection to DB established")

	// Finally, assign the different handlers to Storage
	stg.Conns.AccTableOps = &pgAccountHandler{pool: connPool}
	stg.Conns.EvTableOps = &pgEventHandler{pool: connPool}
	stg.Conns.SocialTableOps = &pgSocialHandler{pool: connPool}
	stg.Conns.MapTableOps = &pgMapHandler{pool: connPool}

	return nil
}

// pgAccountHandler populates the Storage.Conns.AccountStorageHandler field,
// and its methods implement the AccountStorageHandler interface
type pgAccountHandler struct {
	pool *pgxpool.Pool
}

// pgEventHandler populates the Storage.Conns.EventStorageHandler field,
// and its methods implement the EventStorageHandler interface
type pgEventHandler struct {
	pool *pgxpool.Pool
}

// pgSocialHandler populates the Storage.Conns.SocialStoragesHandler field,
// and its methods implement the SocialStoragesHandler interface
type pgSocialHandler struct {
	pool *pgxpool.Pool
}

type pgMapHandler struct {
	pool *pgxpool.Pool
}

func (usrTable *pgAccountHandler) CreateAccount() {}
func (evTable *pgEventHandler) CreateEvent()      {}
func (socTable *pgSocialHandler) FollowUser()     {}
func (socTable *pgMapHandler) GetMap()            {}
