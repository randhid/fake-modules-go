package fakegripper

import (
	"context"
	"errors"
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

func (f *fake) simulateMove(ctx context.Context) error {
	var moveCtx context.Context
	{
		f.mu.Lock()
		moveCtx, f.stopmoving = context.WithTimeout(ctx, 2*time.Second)
		f.moving = true
		f.mu.Unlock()
	}
	<-moveCtx.Done()
	{
		f.mu.Lock()
		f.moving = false
		f.mu.Unlock()
	}
	if errors.Is(moveCtx.Err(), context.DeadlineExceeded) {
		return nil
	}
	return moveCtx.Err()
}

func (f *fake) Grab(ctx context.Context, extra map[string]interface{}) (bool, error) {
	return true, f.simulateMove(ctx)
}

func (f *fake) Open(ctx context.Context, extra map[string]interface{}) error {
	return f.simulateMove(ctx)
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
