package fakeservo

import (
	"context"
	"sync"

	"go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

const increment = 1

type fake struct {
	resource.Named
	resource.AlwaysRebuild

	logger logging.Logger

	mu         sync.Mutex
	moving     bool
	position   uint32
	stopmoving func()
}

func newFakeServo(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	servo.Servo, error,
) {
	f := &fake{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return f, nil
}

func (f *fake) Move(ctx context.Context, pos uint32, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	currPos := f.position
	travel := currPos - pos
	f.moving = true

	var moveCtx context.Context
	moveCtx, f.stopmoving = context.WithCancel(context.Background())
	for travel != 0 {
		if moveCtx.Err() != nil {
			f.moving = false
			f.position = currPos + travel
			break
		}

		switch {
		case travel > 0:
			f.position += increment
			travel -= increment
		case travel < 0:
			f.position -= increment
			travel += increment
		}
	}
	f.position = pos
	return nil
}

func (f *fake) Position(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.position, nil
}

func (f *fake) IsMoving(ctx context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.moving, nil
}

func (f *fake) Stop(ctx context.Context, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.moving = false

	if f.stopmoving != nil {
		f.stopmoving()
	}

	var err error
	f.position, err = f.Position(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (f *fake) Close(ctx context.Context) error {
	if err := f.Stop(ctx, nil); err != nil {
		return err
	}
	return nil
}
