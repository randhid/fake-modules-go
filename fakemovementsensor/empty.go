package fakemovementsensor

import (
	"context"
	"math"

	"github.com/golang/geo/r3"
	geo "github.com/kellydunn/golang-geo"
	"go.viam.com/rdk/components/movementsensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

type empty struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	logger logging.Logger
}

func newEmptyMovementSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	movementsensor.MovementSensor, error,
) {
	f := &empty{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	return f, nil
}

func (e *empty) Position(ctx context.Context, extra map[string]interface{}) (*geo.Point, float64, error) {
	return nil, math.NaN(), nil
}

func (e *empty) Orientation(ctx context.Context, extra map[string]interface{}) (spatialmath.Orientation, error) {
	return nil, nil
}

func (e *empty) CompassHeading(ctx context.Context, extra map[string]interface{}) (float64, error) {
	return math.NaN(), nil

}

func (e *empty) AngularVelocity(ctx context.Context, extra map[string]interface{}) (spatialmath.AngularVelocity, error) {
	return spatialmath.AngularVelocity{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, nil
}

func (e *empty) LinearVelocity(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	return r3.Vector{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, nil
}

func (e *empty) LinearAcceleration(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	return r3.Vector{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, nil
}

func (e *empty) Accuracy(ctx context.Context, extra map[string]interface{}) (*movementsensor.Accuracy, error) {
	return nil, nil
}

func (e *empty) Properties(ctx context.Context, extra map[string]interface{}) (*movementsensor.Properties, error) {
	return &movementsensor.Properties{
		PositionSupported:           true,
		OrientationSupported:        true,
		CompassHeadingSupported:     true,
		LinearVelocitySupported:     true,
		AngularVelocitySupported:    true,
		LinearAccelerationSupported: true,
	}, nil
}

func (e *empty) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
