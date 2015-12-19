package models

import (
	"crypto/md5"
	"crypto/sha256"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string
}

func CreateUser(name, password, email string) *User {
	// create hashes
	passwordHash := sha256.Sum256([]byte(password))
	password = string(passwordHash[:sha256.Size])
	emailHash := md5.Sum([]byte(email))
	email = string(emailHash[:md5.Size])

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
	var result = &User{Name: "error"}

	db.Where(query).First(result)
	return (result.Name != "error")
}
