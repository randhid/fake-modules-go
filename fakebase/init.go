package fakebase

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/base"
	"go.viam.com/rdk/resource"
)

const (
	baseName  = "fake-base"
	emptyName = "empty-base"
)

var (
	Model      = common.FakesFamily.WithModel(baseName)
	EmptyModel = common.FakesFamily.WithModel(emptyName)
)

func init() {
	resource.RegisterComponent(base.API, Model, resource.Registration[base.Base, *Config]{
		Constructor: newFakeBase,
	})

	resource.RegisterComponent(base.API, EmptyModel, resource.Registration[base.Base, resource.NoNativeConfig]{
		Constructor: newEmptyBase,
	})
}
