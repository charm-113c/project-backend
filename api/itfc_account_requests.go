package api

import "github.com/labstack/echo/v4"

// AccountRequests contains the methods that need to be implemented by
// Router types to handle requests concerning user profiles and accounts.
// These requests correspond to the User Accounts/Profile operationIDs in the
// docs/operationIDs.md file (in the docs repo).
type AccountRequests interface {
	LoginUser(c echo.Context) error
}
