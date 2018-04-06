package auth

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptopay-dev/yaga/helpers/testdb"
	"github.com/cryptopay-dev/yaga/logger/log"
	"github.com/cryptopay-dev/yaga/model"
	"github.com/cryptopay-dev/yaga/web"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const fakeResponse = "test"

type fake struct{}

func (fake) Compare(hash, password string) (bool, error) {
	log.Infof("%q == %q", password, hash)
	return password == hash, nil
}
func (fake) Hash(password string) (string, error) { return password, nil }

func testContext(e *web.Engine) (web.Context, *http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(echo.GET, "/", nil)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	return c, req, rec
}

func TestAuth_Middleware(t *testing.T) {
	var (
		e            *web.Engine
		h            web.HandlerFunc
		db           *pg.DB
		err          error
		user         *User
		username     = "user1"
		password     = "password1"
		authenticate *Auth
	)

	defaultHasher = fake{}

	e, err = web.New(web.Options{})
	assert.NoError(t, err)

	db = testdb.GetTestDB().DB

	authenticate = New(db)

	h = authenticate.Middleware()(func(c web.Context) error {
		return c.String(http.StatusOK, fakeResponse)
	})

	if _, err = model.Delete(db, user, model.Equal("username", username)); err != nil {
		t.Fatal(err)
	}

	// Create new user:
	if user, err = NewUser(db, username, password); err != nil {
		t.Fatal(err)
	}

	if _, err = model.Create(db, user); err != nil {
		t.Fatal(err)
	}

	t.Run("Check that user exists", func(t *testing.T) {
		t.Parallel()

		assert.NoError(
			t,
			new(User).ByName(db, username),
		)
	})

	t.Run("Check that we can't add new user with same name", func(t *testing.T) {
		t.Parallel()

		var _, userErr = NewUser(db, username, "password2")
		assert.EqualError(t, errors.Cause(userErr), ErrUsernameAlreadyTaken.Error())
	})

	t.Run("Unauthorized", func(t *testing.T) {
		t.Parallel()

		c, _, _ := testContext(e)
		errH := h(c)
		assert.Error(t, errH)
		assert.EqualError(t, errH, web.ErrUnauthorized.Error())
	})

	t.Run("Valid user credentials", func(t *testing.T) {
		t.Parallel()

		c, req, rec := testContext(e)
		authData := "basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
		req.Header.Set(echo.HeaderAuthorization, authData)

		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, fakeResponse, rec.Body.String())
		}
	})

	t.Run("Bad user credentials - password", func(t *testing.T) {
		t.Parallel()

		c, req, _ := testContext(e)
		authData := "basic " + base64.StdEncoding.EncodeToString([]byte(username+":password2"))
		req.Header.Set(echo.HeaderAuthorization, authData)
		assert.Error(t, h(c))
	})

	t.Run("Bad user credentials - login", func(t *testing.T) {
		t.Parallel()

		c, req, _ := testContext(e)
		authData := "basic " + base64.StdEncoding.EncodeToString([]byte("user2:"+password))
		req.Header.Set(echo.HeaderAuthorization, authData)
		assert.EqualError(t, errors.Cause(h(c)), pg.ErrNoRows.Error())
	})

}
