package push

import (
	"errors"
	"log"
	"math/rand"
	"os"

	"github.com/Coccodrillo/apns"
	"github.com/MoinApp/moinapp-server/models"
	"github.com/alexjlockwood/gcm"
)

type PushNotification struct {
	Message    string
	Sound      string
	BadgeCount int
	Payload    map[string]interface{}
}

const (
	GoogleCloudMessagingRetries = 1
)

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
var gcmSender *gcm.Sender

func getGCMAPIKey() string {
	return os.Getenv("GCM_API_KEY")
}
func getAPNSCertificate() string {
	return os.Getenv("APN_CERT")
}

func InitPushServices(isProduction bool) {
	initApplePushNotificationService(isProduction)
	initGoogleCloudMessaging()
}
func initApplePushNotificationService(isProduction bool) {
	var gateway string
	if isProduction {
		gateway = "gateway.push.apple.com:2195"
	} else {
		gateway = "gateway.sandbox.push.apple.com:2195"
	}

	// TODO: this is an PFX file. Need to separate out PEM and KEY
	certificateBase64 := getAPNSCertificate()
	keyBase64 := getAPNSCertificate()
	apnsClient = apns.BareClient(gateway, certificateBase64, keyBase64)
}
func initGoogleCloudMessaging() {
	gcmSender = &gcm.Sender{
		ApiKey: getGCMAPIKey(),
	}
}

func SendMoinNotificationToUser(receiver, sender *models.User) {
	pushTokens := receiver.GetPushTokens()

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
	//fmt.Printf("Send notification %+v to token %+v...\n", notification, token)

	switch token.Type {
	case models.APNToken:
		sendApplePushNotification(token.Token, notification)
	case models.GCMToken:
		sendGoogleCloudMessagingNotification(token.Token, notification)
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

	pn.PayloadString()
	//fmt.Printf("APNS: %v\n", s)

	response := apnsClient.Send(pn)
	if !response.Success {
		log.Printf("APNS error: %v.", response.Error)
	}
}
func sendGoogleCloudMessagingNotification(registrationID string, notification *PushNotification) {
	message := gcm.NewMessage(notification.Payload, registrationID)

	_, err := gcmSender.Send(message, GoogleCloudMessagingRetries)
	if err != nil {
		log.Printf("GCM error: %v.", err)
	}
	//fmt.Printf("GCM reponse: %+v\n", response)
}
