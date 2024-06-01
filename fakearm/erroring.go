package fakearm

import (
	// for embedding model file.
	"context"
	_ "embed"

	pb "go.viam.com/api/component/arm/v1"
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

func (e *erroring) EndPosition(ctx context.Context, extra map[string]interface{}) (spatialmath.Pose, error) {
	return nil, grpc.UnimplementedError
}

func (e *erroring) MoveToPosition(ctx context.Context, pos spatialmath.Pose, extra map[string]interface{}) error {
	return grpc.UnimplementedError
}

// MoveToJointPositions sets the joints.
func (e *erroring) MoveToJointPositions(ctx context.Context, joints *pb.JointPositions, extra map[string]interface{}) error {
	return nil
}

// JointPositions returns joints.
func (e *erroring) JointPositions(ctx context.Context, extra map[string]interface{}) (*pb.JointPositions, error) {
	return nil, grpc.UnimplementedError
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
