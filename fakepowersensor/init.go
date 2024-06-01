package fakepowersensor

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/powersensor"
	"go.viam.com/rdk/resource"
)

const (
	powersensorName = "fake-powersensor"
	emptyName       = "empty-powersensor"
	staticName      = "static-powersensor"
)

var (
	Model       = common.FakesFamily.WithModel(powersensorName)
	EmptyModel  = common.FakesFamily.WithModel(emptyName)
	StaticModel = common.FakesFamily.WithModel(staticName)
)

func init() {
	resource.RegisterComponent(powersensor.API, Model, resource.Registration[powersensor.PowerSensor, *Config]{
		Constructor: newFakePowerSensor,
	})
	resource.RegisterComponent(powersensor.API, EmptyModel, resource.Registration[powersensor.PowerSensor, *Config]{
		Constructor: newEmptyPowerSensor,
	})
	resource.RegisterComponent(powersensor.API, StaticModel, resource.Registration[powersensor.PowerSensor, *Config]{
		Constructor: newStaticPowerSensor,
	})
}
