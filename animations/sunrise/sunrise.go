// Package sunrise specifies the behavior of the Hue sunrise animation
package sunrise

import (
	"time"

	smarthome "github.com/ragurney/go-alexa-smarthome"
	"github.com/ragurney/huego"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type sunriseHandler struct {
	val interface{}
}

const (
	// StartTransitionTime transition time used to set initial state. Should always be 0
	StartTransitionTime uint16 = 0
	// EndTransitionTime transition time used to set duration of sunrise. Represents multiples of 100ms. Max value is
	// 65535 (~ 109 minutes)
	EndTransitionTime uint16 = 6000
	// StartBrightness is the starting brigntness of the sunrise, can be 0 - 254
	StartBrightness uint8 = 0
	// EndBrightness is the end brigntness of the sunrise, can be 0 - 254
	EndBrightness uint8 = 254
	// StartHue is the start hue value of the sunrise, can be 0 - 65535
	StartHue uint16 = 3091
	// EndHue is the start hue value of the sunrise, can be 0 - 65535
	EndHue uint16 = 11500
	// StartSat is the start saturation of the sunrise, can be 254 (most saturated) - 0 (least saturated, a.k.a white)
	StartSat uint8 = 222
	// EndSat is the end saturation of the sunrise, can be 254 (most saturated) - 0 (least saturated, a.k.a white)
	EndSat uint8 = 0
)

// GetValue handles returning the value of Sunrise. Since we want Sunrise to simply be a trigger, we do nothing
// and always retrn 'off'
func (mockHandler *sunriseHandler) GetValue() (interface{}, error) {
	log.Debug().Str("Handler", "Sunrise").Str("Function", "GetValue").Msgf("Getting value: %+v", mockHandler.val)
	return mockHandler.val, nil
}

// SetValue is called when Alexa changes the value of the Sunrise device. It only responds to the "ON" value
// which triggers the sunrise animation. Since we want sunrise to be a trigger, we do not change the value  of the
// device
func (mockHandler *sunriseHandler) SetValue(val interface{}, token string) error {
	log.Debug().Str("Handler", "Sunrise").Str("Function", "SetValue").Msgf("ReceivedValue value: %+v", val)

	mockHandler.val = val

	if val == "ON" {
		bridge := huego.New(token)
		_, err := bridge.Login("libre-hue")
		if err != nil {
			log.Error().Str("Handler", "Sunrise").Str("Function", "SetValue").Msgf("Error logging into Hue: %s", err)
			return err
		}

		groups, err := bridge.GetGroups()
		if err != nil {
			log.Error().Str("Handler", "Sunrise").Str("Function", "SetValue").Msgf("Error getting Hue light Groups: %s", err)
			return err // TODO: add support for Alexa to say what the issue is to the user
		}

		// \[T]/
		// Sends sunrise command to the first group, whatever that is. Can customize this to your own needs.
		// TODO: make configurable via Alexa Smart Home, either create custom device or add config options to 'light'
		groups[0].SetState(huego.State{
			On:             true,
			TransitionTime: StartTransitionTime,
			Bri:            StartBrightness,
			Hue:            StartHue,
			Sat:            StartSat,
		})
		time.Sleep(1 * time.Second) // Give Hue time to do its thing (group calls need 1 second to propagate)
		groups[0].SetState(huego.State{TransitionTime: EndTransitionTime, Bri: EndBrightness, Hue: EndHue, Sat: EndSat})
	}

	return nil
}

// UpdateChannel is unused, but necessary for the AbstractDevice interface
func (mockHandler *sunriseHandler) UpdateChannel() <-chan interface{} {
	return nil
}

// New creates a new capability for a Sunrise
func New() *smarthome.AbstractDevice {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Defines the device that will be discoverable via Alexa Smart Home
	sunriseDevice := smarthome.NewAbstractDevice(
		"1",                 // id
		"Sunrise",           // name
		"sunrise",           // manufactureName - can be anything unique
		"Sunrise Animation", // description
	)

	sunriseDevice.AddDisplayCategory("LIGHT")
	capability := sunriseDevice.NewCapability("PowerController")
	capability.AddPropertyHandler("powerState", &sunriseHandler{val: "OFF"})

	return sunriseDevice
}
