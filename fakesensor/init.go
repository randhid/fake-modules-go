package fakesensor

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/resource"
)

const (
	sensorName   = "fake-sensor"
	emptyName    = "empty-sensor"
	staticName   = "static-sensor"
	waitingName  = "waiting-sensor"
	erroringName = "erroring-sensor"
)

var (
	Model         = common.FakesFamily.WithModel(sensorName)
	EmptyModel    = common.FakesFamily.WithModel(emptyName)
	StaticModel   = common.FakesFamily.WithModel(staticName)
	WaitingModel  = common.FakesFamily.WithModel(waitingName)
	ErrroingModel = common.FakesFamily.WithModel(erroringName)
)

func init() {
	resource.RegisterComponent(sensor.API, Model, resource.Registration[sensor.Sensor, resource.NoNativeConfig]{
		Constructor: newFakeSensor,
	})
	resource.RegisterComponent(sensor.API, EmptyModel, resource.Registration[sensor.Sensor, resource.NoNativeConfig]{
		Constructor: newEmptySensor,
	})
	resource.RegisterComponent(sensor.API, StaticModel, resource.Registration[sensor.Sensor, resource.NoNativeConfig]{
		Constructor: newStaticSensor,
	})
	resource.RegisterComponent(sensor.API, WaitingModel, resource.Registration[sensor.Sensor, WaitingConfig]{
		Constructor: newWaitingSensor,
	})
	resource.RegisterComponent(sensor.API, ErrroingModel, resource.Registration[sensor.Sensor, resource.NoNativeConfig]{
		Constructor: newErroringSensor,
	})
}
