package auth

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	infAuthLoginRequest = "new auth request: '%s'"
)

type Auth struct {
	Options
}

func New(opts ...Option) *Auth {
	var (
		options = newOptions(opts...)
		auth    = Auth{options}
	)
	return &auth
}

func (a *Auth) Middleware() echo.MiddlewareFunc {
	return middleware.BasicAuth(a.check)
}

func (a *Auth) check(username, password string, ctx echo.Context) (result bool, err error) {
	var user User
	a.Logger.Infof(infAuthLoginRequest, username)
	if err = user.ByName(a.DB, username); err != nil {
		return
	}
	result, err = user.CheckPasswordHash(password)
	return
}
