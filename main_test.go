package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initDB()
	migrate()

	retCode := m.Run()

	os.Exit(retCode)
}

func TestCreateUser(t *testing.T) {
	user := createUser(t)

	authedUser, err := UserAuthenticate(user.Name, "secret")
	assertNoErr(t, err)

	if authedUser.ID != user.ID {
		t.Fatal("GORM Create user failure")
	}
}

func TestUpdateUserPassword(t *testing.T) {
	user := createUser(t)

	newUser := User{Name: "newusername", Password: "newpassword"}

	err := DB.Model(user).Updates(newUser).Error
	assertNoErr(t, err)

	authedUser, err := UserAuthenticate(newUser.Name, newUser.Password)
	assertNoErr(t, err)

	if authedUser.ID != user.ID {
		t.Fatal("GORM Updates password failure")
	}
}

func createUser(t *testing.T) (user *User) {
	user = &User{Name: "username", Password: "secret"}
	err := DB.Create(user).Error
	assertNoErr(t, err)

	return
}

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
