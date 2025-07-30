/* Package handlers defines the request subhandlers as structs
* that implement the handler interfaces
 */
package handlers

import (
	"github.com/charm-113c/project-zero/database"
	"go.uber.org/zap"
)

// AccountHandler implements the AccountRequests interface and is a subhandler
// for all requests relating to accounts
type AccountHandler struct {
	DB     database.AccountStorageHandler
	Logger *zap.Logger
}

// EventHandler implements the EventRequests interface and handles
// requests relating to events
type EventHandler struct {
	DB     database.EventStorageHandler
	Logger *zap.Logger
}

// SocialHandler implements the SocialRequests interface and handles
// requests relating to social things
type SocialHandler struct {
	DB     database.SocialStorageHandler
	Logger *zap.Logger
}

// MapHandler implements the MapRequests interface and handles
// all relevant requests
type MapHandler struct {
	DB     database.MapStorageHandler
	Logger *zap.Logger
}

// NewAccountHandler instantiates an AccountHandler
func NewAccountHandler(db database.AccountStorageHandler, logger *zap.Logger) *AccountHandler {
	return &AccountHandler{
		db,
		logger,
	}
}

// NewEventHandler instantiates an EventHandler
func NewEventHandler(db database.EventStorageHandler, logger *zap.Logger) *EventHandler {
	return &EventHandler{
		db,
		logger,
	}
}

// NewSocialHandler instantiates an SocialHandler
func NewSocialHandler(db database.SocialStorageHandler, logger *zap.Logger) *SocialHandler {
	return &SocialHandler{
		db,
		logger,
	}
}

// NewMapHandler instantiates an MapHandler
func NewMapHandler(db database.MapStorageHandler, logger *zap.Logger) *MapHandler {
	return &MapHandler{
		db,
		logger,
	}
}
