package fakemotor

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/motor"
	"go.viam.com/rdk/resource"
)

const (
	emptyName  = "empty-motor"
	staticName = "static-motor"
	motorName  = "fake-motor"
)

var (
	EmptyModel  = common.FakesFamily.WithModel(emptyName)
	Model       = common.FakesFamily.WithModel(motorName)
	StaticModel = common.FakesFamily.WithModel(staticName)
)

func init() {
	resource.RegisterComponent(motor.API, StaticModel, resource.Registration[motor.Motor, resource.NoNativeConfig]{
		Constructor: newStaticMotor,
	})

	resource.RegisterComponent(motor.API, Model, resource.Registration[motor.Motor, Config]{
		Constructor: newFakeMotor,
	})

	resource.RegisterComponent(motor.API, EmptyModel, resource.Registration[motor.Motor, resource.NoNativeConfig]{
		Constructor: newEmptyMotor,
	})
}
