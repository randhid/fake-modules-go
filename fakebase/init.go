package fakebase

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/base"
	"go.viam.com/rdk/resource"
)

const (
	baseName  = "fake-base"
	emptyName = "empty-base"
	nanName   = "nan-base"
)

var (
	Model      = common.FakesFamily.WithModel(baseName)
	EmptyModel = common.FakesFamily.WithModel(emptyName)
	NanModel   = common.FakesFamily.WithModel(nanName)
)

func init() {
	resource.RegisterComponent(base.API, Model, resource.Registration[base.Base, *Config]{
		Constructor: newFakeBase,
	})

	resource.RegisterComponent(base.API, EmptyModel, resource.Registration[base.Base, resource.NoNativeConfig]{
		Constructor: newEmptyBase,
	})

	resource.RegisterComponent(base.API, NanModel, resource.Registration[base.Base, resource.NoNativeConfig]{
		Constructor: newNanBase,
	})
}
