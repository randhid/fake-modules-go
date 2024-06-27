package fakepowersensor

import (
	"context"
	"math"

	"go.viam.com/rdk/components/powersensor"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type errorring struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger
}

func newErroringPowerSensor(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	powersensor.PowerSensor, error,
) {
	f := &errorring{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return f, nil
}

func (e *errorring) Voltage(context.Context, map[string]interface{}) (float64, bool, error) {
	return math.NaN(), false, grpc.UnimplementedError
}

func (e *errorring) Current(context.Context, map[string]interface{}) (float64, bool, error) {
	return math.NaN(), false, grpc.UnimplementedError
}

func (e *errorring) Power(context.Context, map[string]interface{}) (float64, error) {
	return math.NaN(), grpc.UnimplementedError
}

func (e *errorring) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
	return nil, grpc.UnimplementedError
}
