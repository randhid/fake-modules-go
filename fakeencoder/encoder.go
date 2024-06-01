package fakeencoder

import (
	"context"
	"errors"
	"fake-modules-go/common"
	"fmt"
	"math"
	"sync"

	"go.viam.com/rdk/components/encoder"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

const (
	encoderName  = "fake-encoder"
	setIncrement = "setIncrement" // for DoCommand to make the increment larger or siwtch directions
)

var (
	Model = common.FakesFamily.WithModel(encoderName)
)

func init() {
	resource.RegisterComponent(encoder.API, Model, resource.Registration[encoder.Encoder, resource.NoNativeConfig]{
		Constructor: newFakeEncoder,
	})
}

type fake struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	resource.Named

	logger logging.Logger

	mu        sync.Mutex
	position  float64
	increment float64
}

func newFakeEncoder(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	encoder.Encoder, error,
) {
	f := &fake{
		Named:     conf.ResourceName().AsNamed(),
		logger:    logger,
		increment: 1.0,
	}

	return f, nil
}

func (f *fake) Position(ctx context.Context, encType encoder.PositionType, extra map[string]interface{}) (float64, encoder.PositionType, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.position += f.increment

	switch encType {
	case encoder.PositionTypeDegrees:
		return math.Mod(f.position, 360), encoder.PositionTypeDegrees, nil
	case encoder.PositionTypeTicks:
		return f.position, encoder.PositionTypeTicks, nil
	default:
		return math.NaN(), encoder.PositionTypeUnspecified, errors.New("unsupported encoder position types")
	}
}

func (f *fake) ResetPosition(ctx context.Context, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.position = 0
	return nil
}

func (f *fake) Properties(ctx context.Context, extra map[string]interface{}) (encoder.Properties, error) {
	return encoder.Properties{AngleDegreesSupported: true, TicksCountSupported: true}, nil
}

func (f *fake) DoCommand(ctx context.Context, req map[string]interface{}) (map[string]interface{}, error) {
	reqInc := req[setIncrement]
	increment, ok := reqInc.(float64)
	if !ok {
		return nil, fmt.Errorf("%s is not a float64", setIncrement)
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.increment = increment

	return map[string]interface{}{}, nil
}
