package api

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// echoSessionStorage implements Logto's Storage interface
type echoSessionStorage struct {
	ctx echo.Context
}

func (s *echoSessionStorage) GetItem(key string) string {
	sess, _ := session.Get("logto-session", s.ctx)
	// TODO: understand why and how this works.
	if val, ok := sess.Values[key]; ok {
		return val.(string)
	}
	return ""
}

func (s *echoSessionStorage) SetItem(key, val string) {
	sess, _ := session.Get("logto-session", s.ctx)
	sess.Values[key] = val
	sess.Save(s.ctx.Request(), s.ctx.Response())
}
