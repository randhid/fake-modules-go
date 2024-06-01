package fakebase

import (
	"context"
	"errors"
	"fake-modules-go/common"
	"math"
	"sync"
	"time"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/components/base"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

const (
	baseName                    = "fake-base"
	defaultWidthM               = 100
	defaultWheelCircumferenceMm = 25
)

var Model = common.FakesFamily.WithModel(baseName)

type Config struct {
	resource.TriviallyValidateConfig
	Width              float64 `json:"width_mm,omitempty"`
	WheelCircumference float64 `json:"wheel_circumference_mm,omitempty"`
}

func init() {
	resource.RegisterComponent(base.API, Model, resource.Registration[base.Base, *Config]{
		Constructor: newFakeBase,
	})
}

type fake struct {
	resource.Named
	resource.AlwaysRebuild
	logger logging.Logger

	mu         sync.Mutex
	moving     bool
	properties base.Properties
	stopmoving func()
}

func newFakeBase(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	base.Base, error,
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
	f.mu.Lock()
	defer f.mu.Unlock()
	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}

	width := newConf.Width * 0.001
	if width == 0 {
		width = defaultWidthM * 0.01
	}

	wheelc := newConf.WheelCircumference
	if wheelc == 0 {
		wheelc = defaultWheelCircumferenceMm
	}

	f.properties = base.Properties{WidthMeters: width, WheelCircumferenceMeters: wheelc}

	return nil
}

func (f *fake) SetPower(ctx context.Context, linear, angular r3.Vector, extra map[string]interface{}) error {
	if linear.Norm() == 0 && angular.Norm() == 0 {
		err := f.Stop(ctx, extra)
		return errors.Join(err, errors.New("trying to move base with a {0,0,0} linear and angular power request"))
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.moving = true
	return nil
}

func (f *fake) MoveStraight(ctx context.Context, distance int, mmPerSec float64, extra map[string]interface{}) error {
	if distance == 0 || mmPerSec == 0 {
		return errors.New("trying to move base with 0 mm_per_seconds or 0 distance, not moving")
	}
	f.mu.Lock()
	f.moving = true
	f.mu.Unlock()

	waitSec := math.Abs(float64(distance) / mmPerSec)
	sleepFor := time.Duration(waitSec * 1e9)
	time.Sleep(sleepFor)

	f.mu.Lock()
	defer f.mu.Unlock()
	f.moving = false
	return nil
}

func (f *fake) Spin(ctx context.Context, angle, degsPerSec float64, extra map[string]interface{}) error {
	if angle == 0 || degsPerSec == 0 {
		return errors.New("trying to move base with 0 mm_per_seconds or 0 distance, not moving")
	}
	f.mu.Lock()
	f.moving = true
	f.mu.Unlock()

	waitSec := math.Abs(angle / degsPerSec)
	sleepFor := time.Duration(waitSec * 1e9)
	time.Sleep(sleepFor)

	f.mu.Lock()
	defer f.mu.Unlock()
	f.moving = false
	return nil
}

func (f *fake) SetVelocity(ctx context.Context, linear, angular r3.Vector, extra map[string]interface{}) error {
	if linear.Norm() == 0 && angular.Norm() == 0 {
		err := f.Stop(ctx, extra)
		return errors.Join(err, errors.New("trying to move base with a {0,0,0} linear and angular velocity request"))
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.moving = true
	return nil
}

func (f *fake) IsMoving(ctx context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.moving, nil
}

func (f *fake) Stop(ctx context.Context, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.moving = false
	return nil
}

func (f *fake) Close(ctx context.Context) error {
	return f.Stop(ctx, nil)
}

func (f *fake) Properties(ctx context.Context, extra map[string]interface{}) (base.Properties, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.properties, nil
}

func (f *fake) Geometries(ctx context.Context, extra map[string]interface{}) ([]spatialmath.Geometry, error) {
	box, err := spatialmath.NewBox(
		spatialmath.NewPose(r3.Vector{X: 0, Y: 0, Z: 0}, spatialmath.NewZeroPose().Orientation()),
		r3.Vector{},
		f.Name().ShortName(),
	)

	return []spatialmath.Geometry{box}, err
}
