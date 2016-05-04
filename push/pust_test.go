package push

import (
	"math/rand"
	"os"
	"testing"

	"github.com/MoinApp/moinapp-server/models"
)

func TestMain(m *testing.M) {
	models.InitDB(false)

	os.Exit(m.Run())
}

func TestEnvironmentVariables(t *testing.T) {
	const v = "TEST_VALUE"

	os.Setenv("GCM_API_KEY", v)
	gcmAPIKey := getGCMAPIKey()
	if gcmAPIKey != v {
		t.Errorf("Wrong GCM API key. Expected: %v. Got: %v.", v, gcmAPIKey)
	}

	os.Setenv("APN_CERT", v)
	apnCert := getAPNSCertificate()
	if apnCert != v {
		t.Errorf("Wrong APNS certificate. Expected: %v. Got: %v.", v, apnCert)
	}
}

func TestInit(t *testing.T) {
	InitPushServices(false)
	// this will never fail
}

func TestRandomSoundFileName(t *testing.T) {
	correlations := map[int64]string{
		0:  "moin4.wav",
		1:  "moin3.wav",
		3:  "moin1.wav",
		15: "moin5.wav",
	}

	if len(correlations) < len(SoundFiles) {
		t.Errorf("Incorrect length of test subjects. Expected: %v. Got: %v.", len(SoundFiles), len(correlations))
	}

	for k, v := range correlations {
		rand.Seed(k)

		sound := randomSoundFilename()
		if sound != v {
			t.Errorf("Unexpected sound file name. Expected: %q. Got: %q.", v, sound)
		}
	}
}

func TestSendNotificationToAll_Dry(t *testing.T) {
	n := PushNotification{}
	m := []models.PushToken{
		models.PushToken{
			Type:  models.APNToken,
			Token: "null",
		},
		models.PushToken{
			Type:  models.GCMToken,
			Token: "null",
		},
	}

	SendPushNotificationToAll(m, &n)
}

func TestSendMoinNotificationToUser_Dry(t *testing.T) {
	u := models.CreateUser("TestSendMoinNotificationToUser_Dry", "TestSendMoinNotificationToUser_Dry", "TestSendMoinNotificationToUser_Dry")
	t1 := models.NewPushToken(models.APNToken, "null1")
	t2 := models.NewPushToken(models.GCMToken, "null2")
	u.AddPushToken(t1)
	u.AddPushToken(t2)

	SendMoinNotificationToUser(u, u)
}
