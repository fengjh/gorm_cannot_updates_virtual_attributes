package main

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"

	// Using postgres sql driver
	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
)

var (
	// DB returns a gorm.DB interface, it is used to access to database
	DB *gorm.DB
)

type User struct {
	gorm.Model
	Name              string `sql:"not null"`
	Password          string `sql:"-"`
	EncryptedPassword string `sql:"not null" json:"-"`
}

// BeforeSave is a hook function provided by gorm package. It is used to:
// - update user password
func (user *User) BeforeSave(db *gorm.DB) (err error) {
	if user.Password != "" {
		if err = user.setEncryptedPassword(); err != nil {
			return
		}
	}

	return
}

func (user *User) setEncryptedPassword() error {
	pw, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
	if err != nil {
		return err
	}
	user.EncryptedPassword = string(pw)
	user.Password = ""
	return nil
}

// ErrInvalidEmailOrPassword returns an error. It's using in
// models.UserAuthenticate function when authenticate failure.
var ErrInvalidEmailOrPassword = errors.New("invalid username or password")

// UserAuthenticate receives a name and a password. Then find the user
// with name in database and validate the password. And returns a user
// instance and an error.
func UserAuthenticate(name string, password string) (User, error) {
	var user User
	if err := DB.Where("name = ?", name).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return user, ErrInvalidEmailOrPassword
		}
		return user, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)); err != nil {
		return user, ErrInvalidEmailOrPassword
	}
	return user, nil
}

func init() {
	initDB()
	migrate()
}

func initDB() {
	var err error
	var db *gorm.DB

	dbParams := os.Getenv("DB_PARAMS")
	if dbParams == "" {
		panic(errors.New("DB_PARAMS environment variable not set"))
	}

	db, err = gorm.Open("postgres", fmt.Sprintf(dbParams))
	if err == nil {
		DB = db
	} else {
		panic(err)
	}
}

func migrate() {
	DB.DropTableIfExists(&User{})
	DB.AutoMigrate(&User{})
}
