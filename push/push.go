package push

import (
	"errors"
	"fmt"
	"github.com/Coccodrillo/apns"
	"github.com/MoinApp/moinapp-server/models"
	"log"
	"math/rand"
)

type PushNotification struct {
	Message    string
	Sound      string
	BadgeCount int
	Payload    map[string]interface{}
}

var (
	ErrUnknownPushTokenType = errors.New("Unknown push token type.")
	SoundFiles              = []string{
		"moin1.wav",
		"moin3.wav",
		"moin4.wav",
		"moin5.wav",
	}
)

var apnsClient *apns.Client

func InitPushServices(isProduction bool) {
	initApplePushNotificationService(isProduction)
}
func initApplePushNotificationService(isProduction bool) {
	var gateway string
	if isProduction {
		gateway = "gateway.push.apple.com:2195"
	} else {
		gateway = "gateway.sandbox.push.apple.com:2195"
	}

	certificateFile := "TODO"
	keyFile := "TODO"
	apnsClient = apns.NewClient(gateway, certificateFile, keyFile)
}

func SendMoinNotificationToUser(receiver, sender *models.User) {
	pushTokens := receiver.GetPushTokens()

	/* if len(pushTokens) < 1 {
		return
	} */

	notification := &PushNotification{
		Message:    "by " + sender.Name,
		Sound:      randomSoundFilename(),
		BadgeCount: 1,
	}

	SendPushNotificationToAll(pushTokens, notification)
}

func randomSoundFilename() string {
	return SoundFiles[rand.Intn(len(SoundFiles))]
}

func SendPushNotificationToAll(tokens []models.PushToken, notification *PushNotification) {
	for _, token := range tokens {
		SendPushNotification(token, notification)
	}
}
func SendPushNotification(token models.PushToken, notification *PushNotification) {
	fmt.Printf("Send notification %+v to token %+v...\n", notification, token)

	switch token.Type {
	case models.APNToken:
		sendApplePushNotification(token.Token, notification)
	default:
		panic(ErrUnknownPushTokenType)
	}
}
func sendApplePushNotification(token string, notification *PushNotification) {
	payload := apns.NewPayload()
	payload.Alert = notification.Message
	payload.Badge = notification.BadgeCount
	payload.Sound = notification.Sound

	pn := apns.NewPushNotification()
	pn.AddPayload(payload)

	for k, v := range notification.Payload {
		pn.Set(k, v)
	}

	s, _ := pn.PayloadString()
	fmt.Printf("APNS: %v\n", s)

	response := apnsClient.Send(pn)
	if !response.Success {
		log.Printf("APNS error: %v.", response.Error)
	}
}
