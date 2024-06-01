package fakepowersensor

import (
	"context"
	"errors"
	"math/rand"
	"sync"

	"go.viam.com/rdk/components/powersensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type StaticConfig struct {
	resource.TriviallyValidateConfig
	IsAC bool `json:"make_ac_sensor,omitempty"`
}

type static struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger logging.Logger

	mu sync.Mutex
	ac bool
}

func newStaticPowerSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	powersensor.PowerSensor, error,
) {

	newConf, err := resource.NativeConfig[*StaticConfig](conf)
	if err != nil {
		return nil, err
	}

	f := &static{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
		ac:     newConf.IsAC,
	}

	return f, nil
}

func (s *static) Voltage(ctx context.Context, extra map[string]interface{}) (float64, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	randomfloat := rand.Float64()
	return randomfloat, s.ac, nil
}

func (s *static) Current(ctx context.Context, extra map[string]interface{}) (float64, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return 1.0, s.ac, nil
}

func (s *static) Power(ctx context.Context, extra map[string]interface{}) (float64, error) {
	v, _, _ := s.Voltage(ctx, extra)
	i, _, _ := s.Current(ctx, extra)

	return v * i, nil
}

func (s *static) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	v, vac, vErr := s.Voltage(ctx, extra)
	i, iac, iErr := s.Current(ctx, extra)
	p, pErr := s.Power(ctx, extra)

	err := errors.Join(vErr, iErr, pErr)
	if err != nil {
		return nil, err
	}
	res := map[string]interface{}{
		"voltage_volts": v, "voltage_ac": vac, "current_amperes": i, "current_ac": iac, "power_watts": p,
	}

	return res, nil

}
