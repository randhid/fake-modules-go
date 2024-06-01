package fakemovementsensor

import (
	"context"
	"math/rand"

	"github.com/golang/geo/r3"
	geo "github.com/kellydunn/golang-geo"
	"go.viam.com/rdk/components/movementsensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

type static struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	logger logging.Logger
}

func newStaticMovementSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	movementsensor.MovementSensor, error,
) {
	f := &static{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	return f, nil
}

func (s *static) Position(ctx context.Context, extra map[string]interface{}) (*geo.Point, float64, error) {
	return geo.NewPoint(73.0, 72.0), 64, nil
}

func (s *static) Orientation(ctx context.Context, extra map[string]interface{}) (spatialmath.Orientation, error) {
	return &spatialmath.OrientationVector{OX: 1, OY: 1, OZ: 0, Theta: 25}, nil
}

func (s *static) CompassHeading(ctx context.Context, extra map[string]interface{}) (float64, error) {
	return 355, nil

}

func (s *static) AngularVelocity(ctx context.Context, extra map[string]interface{}) (spatialmath.AngularVelocity, error) {

	return spatialmath.AngularVelocity{X: 7, Y: 8, Z: 9}, nil
}

func (s *static) LinearVelocity(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	return r3.Vector{X: 1, Y: 2, Z: 3}, nil
}

func (s *static) LinearAcceleration(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	return r3.Vector{X: 4, Y: 5, Z: 9.81}, nil
}

func (s *static) Accuracy(ctx context.Context, extra map[string]interface{}) (*movementsensor.Accuracy, error) {
	accmap := map[string]float32{
		"satellites_noise_signal": 11,
		"rand_trust_level":        12,
	}
	return &movementsensor.Accuracy{
		AccuracyMap:        accmap,
		Hdop:               0.8,
		Vdop:               0.7,
		CompassDegreeError: 0.5,
		NmeaFix:            4,
	}, nil
}

func (s *static) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	defaults, err := movementsensor.DefaultAPIReadings(ctx, s, extra)
	if err != nil {
		return nil, err
	}
	defaults["foo"] = "bar"
	defaults["satellites"] = rand.Intn(32)
	return defaults, nil
}

func (s *static) Properties(ctx context.Context, extra map[string]interface{}) (*movementsensor.Properties, error) {
	return &movementsensor.Properties{
		PositionSupported:           true,
		OrientationSupported:        true,
		CompassHeadingSupported:     true,
		LinearVelocitySupported:     true,
		AngularVelocitySupported:    true,
		LinearAccelerationSupported: true,
	}, nil
}
