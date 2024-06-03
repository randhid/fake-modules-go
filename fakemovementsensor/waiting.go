package fakemovementsensor

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"

	"fake-modules-go/common"

	"github.com/golang/geo/r3"
	geo "github.com/kellydunn/golang-geo"
	"go.viam.com/rdk/components/movementsensor"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

type WaitingConfig struct {
	resource.TriviallyValidateConfig
	WaitTime  int  `json:"wait_time_milli_seconds,omitempty"`
	PosOff    bool `json:"turn_off_position,omitempty"`
	OriOff    bool `json:"turn_off_orientation,omitempty"`
	CompOff   bool `json:"turn_off_compass_heading,omitempty"`
	AngVelOff bool `json:"turn_off_angular_velcoity,omitempty"`
	LinVelOff bool `json:"turn_off_linear_velocity,omitempty"`
	LinAccOff bool `json:"turn_off_linear_acceleration,omitempty"`
}

type waiting struct {
	resource.Named
	resource.TriviallyCloseable

	logger logging.Logger

	mu    sync.Mutex
	coord geo.Point

	waitTime   time.Duration
	posOff     bool
	oriOff     bool
	compassOff bool
	linvelOff  bool
	angvelOff  bool
	linaccOff  bool
}

func newWaitingMovementSensor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	movementsensor.MovementSensor, error,
) {
	w := &waiting{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}

	if err := w.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *waiting) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	newConf, err := resource.NativeConfig[WaitingConfig](conf)
	if err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	w.posOff = newConf.PosOff
	w.oriOff = newConf.OriOff
	w.compassOff = newConf.CompOff
	w.linvelOff = newConf.LinVelOff
	w.angvelOff = newConf.AngVelOff
	w.linaccOff = newConf.LinAccOff

	w.waitTime = common.DefaultWaitTimeMs // stes default to 500 milliseconds
	if newConf.WaitTime > 0 {             // otherwise set the wait time to the user-configured value
		w.waitTime = time.Duration(newConf.WaitTime) * time.Millisecond
	}

	return nil
}

func (w *waiting) Position(ctx context.Context, extra map[string]interface{}) (*geo.Point, float64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.posOff {
		return nil, math.NaN(), grpc.UnimplementedError
	}

	latincrement := rand.Float64() / 10e4 * common.Randomsign()
	lngincrement := rand.Float64() / 10e5 * math.Cos(latincrement) * common.Randomsign()
	altitude := rand.Float64() * 100

	newcoord := geo.NewPoint(
		w.coord.Lat()+latincrement,
		w.coord.Lng()+lngincrement,
	)

	time.Sleep(w.waitTime)
	return newcoord, altitude, nil
}

func (w *waiting) Orientation(ctx context.Context, extra map[string]interface{}) (spatialmath.Orientation, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.oriOff {
		return nil, grpc.UnimplementedError
	}

	ox := rand.Float64() * common.Randomsign()
	oy := rand.Float64() * common.Randomsign()
	oz := rand.Float64() * common.Randomsign()
	theta := rand.Float64()

	ov := spatialmath.OrientationVector{OX: ox, OY: oy, OZ: oz, Theta: theta}
	ov.Normalize()
	time.Sleep(w.waitTime)
	return &ov, nil
}

func (w *waiting) CompassHeading(ctx context.Context, extra map[string]interface{}) (float64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.compassOff {
		return math.NaN(), grpc.UnimplementedError
	}
	time.Sleep(w.waitTime)

	return float64(rand.Intn(360)), nil
}

func (w *waiting) AngularVelocity(ctx context.Context, extra map[string]interface{}) (spatialmath.AngularVelocity, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.angvelOff {
		return spatialmath.AngularVelocity{}, grpc.UnimplementedError
	}

	time.Sleep(w.waitTime)
	return spatialmath.AngularVelocity{
		X: rand.Float64() * 10. * common.Randomsign(),
		Y: rand.Float64() * 10. * common.Randomsign(),
		Z: rand.Float64() * 10. * common.Randomsign(),
	}, nil
}

func (w *waiting) LinearVelocity(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.linvelOff {
		return r3.Vector{X: math.NaN(), Y: math.NaN(), Z: math.NaN()}, grpc.UnimplementedError
	}
	time.Sleep(w.waitTime)
	return r3.Vector{
		X: rand.Float64() * 10. * common.Randomsign(),
		Y: rand.Float64() * 10. * common.Randomsign(),
		Z: rand.Float64() * 10. * common.Randomsign(),
	}, nil
}

func (w *waiting) LinearAcceleration(ctx context.Context, extra map[string]interface{}) (r3.Vector, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.linaccOff {
		return r3.Vector{}, grpc.UnimplementedError
	}
	time.Sleep(w.waitTime)
	return r3.Vector{
		X: rand.Float64() * 1.2 * common.Randomsign(),
		Y: rand.Float64() * 5 * common.Randomsign(),
		Z: rand.Float64() * 9.81 * common.Randomsign(),
	}, nil
}

func (w *waiting) Accuracy(ctx context.Context, extra map[string]interface{}) (*movementsensor.Accuracy, error) {
	accmap := map[string]float32{
		"satellites_noise_signal": rand.Float32() * 20,
		"rand_trust_level":        rand.Float32() * 5,
	}
	time.Sleep(w.waitTime)
	return &movementsensor.Accuracy{
		AccuracyMap:        accmap,
		Hdop:               rand.Float32() * 2.,
		Vdop:               rand.Float32() * 5.,
		CompassDegreeError: rand.Float32(),
		NmeaFix:            rand.Int31n(5),
	}, nil
}

func (w *waiting) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	defaults, err := movementsensor.DefaultAPIReadings(ctx, w, extra)
	if err != nil {
		return nil, err
	}
	defaults["foo"] = "bar"
	defaults["satellites"] = rand.Intn(32)
	time.Sleep(w.waitTime)
	return defaults, nil
}

func (w *waiting) Properties(ctx context.Context, extra map[string]interface{}) (*movementsensor.Properties, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	time.Sleep(w.waitTime)
	return &movementsensor.Properties{
		PositionSupported:           !w.posOff,
		OrientationSupported:        !w.oriOff,
		CompassHeadingSupported:     !w.compassOff,
		LinearVelocitySupported:     !w.linvelOff,
		AngularVelocitySupported:    !w.angvelOff,
		LinearAccelerationSupported: !w.linaccOff,
	}, nil
}
