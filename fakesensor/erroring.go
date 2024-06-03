package fakesensor

import (
	"context"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type erroring struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger
}

func newErroringSensor(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	sensor.Sensor, error,
) {
	e := &erroring{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return e, nil
}

func (e *erroring) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
