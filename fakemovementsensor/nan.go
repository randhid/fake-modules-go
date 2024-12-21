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

type nanMS struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	logger logging.Logger
}

func newNanMovementSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	movementsensor.MovementSensor, error,
) {
	f := &empty{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	return f, nil
}

func (n *nanMS) Position(ctx context.Context, extra map[string]interface{}) (*geo.Point, float64, error) {
	return nil, math.NaN(), nil
}

func (n *nanMS) Orientation(ctx context.Context, extra map[string]interface{}) (spatialmath.Orientation, error) {
	return &spatialmath.OrientationVector{OX: math.NaN(), OY: math.NaN(), OZ: math.NaN(), Theta: math.NaN()}, nil
}

func (n *nanMS) CompassHeading(ctx context.Context, extra map[string]interface{}) (float64, error) {
	return math.NaN(), nil

}

func (n *nanMS) AngularVelocity(ctx context.Context, extra map[string]interface{}) (spatialmath.AngularVelocity, error) {
	return spatialmath.AngularVelocity{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, nil
}

func (n *nanMS) LinearVelocity(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	return r3.Vector{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, nil
}

func (n *nanMS) LinearAcceleration(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	return r3.Vector{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, nil
}

func (n *nanMS) Accuracy(ctx context.Context, extra map[string]interface{}) (*movementsensor.Accuracy, error) {
	return &movementsensor.Accuracy{
		Hdop:               float32(math.NaN()),
		Vdop:               float32(math.NaN()),
		CompassDegreeError: float32(math.NaN()),
		NmeaFix:            -1,
		AccuracyMap: map[string]float32{
			"one": float32(math.NaN()),
		},
	}, nil
}

func (n *nanMS) Properties(ctx context.Context, extra map[string]interface{}) (*movementsensor.Properties, error) {
	return &movementsensor.Properties{
		PositionSupported:           true,
		OrientationSupported:        true,
		CompassHeadingSupported:     true,
		LinearVelocitySupported:     true,
		AngularVelocitySupported:    true,
		LinearAccelerationSupported: true,
	}, nil
}

func (n *nanMS) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}
