package fakeservo

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/resource"
)

const (
	servoName  = "fake-servo"
	staticName = "static-servo"
)

var (
	Model       = common.FakesFamily.WithModel(servoName)
	StaticModel = common.FakesFamily.WithModel(staticName)
)

func init() {
	resource.RegisterComponent(servo.API, Model, resource.Registration[servo.Servo, resource.NoNativeConfig]{
		Constructor: newFakeServo,
	})
	resource.RegisterComponent(servo.API, StaticModel, resource.Registration[servo.Servo, resource.NoNativeConfig]{
		Constructor: newStaticServo,
	})
}
