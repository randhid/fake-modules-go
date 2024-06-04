package fakepowersensor

import (
	"context"
	"errors"
	"fake-modules-go/common"
	"math"
	"math/rand"
	"sync"
	"time"

	"go.viam.com/rdk/components/powersensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type WaitingConfig struct {
	resource.TriviallyValidateConfig
	IsAC     bool `json:"make_ac_sensor,omitempty"`
	WaitTime int  `json:"wait_time_milli_seconds,omitempty"`
}

type waiting struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger

	mu       sync.Mutex
	ac       bool
	waitTime time.Duration
}

func newWaitingPowerSensor(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	powersensor.PowerSensor, error,
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

func (w *waiting) Voltage(ctx context.Context, extra map[string]interface{}) (float64, bool, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	randomfloat := rand.Float64()
	time.Sleep(w.waitTime)
	return randomfloat, w.ac, nil
}

func (w *waiting) Current(ctx context.Context, extra map[string]interface{}) (float64, bool, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	randomfloat := rand.Float64()
	time.Sleep(w.waitTime)
	return randomfloat, w.ac, nil
}

func (w *waiting) Power(ctx context.Context, extra map[string]interface{}) (float64, error) {
	v, _, vErr := w.Voltage(ctx, extra)
	i, _, iErr := w.Current(ctx, extra)

	err := errors.Join(vErr, iErr)
	if err != nil {
		return math.NaN(), err
	}

	time.Sleep(w.waitTime)
	return v * i, nil
}

func (w *waiting) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	v, vac, vErr := w.Voltage(ctx, extra)
	i, iac, iErr := w.Current(ctx, extra)
	p, pErr := w.Power(ctx, extra)

	err := errors.Join(vErr, iErr, pErr)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{
		"voltage_volts": v, "voltage_ac": vac, "current_amperes": i, "current_ac": iac, "power_watts": p,
	}

	time.Sleep(w.waitTime)
	return res, nil

}
