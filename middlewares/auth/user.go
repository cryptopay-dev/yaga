package auth

import (
	"errors"
	"time"

	"github.com/go-pg/pg/orm"
	"golang.org/x/crypto/bcrypt"
)

var (
	errUsernameAlreadyTaken = errors.New("username already taken")
)

type User struct {
	ID        string `sql:",pk"`
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func NewUser(db orm.DB, username, password string) (*User, error) {
	var (
		err   error
		count int
		user  = &User{
			Username:  username,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)

	if count, err = db.Model(&User{}).Where("username = ?", username).Count(); err != nil {
		return nil, err
	} else if count > 0 {
		return nil, errUsernameAlreadyTaken
	}

	if user.Password, err = user.HashPassword(password); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) ByName(db orm.DB, username string) error {
	return db.Model(u).Where("username = ?", username).First()
}

func (u *User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *User) CheckPasswordHash(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil, err
}
