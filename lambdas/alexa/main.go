package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	smarthome "github.com/ragurney/go-alexa-smarthome"
	"github.com/ragurney/hue-config/animations/sunrise"
)

func main() {
	sm := smarthome.New(smarthome.AuthorizationFunc(func(req smarthome.AcceptGrantRequest) error {
		return nil
	}))

	// Initiate Sunrise
	sm.AddDevice(sunrise.New())

	lambda.Start(sm.Handle)
}
