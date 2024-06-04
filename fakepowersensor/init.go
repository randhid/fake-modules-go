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
	waitingName     = "waiting-powersensor"
	errroingName    = "erroring-powersensor"
)

var (
	Model         = common.FakesFamily.WithModel(powersensorName)
	EmptyModel    = common.FakesFamily.WithModel(emptyName)
	StaticModel   = common.FakesFamily.WithModel(staticName)
	WaitingModel  = common.FakesFamily.WithModel(waitingName)
	ErroringModel = common.FakesFamily.WithModel(errroingName)
)

func init() {
	resource.RegisterComponent(powersensor.API, Model, resource.Registration[powersensor.PowerSensor, *Config]{
		Constructor: newFakePowerSensor,
	})
	resource.RegisterComponent(powersensor.API, EmptyModel, resource.Registration[powersensor.PowerSensor, resource.NoNativeConfig]{
		Constructor: newEmptyPowerSensor,
	})
	resource.RegisterComponent(powersensor.API, StaticModel, resource.Registration[powersensor.PowerSensor, *StaticConfig]{
		Constructor: newStaticPowerSensor,
	})
	resource.RegisterComponent(powersensor.API, WaitingModel, resource.Registration[powersensor.PowerSensor, *WaitingConfig]{
		Constructor: newWaitingPowerSensor,
	})
	resource.RegisterComponent(powersensor.API, ErroringModel, resource.Registration[powersensor.PowerSensor, resource.NoNativeConfig]{
		Constructor: newErroringPowerSensor,
	})
}
