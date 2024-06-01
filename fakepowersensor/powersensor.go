package fakepowersensor

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"sync"

	"go.viam.com/rdk/components/powersensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type Config struct {
	resource.TriviallyValidateConfig
	IsAC bool `json:"make_ac_sensor,omitempty"`
}

type fake struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger

	mu sync.Mutex
	ac bool
}

func newFakePowerSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	powersensor.PowerSensor, error,
) {
	f := &fake{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return f, nil
}

func (f *fake) Voltage(ctx context.Context, extra map[string]interface{}) (float64, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	randomfloat := rand.Float64()
	return randomfloat, f.ac, nil
}

func (f *fake) Current(ctx context.Context, extra map[string]interface{}) (float64, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	randomfloat := rand.Float64()
	return randomfloat, f.ac, nil
}

func (f *fake) Power(ctx context.Context, extra map[string]interface{}) (float64, error) {
	v, _, vErr := f.Voltage(ctx, extra)
	i, _, iErr := f.Current(ctx, extra)

	err := errors.Join(vErr, iErr)
	if err != nil {
		return math.NaN(), err
	}

	return v * i, nil
}

func (f *fake) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	v, vac, vErr := f.Voltage(ctx, extra)
	i, iac, iErr := f.Current(ctx, extra)
	p, pErr := f.Power(ctx, extra)

	err := errors.Join(vErr, iErr, pErr)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{
		"voltage_volts": v, "voltage_ac": vac, "current_amperes": i, "current_ac": iac, "power_watts": p,
	}

	return res, nil

}
