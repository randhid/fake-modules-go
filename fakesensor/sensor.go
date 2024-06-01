package fakesensor

import (
	"context"
	"math/rand"
	"sync"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type fake struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	logger logging.Logger

	mu sync.Mutex
}

func newFakeSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	sensor.Sensor, error,
) {
	f := &fake{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return f, nil
}

func (f *fake) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	randomBool := rand.Int() % 2
	randomfloat := rand.Float64()

	return map[string]interface{}{
		"on":     randomBool,
		"value:": randomfloat}, nil

}
