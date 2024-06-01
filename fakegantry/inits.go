package fakegantry

import (
	"fake-modules-go/common"
	"fmt"

	"go.viam.com/rdk/components/gantry"
	"go.viam.com/rdk/resource"
)

var (
	Model       = common.FakesFamily.WithModel(gantryName)
	StaticModel = common.FakesFamily.WithModel(staticName)
	EmptyModel  = common.FakesFamily.WithModel(emptyName)
)

const (
	staticName = "static-gantry"
	emptyName  = "empty-gantry"
	gantryName = "fake-gantry"
)

func init() {
	fmt.Println(EmptyModel)
	resource.RegisterComponent(gantry.API, EmptyModel, resource.Registration[gantry.Gantry, resource.NoNativeConfig]{
		Constructor: newEmptyGantry,
	})

	resource.RegisterComponent(gantry.API, Model, resource.Registration[gantry.Gantry, Config]{
		Constructor: newFakeGantry,
	})

	resource.RegisterComponent(gantry.API, StaticModel, resource.Registration[gantry.Gantry, resource.NoNativeConfig]{
		Constructor: newStaticGantry,
	})
}
