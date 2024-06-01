package fakesensor

import (
	"context"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type empty struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	resource.Sensor

	logger logging.Logger
}

func newEmptySensor(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	sensor.Sensor, error,
) {
	f := &empty{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return f, nil
}
