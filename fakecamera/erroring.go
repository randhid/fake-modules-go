package fakecamera

import (
	"context"

	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage"
	"go.viam.com/rdk/rimage/transform"
)

type erroring struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logging logging.Logger
}

func newErroringCamera(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	camera.Camera, error,
) {
	return &erroring{
		Named:   conf.ResourceName().AsNamed(),
		logging: logger,
	}, nil
}

func (e *erroring) Image(context.Context, string, map[string]interface{}) ([]byte, camera.ImageMetadata, error) {
	img := rimage.NewImage(320, 240)
	return rimage.ImageToUInt8Buffer(img, false), camera.ImageMetadata{}, grpc.UnimplementedError
}

func (e *erroring) Images(context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	return []camera.NamedImage{
		{SourceName: "cam", Image: rimage.NewImage(320, 240)},
	}, resource.ResponseMetadata{}, grpc.UnimplementedError
}

func (e *erroring) NextPointCloud(context.Context) (pointcloud.PointCloud, error) {
	return pointcloud.New(), grpc.UnimplementedError
}

func (e *erroring) Properties(context.Context) (camera.Properties, error) {
	return camera.Properties{
		IntrinsicParams: &transform.PinholeCameraIntrinsics{
			Width:  320,
			Height: 240,
		},
		SupportsPCD: true,
	}, grpc.UnimplementedError
}
