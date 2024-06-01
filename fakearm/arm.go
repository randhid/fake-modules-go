package fakearm

import (
	"context"
	"fake-modules-go/common"
	"sync"
	"time"

	pb "go.viam.com/api/component/arm/v1"
	"go.viam.com/rdk/components/arm"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
	"go.viam.com/rdk/utils"
)

type fake struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	logger logging.Logger

	mu          sync.Mutex
	moving      bool
	stopmoving  func()
	model       referenceframe.Model
	jointValues []float64
}

func newFakeArm(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	arm.Arm, error,
) {
	model, err := makeModelFrame(conf.Name)
	if err != nil {
		return nil, err
	}

	dof := len(model.DoF())
	jointValues := make([]float64, dof)

	f := &fake{
		Named:       conf.ResourceName().AsNamed(),
		logger:      logger,
		model:       model,
		jointValues: jointValues,
	}

	return f, nil
}

func (f *fake) MoveToJointPositions(ctx context.Context, joints *pb.JointPositions, extra map[string]interface{}) error {
	// Extract the target joint positions
	targetJoints := joints.Values

	// Lock the mutex to ensure thread safety
	f.mu.Lock()
	defer f.mu.Unlock()

	var moveCtx context.Context
	moveCtx, f.stopmoving = context.WithCancel(context.Background())
	currPos := append([]float64(nil), f.jointValues...)
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
				allReached := true
				// Loop through the lengths and update positions
				for i := range f.model.DoF() {
					increment := (targetJoints[i] - currPos[i]) * 0.1
					if !utils.Float64AlmostEqual(targetJoints[i], currPos[i], common.GoalWithinRange) {
						f.moving = true
						currPos[i] += increment
						f.jointValues = currPos
						allReached = false
					}
				}
				if allReached {
					// All positions are within the goal range
					f.moving = false
					return
				}
			}
			f.logger.Debugf("joint positions %v", f.jointValues)
		}
	}()

	return nil
}

func (f *fake) JointPositions(ctx context.Context, extra map[string]interface{}) (*pb.JointPositions, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	retJoint := &pb.JointPositions{Values: f.jointValues}
	return retJoint, nil
}

func (f *fake) Stop(ctx context.Context, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.stopmoving != nil {
		f.stopmoving()
	}
	f.moving = false
	return nil
}

func (f *fake) IsMoving(ctx context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.moving, nil
}

func (f *fake) ModelFrame() referenceframe.Model {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.model
}

func (f *fake) Geometries(ctx context.Context, extra map[string]interface{}) ([]spatialmath.Geometry, error) {
	res, err := f.JointPositions(ctx, nil)
	if err != nil {
		return nil, err
	}
	inputs := f.model.InputFromProtobuf(res)

	gif, err := f.model.Geometries(inputs)
	if err != nil {
		return nil, err
	}
	return gif.Geometries(), nil
}

func (f *fake) EndPosition(ctx context.Context, extra map[string]interface{}) (spatialmath.Pose, error) {
	return spatialmath.NewZeroPose(), grpc.UnimplementedError
}

func (f *fake) MoveToPosition( ctx context.Context, pos spatialmath.Pose, extra map[string]interface{}) error {
	return grpc.UnimplementedError
}

func (f *fake) CurrentInputs(ctx context.Context) ([]referenceframe.Input, error) {
	pos, err := f.JointPositions(ctx, nil)
	if err != nil {
		return []referenceframe.Input{}, err
	}
	inputs := referenceframe.FloatsToInputs(pos.Values)
	return inputs, nil
}

func (f *fake) GoToInputs(ctx context.Context, inputs ...[]referenceframe.Input) error {
	for _, input := range inputs {
		pos := referenceframe.InputsToFloats(input)
		jp := pb.JointPositions{Values: pos}
		if err := f.MoveToJointPositions(ctx, &jp, nil); err != nil {
			return err
		}
	}
	return nil
}
