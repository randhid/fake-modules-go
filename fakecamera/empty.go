package fakecamera

import (
	"context"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
)

type empty struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger
}

func newEmptyCamera(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	camera.Camera, error,
) {
	return &empty{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}, nil
}

func (e *empty) Image(context.Context, string, map[string]interface{}) ([]byte, camera.ImageMetadata, error) {
	return nil, camera.ImageMetadata{}, nil
}

func (e *empty) Images(context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	return nil, resource.ResponseMetadata{}, nil
}
func (e *empty) NextPointCloud(context.Context) (pointcloud.PointCloud, error) {
	return nil, nil
}

func (e *empty) Properties(context.Context) (camera.Properties, error) {
	return camera.Properties{}, nil
}
