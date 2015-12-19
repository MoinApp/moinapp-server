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
	password = string(sha256.Sum256(password))
	email = string(md5.Sum(email))

	user := &User{
		Name:     name,
		Password: password,
		Email:    email,
	}

	db.Create(user)

	return user
}
