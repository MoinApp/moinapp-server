package push

import (
	"log"
)

type GoogleCloudMessages struct {
}

func (gcm GoogleCloudMessages) Init() error {
	return nil
}

func (gcm GoogleCloudMessages) SendPush(data ...interface{}) error {
	log.Printf("GCM: %v\n", data...)
	return nil
}
