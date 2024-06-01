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
	resource.Sensor
}

func newEmptyPowerSensor(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	powersensor.PowerSensor, error,
) {
	f := &static{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return f, nil
}

func (e *empty) Voltage(context.Context, map[string]interface{}) (float64, bool, error) {
	return math.NaN(), false, nil
}

func (e *empty) Current(context.Context, map[string]interface{}) (float64, bool, error) {
	return math.NaN(), false, nil
}

func (e *empty) Power(ctx context.Context, extra map[string]interface{}) (float64, error) {
	return math.NaN(), nil
}
