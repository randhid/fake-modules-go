package fakegantry

import (
	"context"

	"go.viam.com/rdk/components/gantry"
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

func newErroringGantry(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	gantry.Gantry, error,
) {
	e := &erroring{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	return e, nil
}

func (e *erroring) Position(ctx context.Context, extra map[string]interface{}) ([]float64, error) {
	return nil, grpc.UnimplementedError
}

func (e *erroring) Lengths(ctx context.Context, extra map[string]interface{}) ([]float64, error) {
	return nil, grpc.UnimplementedError
}

func (e *erroring) Home(ctx context.Context, extra map[string]interface{}) (bool, error) {
	return true, nil
}

func (e *erroring) MoveToPosition(ctx context.Context, target []float64, speeds []float64, extra map[string]interface{}) error {
	return nil
}

func (e *erroring) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return nil, grpc.UnimplementedError
}

func (e *erroring) IsMoving(context.Context) (bool, error) {
	return false, grpc.UnimplementedError
}

func (e *erroring) Stop(context.Context, map[string]interface{}) error {
	return nil
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
