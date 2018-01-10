package auth

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cryptopay-dev/yaga/logger/nop"
	"github.com/cryptopay-dev/yaga/testdb"
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

const fakeResponse = "test"

func TestAuth_Middleware(t *testing.T) {
	e := echo.New()
	e.HideBanner = true

	d := testdb.GetTestDB()

	authenticate := New(
		Logger(nop.New()),
		DB(d.DB),
	)

	e.Use(authenticate.Middleware())

	req := httptest.NewRequest(echo.POST, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	h := authenticate.Middleware()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		var err = h(c)
		assert.Error(t, err)
		assert.EqualError(t, err, echo.ErrUnauthorized.Error())
	})

	var (
		err      error
		user     *User
		username = "user1"
		password = "password1"
	)

	// Remove user if it exists:
	if _, err = d.DB.Model(user).Where("username = ?", username).Delete(); err != nil {
		t.Fatal(err)
	}

	// Create new user:
	if user, err = NewUser(d.DB, username, password); err != nil {
		t.Fatal(err)
	}

	if err = d.DB.Insert(user); err != nil {
		t.Fatal(err)
	}

	t.Run("Check that user exists", func(t *testing.T) {
		assert.NoError(
			t,
			new(User).ByName(d.DB, username),
		)
	})

	t.Run("Check that we can't add new user with same name", func(t *testing.T) {
		var _, userErr = NewUser(d.DB, username, "password2")
		assert.EqualError(t, userErr, errUsernameAlreadyTaken.Error())
	})

	t.Run("Valid user credentials", func(t *testing.T) {
		authData := "basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
		req.Header.Set(echo.HeaderAuthorization, authData)

		if assert.NoError(t, h(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, fakeResponse, rec.Body.String())
		}
	})

	t.Run("Bad user credentials - password", func(t *testing.T) {
		authData := "basic " + base64.StdEncoding.EncodeToString([]byte(username+":password2"))
		req.Header.Set(echo.HeaderAuthorization, authData)
		assert.EqualError(t, h(c), bcrypt.ErrMismatchedHashAndPassword.Error())
	})

	t.Run("Bad user credentials - login", func(t *testing.T) {

		authData := "basic " + base64.StdEncoding.EncodeToString([]byte("user2:"+password))
		req.Header.Set(echo.HeaderAuthorization, authData)
		assert.EqualError(t, h(c), pg.ErrNoRows.Error())
	})

}
