package fakepowersensor

import (
	"context"
	"math"

	"go.viam.com/rdk/components/powersensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type empty struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger
}

func newEmptyPowerSensor(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	powersensor.PowerSensor, error,
) {
	e := &empty{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return e, nil
}

func (e *empty) Voltage(context.Context, map[string]interface{}) (float64, bool, error) {
	return math.NaN(), false, nil
}

func (e *empty) Current(context.Context, map[string]interface{}) (float64, bool, error) {
	return math.NaN(), false, nil
}

func (e *empty) Power(context.Context, map[string]interface{}) (float64, error) {
	return math.NaN(), nil
}

func (e *empty) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
