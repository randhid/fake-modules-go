package fakemovementsensor

import (
	"context"
	"math"

	"github.com/golang/geo/r3"
	geo "github.com/kellydunn/golang-geo"
	"go.viam.com/rdk/components/movementsensor"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

type erroring struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	logger logging.Logger
}

func newErroringMovementSensor(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	movementsensor.MovementSensor, error,
) {
	e := &erroring{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	return e, nil
}

func (e *erroring) Position(ctx context.Context, extra map[string]interface{}) (*geo.Point, float64, error) {
	return geo.NewPoint(math.NaN(), math.NaN()), math.NaN(), grpc.UnimplementedError
}

func (e *erroring) Orientation(ctx context.Context, extra map[string]interface{}) (spatialmath.Orientation, error) {
	return &spatialmath.OrientationVector{
		OX:    math.NaN(),
		OY:    math.NaN(),
		OZ:    math.NaN(),
		Theta: math.NaN()}, grpc.UnimplementedError
}

func (e *erroring) CompassHeading(ctx context.Context, extra map[string]interface{}) (float64, error) {
	return math.NaN(), grpc.UnimplementedError

}

func (e *erroring) AngularVelocity(ctx context.Context, extra map[string]interface{}) (spatialmath.AngularVelocity, error) {
	return spatialmath.AngularVelocity{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, grpc.UnimplementedError
}

func (e *erroring) LinearVelocity(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	return r3.Vector{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, grpc.UnimplementedError
}

func (e *erroring) LinearAcceleration(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	return r3.Vector{}, grpc.UnimplementedError
}

func (e *erroring) Accuracy(ctx context.Context, extra map[string]interface{}) (*movementsensor.Accuracy, error) {
	accmap := map[string]float32{
		"foo": float32(math.NaN()),
		"bar": float32(math.NaN()),
	}
	return &movementsensor.Accuracy{
		AccuracyMap:        accmap,
		Hdop:               float32(math.NaN()),
		Vdop:               float32(math.NaN()),
		CompassDegreeError: float32(math.NaN()),
	}, grpc.UnimplementedError
}

func (e *erroring) Properties(ctx context.Context, extra map[string]interface{}) (*movementsensor.Properties, error) {
	return &movementsensor.Properties{}, grpc.UnimplementedError
}

func (e *erroring) Readings(context.Context, map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"foo": math.NaN(),
	}, grpc.UnimplementedError
}
