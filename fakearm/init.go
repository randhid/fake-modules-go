package fakearm

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/arm"
	"go.viam.com/rdk/resource"
)

const (
	armName      = "fake-arm"
	emptyName    = "empty-arm"
	staticName   = "static-arm"
	erroringName = "erroring-arm"
)

var (
	Model         = common.FakesFamily.WithModel(armName)
	EmptyModel    = common.FakesFamily.WithModel(emptyName)
	StaticModel   = common.FakesFamily.WithModel(staticName)
	ErroringModel = common.FakesFamily.WithModel(erroringName)
)

func init() {
	resource.RegisterComponent(arm.API, Model, resource.Registration[arm.Arm, resource.NoNativeConfig]{
		Constructor: newFakeArm,
	})

	resource.RegisterComponent(arm.API, EmptyModel, resource.Registration[arm.Arm, resource.NoNativeConfig]{
		Constructor: newEmptyArm,
	})
	resource.RegisterComponent(arm.API, StaticModel, resource.Registration[arm.Arm, resource.NoNativeConfig]{
		Constructor: newStaticArm,
	})
	resource.RegisterComponent(arm.API, ErroringModel, resource.Registration[arm.Arm, resource.NoNativeConfig]{
		Constructor: newErroringArm,
	})
}
