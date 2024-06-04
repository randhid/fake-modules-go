package fakesensor

import (
	"context"
	"fake-modules-go/common"
	"math/rand"
	"sync"
	"time"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type WaitingConfig struct {
	resource.TriviallyValidateConfig
	WaitTime int `json:"wait_time_milli_seconds,omitempty"`
}

type waiting struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	logger logging.Logger

	mu       sync.Mutex
	waitTime time.Duration
}

func newWaitingSensor(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	sensor.Sensor, error,
) {

	newConf, err := resource.NativeConfig[WaitingConfig](conf)
	if err != nil {
		return nil, err
	}

	waitTime := common.DefaultWaitTimeMs // stes default to 500 milliseconds
	if newConf.WaitTime > 0 {            // otherwise set the wait time to the user-configured value
		waitTime = time.Duration(newConf.WaitTime) * time.Millisecond
	}

	w := &waiting{
		Named:    conf.ResourceName().AsNamed(),
		logger:   logger,
		waitTime: waitTime,
	}

	return w, nil
}

func (w *waiting) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	randomBool := rand.Int() % 2
	randomfloat := rand.Float64()

	time.Sleep(w.waitTime)

	return map[string]interface{}{
		"on":     randomBool,
		"value:": randomfloat}, nil

}
