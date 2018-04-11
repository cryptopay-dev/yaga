package auth

import (
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/cryptopay-dev/yaga/web"
	"github.com/go-pg/pg"
)

const (
	infAuthLoginRequest = "new auth request: '%s'"
)

// Auth struct
type Auth struct {
	DB *pg.DB
}

// New creates new Auth
func New(db *pg.DB) *Auth {
	return &Auth{DB: db}
}

// Middleware for web-application
func (a *Auth) Middleware() web.MiddlewareFunc {
	return web.BasicAuth(a.check)
}

// check username and password from request
func (a *Auth) check(username, password string, ctx web.Context) (result bool, err error) {
	var user User
	log.Infof(infAuthLoginRequest, username)
	if err = user.ByName(a.DB, username); err != nil {
		return
	}
	result, err = user.CheckPasswordHash(password)
	return
}
