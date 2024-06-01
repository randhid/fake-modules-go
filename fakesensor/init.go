package fakesensor

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/resource"
)

const (
	sensorName = "fake-sensor"
	emptyName  = "empty-sensor"
	staticName = "static-sensor"
)

var (
	Model = common.FakesFamily.WithModel(sensorName)
	EmptyModel = common.FakesFamily.WithModel(emptyName)
	StaticModel = common.FakesFamily.WithModel(staticName)
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
}
