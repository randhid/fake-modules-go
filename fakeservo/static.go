package fakeservo

import (
	"context"
	"sync"

	"go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type static struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	logger   logging.Logger
	mu       sync.Mutex
	position uint32
}

func newStaticServo(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	servo.Servo, error,
) {
	s := &static{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return s, nil
}

func (s *static) Move(ctx context.Context, pos uint32, extra map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.position = pos
	return nil
}

func (s *static) Position(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.position, nil
}

func (s *static) Stop(context.Context, map[string]interface{}) error {
	return nil
}

func (s *static) IsMoving(context.Context) (bool, error) {
	return false, nil
}
