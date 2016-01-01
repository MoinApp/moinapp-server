package models

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Name       string `sql:"unique; index"`
	Password   string `sql:"not null"`
	Email      string `sql:"not null"`
	PushTokens []PushToken
	PrivateKey string `sql:"size:4096"`
}

func (u *User) IsResult() bool {
	return (u.Password != nilUser().Password)
}

func nilUser() *User {
	return &User{
		Password: "~error~", // this should be a never-reached hash
	}
}

func (u *User) AddPushToken(token *PushToken) *User {
	db.Model(&u).Association("PushTokens").Append(token)

	return u
}

func (u *User) GetPushTokens() []PushToken {
	var tokens []PushToken

	db.Model(&u).Related(&tokens)

	return tokens
}

func IsUsernameTaken(username string) bool {
	var query = &User{
		Name: username,
	}
	var result = nilUser()

	db.Where(query).First(result)
	return result.IsResult()
}

func CreateUser(name, password, email string) *User {
	// create hashes
	password = getPasswordHash(password)
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

func getPasswordHash(password string) string {
	passwordHash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", passwordHash[:sha256.Size])
}

func FindUserByName(username string) *User {
	var result = nilUser()

	db.Where("name LIKE ?", username).First(result)
	return result
}

func FindUsersByName(username string) []*User {
	var result []*User

	db.Where("name LIKE ?", "%"+username+"%").Find(&result)
	return result
}

func FindUserById(id interface{}) *User {
	var result = nilUser()

	db.First(result, id)

	return result
}

func FindUserWithCredentials(username, password string) *User {
	var result = nilUser()

	query := &User{
		Name:     username,
		Password: getPasswordHash(password),
	}

	db.Where(query).First(result)

	return result
}

func SaveUser(user *User) {
	db.Save(user)
}
