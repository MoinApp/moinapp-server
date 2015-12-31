package push

import (
	"log"
)

type ApplePushNotifications struct {
}

func (apns ApplePushNotifications) Init() error {
	return nil
}

func (apns ApplePushNotifications) SendPush(data ...interface{}) error {
	log.Printf("APN: %v\n", data...)
	return nil
}
