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

type empty struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable

	logger logging.Logger
}

func newEmptyArm(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	arm.Arm, error,
) {

	e := &empty{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return e, nil
}

func (e *empty) EndPosition(ctx context.Context, extra map[string]interface{}) (spatialmath.Pose, error) {
	return nil, grpc.UnimplementedError
}

func (e *empty) MoveToPosition(ctx context.Context, pos spatialmath.Pose, extra map[string]interface{}) error {
	return nil
}

// MoveToJointPositions sets the joints.
func (e *empty) MoveToJointPositions(ctx context.Context, joints []referenceframe.Input, extra map[string]interface{}) error {
	return nil
}

// JointPositions returns joints.
func (e *empty) JointPositions(ctx context.Context, extra map[string]interface{}) ([]referenceframe.Input, error) {
	return nil, nil
}

func (e *empty) MoveThroughJointPositions(context.Context, [][]referenceframe.Input, *arm.MoveOptions, map[string]any) error {
	return nil
}

func (e *empty) IsMoving(context.Context) (bool, error) {
	return false, nil
}

func (e *empty) Stop(context.Context, map[string]interface{}) error {
	return nil
}

func (e *empty) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return nil, nil
}

func (e *empty) CurrentInputs(ctx context.Context) ([]referenceframe.Input, error) {
	return nil, nil
}
func (e *empty) GoToInputs(context.Context, ...[]referenceframe.Input) error {
	return nil
}

func (e *empty) ModelFrame() referenceframe.Model {
	return nil
}
