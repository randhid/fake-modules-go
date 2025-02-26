package fakecamera

import (
	"bytes"
	"context"
	"image"
	"image/png"
	"time"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/pointcloud"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/rimage/transform"
)

type fake struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger
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
	img := image.NewGray(image.Rect(0, 0, 320, 240))
	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, camera.ImageMetadata{}, err
	}
	return buf.Bytes(),
		camera.ImageMetadata{MimeType: "image/png"},
		nil
}

func (f *fake) Images(context.Context) ([]camera.NamedImage, resource.ResponseMetadata, error) {
	return []camera.NamedImage{
			{SourceName: "cam", Image: image.NewGray(image.Rect(0, 0, 320, 240))},
			{SourceName: "cam", Image: image.NewAlpha(image.Rect(0, 0, 320, 240))},
			{SourceName: "cam", Image: image.NewRGBA(image.Rect(0, 0, 320, 240))},
			{SourceName: "cam", Image: image.NewYCbCr(image.Rect(0, 0, 320, 240), image.YCbCrSubsampleRatio444)},
		},
		resource.ResponseMetadata{CapturedAt: time.Now()},
		nil
}

func (f *fake) NextPointCloud(context.Context) (pointcloud.PointCloud, error) {
	return pointcloud.NewBasicOctree(r3.Vector{X: 0, Y: 0, Z: 0}, 5.0)
}

func (f *fake) Properties(context.Context) (camera.Properties, error) {
	return camera.Properties{
		IntrinsicParams: &transform.PinholeCameraIntrinsics{
			Width:  320,
			Height: 240,
			Fx:     10,
			Fy:     10,
			Ppx:    10,
			Ppy:    10,
		},
		DistortionParams: &transform.BrownConrady{
			RadialK1:     10,
			RadialK2:     20,
			RadialK3:     30,
			TangentialP1: 5,
			TangentialP2: 7,
		},
		MimeTypes:   []string{"image/png", "image/rgba", "image/ycbcr", "image/gray"},
		SupportsPCD: true,
		FrameRate:   30,
		ImageType:   "gray",
	}, nil
}
