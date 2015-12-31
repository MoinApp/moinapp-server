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
	User   User
	UserID uint

	Token string
	Type  TokenType
}

func NewPushToken(t TokenType, token string) *PushToken {
	pushToken := &PushToken{
		Token: token,
		Type:  t,
	}

	db.Create(pushToken)

	return pushToken
}
