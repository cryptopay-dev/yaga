package auth

import (
	"time"

	"github.com/cryptopay-dev/yaga/model"
	"github.com/go-pg/pg/orm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUsernameAlreadyTaken when username exists in DB
	ErrUsernameAlreadyTaken = errors.New("username already taken")
)

// User model
type User struct {
	ID        string `sql:",pk"`
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRequest to create new User
type UserRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

// NewUser from params, that check unique `username` and
// return new user, or error
func NewUser(db orm.DB, username, password string) (*User, error) {
	var (
		err   error
		found bool
		dt    = time.Now()
		user  = &User{
			Username:  username,
			CreatedAt: dt,
			UpdatedAt: dt,
		}
	)

	if found, err = model.Exist(db, model.Conditions{"username": username}, &User{}); err != nil {
		return nil, errors.Wrap(err, "auth NewUser failed")
	} else if found {
		return nil, ErrUsernameAlreadyTaken
	}

	if err = user.HashPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

// ByName gets user from DB
func (u *User) ByName(db orm.DB, username string) error {
	return model.FindOne(db, model.Conditions{"username": username}, u)
}

// HashPassword from input to user-model
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return errors.Wrap(err, "auth HasPassword failed")
	}
	u.Password = string(bytes)
	return nil
}

// CheckPasswordHash compare input and user password
func (u *User) CheckPasswordHash(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil, errors.Wrap(err, "auth CheckPasswordHash failed")
}
