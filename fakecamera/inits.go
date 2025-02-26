package fakecamera

import (
	"fake-modules-go/common"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/resource"
)

const cameraName = "fake-camera"
const emptyName = "empty"
const erroringName = "erroring"
const malformedName = "malformed"

var (
	Model     = common.FakesFamily.WithModel(cameraName)
	Empty     = common.FakesFamily.WithModel(emptyName)
	Erroring  = common.FakesFamily.WithModel(erroringName)
	Malformed = common.FakesFamily.WithModel(malformedName)
)

func init() {
	resource.RegisterComponent(camera.API, Model, resource.Registration[camera.Camera, resource.NoNativeConfig]{
		Constructor: newFakeCamera,
	})
	resource.RegisterComponent(camera.API, Empty, resource.Registration[camera.Camera, resource.NoNativeConfig]{
		Constructor: newEmptyCamera,
	})
	resource.RegisterComponent(camera.API, Erroring, resource.Registration[camera.Camera, resource.NoNativeConfig]{
		Constructor: newErroringCamera,
	})
	resource.RegisterComponent(camera.API, Malformed, resource.Registration[camera.Camera, resource.NoNativeConfig]{
		Constructor: newMalformedCamera,
	})
}
