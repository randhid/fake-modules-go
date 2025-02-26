package fakecamera

import (
	"context"
	"image"
	"math"
	"time"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage"
	"go.viam.com/rdk/rimage/transform"
)

type malformed struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger
}

func newMalformedCamera(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	camera.Camera, error,
) {
	return &malformed{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}, nil
}

func (m *malformed) Image(context.Context, string, map[string]interface{}) ([]byte, camera.ImageMetadata, error) {
	return []byte{0}, camera.ImageMetadata{MimeType: "rand"}, nil
}

func (m *malformed) Images(context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	return []camera.NamedImage{
			{SourceName: "rand", Image: rimage.NewImage(-1, -2)},
			{SourceName: "rand", Image: rimage.NewImage(int(math.NaN()), int(math.NaN()))},
			{SourceName: "rand", Image: image.NewAlpha(image.Rect(int(math.NaN()), int(math.NaN()), int(math.NaN()), int(math.NaN())))},
		}, resource.ResponseMetadata{
			CapturedAt: time.Time{},
		}, nil
}

func (m *malformed) NextPointCloud(context.Context) (pointcloud.PointCloud, error) {
	octree, _ := pointcloud.NewBasicOctree(r3.Vector{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, math.NaN())
	return octree, nil
}

func (m *malformed) Properties(context.Context) (camera.Properties, error) {
	return camera.Properties{
		IntrinsicParams: &transform.PinholeCameraIntrinsics{
			Width:  int(math.NaN()),
			Height: int(math.NaN()),
			Fx:     math.NaN(),
			Fy:     math.NaN(),
			Ppx:    math.NaN(),
			Ppy:    math.NaN(),
		},
		DistortionParams: &transform.BrownConrady{
			RadialK1:     math.NaN(),
			RadialK2:     math.NaN(),
			RadialK3:     math.NaN(),
			TangentialP1: math.NaN(),
			TangentialP2: math.NaN(),
		},
		SupportsPCD: true,
		ImageType:   "sleep-no-more",
		MimeTypes:   []string{"rand", "rand", "rand"},
		FrameRate:   float32(math.NaN()),
	}, nil
}
