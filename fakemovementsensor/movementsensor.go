package fakemovementsensor

import (
	"context"
	"math"
	"math/rand"
	"sync"

	"fake-modules-go/common"

	"github.com/golang/geo/r3"
	geo "github.com/kellydunn/golang-geo"
	"go.viam.com/rdk/components/movementsensor"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)


type Config struct {
	resource.TriviallyValidateConfig
	PosOff    bool `json:"turn_off_position,omitempty"`
	OriOff    bool `json:"turn_off_orientation,omitempty"`
	CompOff   bool `json:"turn_off_compass_heading,omitempty"`
	AngVelOff bool `json:"turn_off_angular_velcoity,omitempty"`
	LinVelOff bool `json:"turn_off_linear_velocity,omitempty"`
	LinAccOff bool `json:"turn_off_linear_acceleration,omitempty"`
}

type fake struct {
	resource.Named
	resource.TriviallyCloseable

	logger logging.Logger

	mu    sync.Mutex
	coord geo.Point

	posOff     bool
	oriOff     bool
	compassOff bool
	linvelOff  bool
	angvelOff  bool
	linaccOff  bool
}

func newFakeMovementSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	movementsensor.MovementSensor, error,
) {
	f := &fake{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	if err := f.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *fake) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	newConf, err := resource.NativeConfig[Config](conf)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.posOff = newConf.PosOff
	f.oriOff = newConf.OriOff
	f.compassOff = newConf.CompOff
	f.linvelOff = newConf.LinVelOff
	f.angvelOff = newConf.AngVelOff
	f.linaccOff = newConf.LinAccOff

	return nil
}

func (f *fake) Position(ctx context.Context, extra map[string]interface{}) (*geo.Point, float64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.posOff {
		return nil, math.NaN(), grpc.UnimplementedError
	}

	latincrement := rand.Float64() / 10e4 * common.Randomsign()
	lngincrement := rand.Float64() / 10e5 * math.Cos(latincrement) * common.Randomsign()
	altitude := rand.Float64() * 100

	newcoord := geo.NewPoint(
		f.coord.Lat()+latincrement,
		f.coord.Lng()+lngincrement,
	)
	return newcoord, altitude, nil
}

func (f *fake) Orientation(ctx context.Context, extra map[string]interface{}) (spatialmath.Orientation, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.oriOff {
		return nil, grpc.UnimplementedError
	}

	ox := rand.Float64() * common.Randomsign()
	oy := rand.Float64() * common.Randomsign()
	oz := rand.Float64() * common.Randomsign()
	theta := rand.Float64()

	ov := spatialmath.OrientationVector{OX: ox, OY: oy, OZ: oz, Theta: theta}
	ov.Normalize()
	return &ov, nil
}

func (f *fake) CompassHeading(ctx context.Context, extra map[string]interface{}) (float64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.compassOff {
		return math.NaN(), grpc.UnimplementedError
	}

	return float64(rand.Intn(360)), nil
}

func (f *fake) AngularVelocity(ctx context.Context, extra map[string]interface{}) (spatialmath.AngularVelocity, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.angvelOff {
		return spatialmath.AngularVelocity{}, grpc.UnimplementedError
	}

	return spatialmath.AngularVelocity{
		X: rand.Float64() * 10. * common.Randomsign(),
		Y: rand.Float64() * 10. * common.Randomsign(),
		Z: rand.Float64() * 10. * common.Randomsign(),
	}, nil
}

func (f *fake) LinearVelocity(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.linvelOff {
		return r3.Vector{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, grpc.UnimplementedError
	}
	return r3.Vector{
		X: rand.Float64() * 10. * common.Randomsign(),
		Y: rand.Float64() * 10. * common.Randomsign(),
		Z: rand.Float64() * 10. * common.Randomsign(),
	}, nil
}

func (f *fake) LinearAcceleration(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.linaccOff {
		return r3.Vector{}, grpc.UnimplementedError
	}

	return r3.Vector{
		X: rand.Float64() * 1.2 * common.Randomsign(),
		Y: rand.Float64() * 5 * common.Randomsign(),
		Z: rand.Float64() * 9.81 * common.Randomsign(),
	}, nil
}

func (f *fake) Accuracy(ctx context.Context, extra map[string]interface{}) (*movementsensor.Accuracy, error) {
	accmap := map[string]float32{
		"satellites_noise_signal": rand.Float32() * 20,
		"rand_trust_level":        rand.Float32() * 5,
	}
	return &movementsensor.Accuracy{
		AccuracyMap:        accmap,
		Hdop:               rand.Float32() * 2.,
		Vdop:               rand.Float32() * 5.,
		CompassDegreeError: rand.Float32(),
		NmeaFix:            rand.Int31n(5),
	}, nil
}

func (f *fake) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	defaults, err := movementsensor.DefaultAPIReadings(ctx, f, extra)
	if err != nil {
		return nil, err
	}
	defaults["foo"] = "bar"
	defaults["satellites"] = rand.Intn(32)
	return defaults, nil
}

func (f *fake) Properties(ctx context.Context, extra map[string]interface{}) (*movementsensor.Properties, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return &movementsensor.Properties{
		PositionSupported:           !f.posOff,
		OrientationSupported:        !f.oriOff,
		CompassHeadingSupported:     !f.compassOff,
		LinearVelocitySupported:     !f.linvelOff,
		AngularVelocitySupported:    !f.angvelOff,
		LinearAccelerationSupported: !f.linaccOff,
	}, nil
}
