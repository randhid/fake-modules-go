package fakesensor

import (
	"context"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type static struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	logger logging.Logger
}

func newStaticSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	sensor.Sensor, error,
) {
	s := &static{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return s, nil
}

func (s *static) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"on":     true,
		"value:": 2.14159}, nil

}
