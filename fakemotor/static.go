package fakemotor

import (
	"context"
	"sync"

	"go.viam.com/rdk/components/motor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type static struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	logger logging.Logger

	mu       sync.Mutex
	power    float64
	position float64
	moving   bool
}

func newStaticMotor(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	motor.Motor, error,
) {
	s := &static{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	return s, nil
}

func (s *static) SetPower(ctx context.Context, power float64, extra map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.power = power
	s.moving = true
	return nil
}

func (s *static) GoFor(ctx context.Context, rpm, revolutions float64, extra map[string]interface{}) error {
	if rpm == 0 {
		return motor.NewZeroRPMError()
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.position += revolutions
	return nil
}

func (s *static) GoTo(ctx context.Context, rpm, targetPos float64, extra map[string]interface{}) error {
	pos, err := s.Position(ctx, extra)
	if err != nil {
		return err
	}
	rev := targetPos - pos

	return s.GoFor(ctx, rev, rpm, extra)
}

func (s *static) Position(ctx context.Context, extra map[string]interface{}) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.position, nil
}

func (s *static) ResetZeroPosition(ctx context.Context, offset float64, extra map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.position = 0
	return nil
}

func (s *static) Properties(ctx context.Context, extra map[string]interface{}) (motor.Properties, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return motor.Properties{PositionReporting: true}, nil
}

func (s *static) IsMoving(ctx context.Context) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.moving, nil
}

func (s *static) IsPowered(ctx context.Context, extra map[string]interface{}) (bool, float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.moving, s.power, nil
}

func (s *static) Stop(ctx context.Context, extra map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.moving = false
	s.power = 0
	return nil
}
