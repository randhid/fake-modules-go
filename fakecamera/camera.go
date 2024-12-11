package fakecamera

import (
	"context"
	"fake-modules-go/common"
	"image"

	"sync"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/gostream"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
)

const cameraName = "fake-camera"

var Model = common.FakesFamily.WithModel(cameraName)

func init() {
	resource.RegisterComponent(camera.API, Model, resource.Registration[camera.Camera, resource.NoNativeConfig]{
		Constructor: newFakeCamera,
	})
}

type fake struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger

	mu sync.Mutex
}

func newFakeCamera(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	camera.Camera, error,
) {
	f := &fake{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return f, nil
}

func (f *fake) Image(context.Context, string, map[string]interface{}) ([]byte, camera.ImageMetadata, error) {
	return nil, camera.ImageMetadata{}, grpc.UnimplementedError
}

func (f *fake) Images(context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	return []camera.NamedImage{}, resource.ResponseMetadata{}, grpc.UnimplementedError
}

func (f *fake) NextPointCloud(context.Context) (pointcloud.PointCloud, error) {
	return pointcloud.New(), grpc.UnimplementedError
}

func (f *fake) Properties(context.Context) (camera.Properties, error) {
	return camera.Properties{}, grpc.UnimplementedError
}

func (f *fake) Stream(context.Context, ...gostream.ErrorHandler) (gostream.MediaStream[image.Image], error) {
	return nil, grpc.UnimplementedError
}
