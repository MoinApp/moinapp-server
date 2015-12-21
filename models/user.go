package models

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string
}

func (u *User) IsResult() bool {
	return (u.Password != nilUser.Password)
}

var (
	nilUser = &User{
		Password: "~error~", // this should be a never-reached hash
	}
)

func CreateUser(name, password, email string) *User {
	// create hashes
	passwordHash := sha256.Sum256([]byte(password))
	password = fmt.Sprintf("%x", passwordHash[:sha256.Size])
	emailHash := md5.Sum([]byte(email))
	email = fmt.Sprintf("%x", emailHash[:md5.Size])

	user := &User{
		Name:     name,
		Password: password,
		Email:    email,
	}

	db.Create(user)

	return user
}

func IsUsernameTaken(username string) bool {
	var query = &User{
		Name: username,
	}
	var result = &User{Password: "error"}

	db.Where(query).First(result)
	return (result.Password != "error")
}
