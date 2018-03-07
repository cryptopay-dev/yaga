package auth

import (
	"github.com/cryptopay-dev/yaga/web"
)

const (
	infAuthLoginRequest = "new auth request: '%s'"
)

// Auth struct
type Auth struct {
	Options
}

// New creates new Auth
func New(opts ...Option) *Auth {
	var (
		options = newOptions(opts...)
		auth    = Auth{options}
	)
	return &auth
}

// Middleware for web-application
func (a *Auth) Middleware() web.MiddlewareFunc {
	return web.BasicAuth(a.check)
}

// check username and password from request
func (a *Auth) check(username, password string, ctx web.Context) (result bool, err error) {
	var user User
	a.Logger.Infof(infAuthLoginRequest, username)
	if err = user.ByName(a.DB, username); err != nil {
		return
	}
	result, err = user.CheckPasswordHash(password)
	return
}
