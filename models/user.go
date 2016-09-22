package models

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Name       string `sql:"unique; index"`
	Password   string `sql:"not null"`
	Email      string `sql:"not null"`
	PrivateKey string `sql:"size:4096"` // Needed for API v4

	PushTokens []PushToken
	Recents    []RecentMoin
}

type RecentMoin struct {
	gorm.Model
	UserID uint

	User     User      `sql:"not null"`
	LastMoin time.Time `sql:"not null"`
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
	db.Model(u).Association("PushTokens").Append(token)

	return u
}

func (u *User) GetPushTokens() []PushToken {
	var tokens []PushToken

	db.Model(u).Related(&tokens)

	return tokens
}

func (u *User) GetRecents() []*User {
	var recents []RecentMoin

	db.Model(u).Related(&recents)

	users := make([]*User, len(recents))
	for i, recent := range recents {
		user := nilUser()
		db.Model(&recent).Related(user)
		users[i] = user
	}

	return users
}

func (u *User) AddRecentUser(newRecent *User) *User {
	moin := RecentMoin{
		User:     *newRecent,
		LastMoin: time.Now(),
	}

	var recents []RecentMoin

	db.Model(u).Related(&recents)

	for i, recent := range recents {
		recentUser := nilUser()
		db.Model(&recent).Related(recentUser)

		if recentUser.Name == newRecent.Name {
			recents = append(recents[:i], recents[i+1:]...)
			break
		}
	}

	recents = append([]RecentMoin{moin}, recents...)

	db.Model(u).Association("Recents").Replace(recents)

	return u
}

func IsUsernameTaken(username string) bool {
	result := FindUserByName(username)
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

	db.Where("id = ?", id).First(&result)

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
