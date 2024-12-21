package fakearm

import (
	// for embedding model file.
	"context"
	_ "embed"
	"math"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/components/arm"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

type nanArm struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable

	logger logging.Logger
}

var nan = math.NaN()

func newNanArm(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	arm.Arm, error,
) {

	n := &nanArm{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	return n, nil
}

func (n *nanArm) EndPosition(ctx context.Context, extra map[string]interface{}) (spatialmath.Pose, error) {
	return spatialmath.NewPose(
		r3.Vector{X: nan, Y: nan, Z: nan},
		&spatialmath.OrientationVector{OX: nan, OY: nan, OZ: nan, Theta: nan}), nil
}

func (n *nanArm) MoveToPosition(ctx context.Context, pos spatialmath.Pose, extra map[string]interface{}) error {
	return nil
}

// MoveToJointPositions sets the joints.
func (n *nanArm) MoveToJointPositions(ctx context.Context, joints []referenceframe.Input, extra map[string]interface{}) error {
	return nil
}

// JointPositions returns joints.
func (n *nanArm) JointPositions(ctx context.Context, extra map[string]interface{}) ([]referenceframe.Input, error) {
	return referenceframe.FloatsToInputs([]float64{nan, nan, nan, nan}), nil
}

func (n *nanArm) MoveThroughJointPositions(context.Context, [][]referenceframe.Input, *arm.MoveOptions, map[string]any) error {
	return nil
}

func (n *nanArm) IsMoving(context.Context) (bool, error) {
	return false, nil
}

func (n *nanArm) Stop(context.Context, map[string]interface{}) error {
	return nil
}

func (n *nanArm) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	box, err := spatialmath.NewBox(
		spatialmath.NewPose(r3.Vector{X: nan, Y: nan, Z: nan}, &spatialmath.OrientationVector{OX: nan, OY: nan, OZ: nan, Theta: nan}),
		r3.Vector{X: nan, Y: nan, Z: nan},
		"box",
	)
	return []spatialmath.Geometry{box}, err
}

func (n *nanArm) CurrentInputs(ctx context.Context) ([]referenceframe.Input, error) {
	return referenceframe.FloatsToInputs([]float64{nan, nan, nan, nan}), nil
}
func (n *nanArm) GoToInputs(context.Context, ...[]referenceframe.Input) error {
	return nil
}

func (n *nanArm) ModelFrame() referenceframe.Model {
	return nil
}
