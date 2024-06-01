package fakegripper

import (
	"context"
	"fake-modules-go/common"
	"sync"
	"time"

	"go.viam.com/rdk/components/gripper"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

const (
	servoName = "fake-gripper"
	increment = 1
)

var Model = common.FakesFamily.WithModel(servoName)

func init() {
	resource.RegisterComponent(gripper.API, Model, resource.Registration[gripper.Gripper, resource.NoNativeConfig]{
		Constructor: newFakeGripper,
	})
}

type fake struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	logger logging.Logger

	mu         sync.Mutex
	moving     bool
	position   uint32
	stopmoving func()
}

func newFakeGripper(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	gripper.Gripper, error,
) {
	f := &fake{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return f, nil
}

func (f *fake) simulateMove(ctx context.Context, extra map[string]interface{}) {
	var moveCtx context.Context
	moveCtx, f.stopmoving = context.WithCancel(context.Background())

	// simulate motion
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100) // Update every 100ms
		defer ticker.Stop()

		for {
			select {
			case <-moveCtx.Done():
				// Exit the goroutine if the context is canceled
				return
			case <-ticker.C:
				// Check if the context has been canceled
				if moveCtx.Err() != nil {
					err := f.Stop(ctx, extra)
					f.logger.Error(err)
					return // Ensure we exit the goroutine
				}

				// Loop through a simple busy loop
				for i := 0; i < 50; i++ {
					f.moving = true
				}
			}
		}
	}()
}

func (f *fake) Grab(ctx context.Context, extra map[string]interface{}) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.simulateMove(ctx, extra)
	// we're done moving
	f.moving = false
	grabbed := true
	return grabbed, nil
}

func (f *fake) Open(ctx context.Context, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.simulateMove(ctx, extra)
	// we're done moving
	f.moving = false

	return nil
}

func (f *fake) Stop(context.Context, map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.stopmoving != nil {
		f.stopmoving()
	}
	return nil
}

func (f *fake) IsMoving(context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.moving, nil
}

func (f *fake) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return []spatialmath.Geometry{}, nil
}

func (f *fake) ModelFrame() referenceframe.Model {
	return nil
}
