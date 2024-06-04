package fakemotor

import (
	"context"
	"math"

	"go.viam.com/rdk/components/motor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type empty struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	logger logging.Logger

	// nil interfaces
	resource.Actuator
}

func newEmptyMotor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	motor.Motor, error,
) {
	e := &empty{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	return e, nil
}

func (e *empty) SetPower(context.Context, float64, map[string]interface{}) error {
	return nil
}

func (e *empty) GoFor(context.Context, float64, float64, map[string]interface{}) error {
	return nil
}

func (e *empty) GoTo(context.Context, float64, float64, map[string]interface{}) error {
	return nil
}

func (e *empty) Position(context.Context, map[string]interface{}) (float64, error) {
	return math.NaN(), nil
}

func (e *empty) ResetZeroPosition(context.Context, float64, map[string]interface{}) error {
	return nil
}

func (e *empty) IsPowered(context.Context, map[string]interface{}) (bool, float64, error) {
	return false, math.NaN(), nil
}

func (e *empty) Properties(context.Context, map[string]interface{}) (motor.Properties, error) {
	return motor.Properties{}, nil
}
