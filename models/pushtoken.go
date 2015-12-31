package models

import (
	"github.com/jinzhu/gorm"
)

type TokenType uint

const (
	APNType TokenType = 0
	GCMType           = 1
)

type PushToken struct {
	gorm.Model
	UserID uint

	Token string    `sql:"not null;unique"`
	Type  TokenType `sql:"not null"`
}

func NewPushToken(t TokenType, token string) *PushToken {
	pushToken := &PushToken{
		Token: token,
		Type:  t,
	}

	db.Create(pushToken)

	return pushToken
}
