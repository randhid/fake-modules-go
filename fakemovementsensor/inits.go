package fakemovementsensor

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/movementsensor"
	"go.viam.com/rdk/resource"
)

const (
	movementsensorName = "fake-movementsensor"
	emptyName          = "empty-movementsensor"
	staticName         = "static-movementsensor"
)

var (
	EmptyModel  = common.FakesFamily.WithModel(emptyName)
	Model       = common.FakesFamily.WithModel(movementsensorName)
	StaticModel = common.FakesFamily.WithModel(staticName)
)

func init() {
	resource.RegisterComponent(movementsensor.API, StaticModel, resource.Registration[movementsensor.MovementSensor, resource.NoNativeConfig]{
		Constructor: newStaticMovementSensor,
	})

	resource.RegisterComponent(movementsensor.API, EmptyModel, resource.Registration[movementsensor.MovementSensor, resource.NoNativeConfig]{
		Constructor: newEmptyMovementSensor,
	})

	resource.RegisterComponent(movementsensor.API, Model, resource.Registration[movementsensor.MovementSensor, Config]{
		Constructor: newFakeMovementSensor,
	})
}
