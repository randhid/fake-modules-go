//go:build !no_cgo

// Package fake implements a fake audio input.
package fake

import (
	"context"
	"fake-modules-go/common"

	"github.com/pion/mediadevices/pkg/prop"
	"github.com/pion/mediadevices/pkg/wave"

	"go.viam.com/rdk/components/audioinput"
	"go.viam.com/rdk/gostream"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

const armName = "fake-audio-in"

var Model = common.FakesFamily.WithModel(armName)

func init() {
	resource.RegisterComponent(audioinput.API, Model, resource.Registration[audioinput.AudioInput, resource.NoNativeConfig]{
		Constructor: newFakeAudio,
	})
}

func newFakeAudio(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	audioinput.AudioInput, error,
) {

	i := &fake{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return i, nil

}

// audioInput is a fake audioinput that always returns the same chunk.
type fake struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	logger logging.Logger
}

const (
	latencyMillis = 20
	samplingRate  = 48000
	channelCount  = 1
)

func (f *fake) Read(ctx context.Context) (wave.Audio, func(), error) {
	return nil, func() {}, grpc.UnimplementedError
}

func (f *fake) MediaProperties(_ context.Context) (prop.Audio, error) {
	return prop.Audio{}, grpc.UnimplementedError
}

func (f *fake) Stream(context.Context, ...gostream.ErrorHandler) (gostream.MediaStream[wave.Audio], error) {
	return nil, grpc.UnimplementedError
}
