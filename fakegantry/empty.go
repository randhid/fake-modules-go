package fakegantry

import (
	"context"

	"go.viam.com/rdk/components/gantry"
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

func newEmptyGantry(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	gantry.Gantry, error,
) {
	e := &empty{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	return e, nil
}

func (e *empty) Position(ctx context.Context, extra map[string]interface{}) ([]float64, error) {
	return nil, nil
}

func (e *empty) Lengths(ctx context.Context, extra map[string]interface{}) ([]float64, error) {
	return nil, nil
}

func (e *empty) Home(ctx context.Context, extra map[string]interface{}) (bool, error) {
	return true, nil
}

func (e *empty) MoveToPosition(ctx context.Context, target []float64, speeds []float64, extra map[string]interface{}) error {
	return nil
}

func (e *empty) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return nil, nil
}

func (e *empty) IsMoving(context.Context) (bool, error) {
	return false, nil
}

func (e *empty) Stop(context.Context, map[string]interface{}) error {
	return nil
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
