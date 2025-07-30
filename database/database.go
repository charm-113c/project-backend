package database

import (
	"context"
	"fmt"
	"time"

	"github.com/charm-113c/project-zero/config"
	"go.uber.org/zap"
)

// Storage struct is the database component; its Conns field represents the connection pool
// it establishes with the actual DB: Conns must therefore implement all the necessary DB operations
// and those operations are defined in the StorageHandler interfaces.
type Storage struct {
	Conns struct {
		Close          GracefulShutdown
		AccTableOps    AccountStorageHandler
		EvTableOps     EventStorageHandler
		SocialTableOps SocialStorageHandler
		MapTableOps    MapStorageHandler
	}
	Cache  KeyValCache
	logger *zap.Logger
}

// TODO: create Redis cache

// KeyValCache caches the most read fields from the database. So doing it allows
// leveraging the high performance of Golang routers (by reducing interactions with
// the DB, often the bottleneck of a system)
type KeyValCache interface {
	Set(key string, value any, date time.Duration, ttl time.Duration)
	Get(key string) (any, bool)
	Invalidate(key string) error // Invalidate a cache entry
	Close() error
}

// StartStorage initializes and connects to the DB and also instantiates a DB-specific logger.
// It is designed with flexibility in mind, and abstracts away from the implementation of teh DB
func StartStorage(ctx context.Context, cfg config.Config, storage *Storage, parentLogger *zap.Logger) error {
	// var storage Storage

	// Launch DB logger
	dbLogger := parentLogger.With(zap.String("component", "database"))
	dbLogger.Info("Database logger initialized, creating DB connection pool")
	storage.logger = dbLogger

	// NOTE: Caching is now implemented as a decorator over the *StorageHandler interfaces.
	// This means the storage handlers are responsible for searching some data in cache before checking the DB

	// Start the DB
	dbLogger.Info("Database type: " + cfg.Database.Type)
	switch cfg.Database.Type {
	case "postgres", "sql":
		if err := startPostgres(ctx, cfg, storage); err != nil {
			return fmt.Errorf("could not start the DB: %w", err)
		}
	default:
		return fmt.Errorf("database of type %s is unsupported", cfg.Database.Type)
	}

	// TODO: implement cache
	dbLogger.Warn("Cache has not yet been implemented")

	return nil
}

// GracefulShutdown is the interface for closing connections with the Storage (e.g. with postgres)
type GracefulShutdown interface {
	CloseConns() error
}

// AccountStorageHandler is responsible for defining the operations on the User table
type AccountStorageHandler interface {
	CreateAccount()
}

// EventStorageHandler is responsible for defining the operations on the Event table
type EventStorageHandler interface {
	CreateEvent()
}

// SocialStoragesHandler is responsible for defining the operations on the tables that
// relate to social interactions between users
type SocialStorageHandler interface {
	FollowUser()
}

// MapStorageHandler is responsible for defining the operations on the tables that
// relate to social interactions between users
type MapStorageHandler interface {
	GetMap()
}
