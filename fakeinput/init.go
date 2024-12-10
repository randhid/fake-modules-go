package fakeinput

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/input"
	"go.viam.com/rdk/resource"
)

const (
	inputName       = "fake-input"
	emptyInputeName = "empty-input"
)

var (
	Model      = common.FakesFamily.WithModel(inputName)
	EmptyModel = common.FakesFamily.WithModel(emptyInputeName)
)

func init() {
	resource.RegisterComponent(input.API, Model, resource.Registration[input.Controller, *Config]{
		Constructor: newFakeInput,
	})

	resource.RegisterComponent(input.API, EmptyModel, resource.Registration[input.Controller, resource.NoNativeConfig]{
		Constructor: newEmptyInput,
	},
	)
}
