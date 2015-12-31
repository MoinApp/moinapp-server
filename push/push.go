package push

import (
	"errors"
	"github.com/MoinApp/moinapp-server/models"
)

type PushNotificationProvider interface {
	Init() error
	// TODO define "data" as API
	SendPush(data ...interface{}) error
}

var (
	ErrUnknownPushTokenType = errors.New("Unknown push token type.")
)

func SendPushNotification(token models.PushToken, message string) {
	var provider PushNotificationProvider

	switch token.Type {
	case models.APNToken:
		provider = ApplePushNotifications{}
	case models.GCMToken:

	default:
		panic(ErrUnknownPushTokenType)
	}

	provider.Init()
	provider.SendPush(token, message)
}
