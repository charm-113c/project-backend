package api

import (
	// echojwt "github.com/labstack/echo-jwt/v4"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/logto-io/go/v2/client"
	"go.uber.org/zap"
)

func setUpRoutes(e *echo.Echo, rh *RequestHandler, logtoCfg *client.LogtoConfig, logger *zap.Logger) error {
	useLoggerMiddleware(e, logger)

	// TODO: Have a look at Echo's middleware arsenal, put in what is necessary

	// Set up routes for log-in
	// NOTE: the following is for development purposes, modify for prod!

	// Create simple home page
	e.GET("/", func(c echo.Context) error {
		client := client.NewLogtoClient(
			logtoCfg,
			&echoSessionStorage{c},
		)

		sess, err := session.Get("session", c)
		if err != nil {
			c.Logger().Printf("There was an error creating the session: %v", err)
			return err
		}
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   false, // NOTE: set to !DevMode
		}
		sess.Values["logtoEndpoint"] = logtoCfg.Endpoint
		sess.Values["logtoAppID"] = logtoCfg.AppId
		sess.Values["logtoAppSecret"] = logtoCfg.AppSecret
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			c.Logger().Printf("There was an error saving the logto config: %v", err)
			return err
		}

		authState := "You are not logged in :("

		if client.IsAuthenticated() {
			authState = "You're logged in! :D"
		}

		return c.HTML(http.StatusOK, "<h1>Hello with Logto</h1>"+"<div>"+authState+"</div>")
	})
	// Create sign-in route
	e.GET("account/login", rh.AccountReqs.LoginUser)

	return nil
}

// useLoggerMiddleware applies Echo's routing logger to our current logger,
// giving it access to specific request fields
func useLoggerMiddleware(e *echo.Echo, logger *zap.Logger) {
	// Echo middleware allows access to detailed request data
	// See their docs for more options
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Sugar().Info(zap.String("request", "successful"),
					zap.Int("Status", v.Status),
					zap.String("URI", v.URI))
			} else {
				logger.Sugar().Info(zap.String("request", "error"),
					zap.Int("Status", v.Status),
					zap.String("URI", v.URI))
				zap.Error(v.Error)
			}
			return nil
		},
	}))
}
