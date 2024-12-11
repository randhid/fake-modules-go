package fakearm

import (
	// for embedding model file.
	"context"
	_ "embed"

	"go.viam.com/rdk/components/arm"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

type erroring struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable

	logger logging.Logger
}

func newErroringArm(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	arm.Arm, error,
) {

	e := &erroring{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return e, nil
}

func (e *erroring) EndPosition(context.Context, map[string]interface{}) (spatialmath.Pose, error) {
	return nil, grpc.UnimplementedError
}

func (e *erroring) MoveToPosition(context.Context, spatialmath.Pose, map[string]interface{}) error {
	return grpc.UnimplementedError
}

// MoveToJointPositions sets the joints.
func (e *erroring) MoveToJointPositions(context.Context, []referenceframe.Input, map[string]interface{}) error {
	return grpc.UnimplementedError
}

// JointPositions returns joints.
func (e *erroring) JointPositions(ctx context.Context, extra map[string]interface{}) ([]referenceframe.Input, error) {
	return nil, grpc.UnimplementedError
}

func (e *erroring) MoveThroughJointPositions(context.Context, [][]referenceframe.Input, *arm.MoveOptions, map[string]interface{}) error {
	return grpc.UnimplementedError
}

func (e *erroring) IsMoving(context.Context) (bool, error) {
	return false, grpc.UnimplementedError
}

func (e *erroring) Stop(context.Context, map[string]interface{}) error {
	return grpc.UnimplementedError
}

func (e *erroring) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return nil, grpc.UnimplementedError
}

func (e *erroring) CurrentInputs(ctx context.Context) ([]referenceframe.Input, error) {
	return nil, grpc.UnimplementedError
}
func (e *erroring) GoToInputs(context.Context, ...[]referenceframe.Input) error {
	return grpc.UnimplementedError
}

func (e *erroring) ModelFrame() referenceframe.Model {
	return nil
}
