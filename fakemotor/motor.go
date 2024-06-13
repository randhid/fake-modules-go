package fakemotor

import (
	"context"
	"errors"
	"math"
	"sync"
	"time"

	"fake-modules-go/common"

	"go.viam.com/rdk/components/motor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/utils"
)

const (
	defaultRPM              = 100
	defaultTicksPerRotation = 1
	goalWithinRange         = 0.2
)

type Config struct {
	PositionOn bool    `json:"enable_position"`
	MaxSpeed   float64 `json:"max_speed_rpm,omitempty"`
}

func (c Config) Validate(path string) ([]string, error) {
	if c.MaxSpeed < 0 {
		return nil, errors.New("cannot have negative max_speed_rpm")
	}

	return []string{}, nil
}

type fake struct {
	resource.Named
	logger logging.Logger

	mu                sync.Mutex
	positionReporting bool
	moving            bool
	power             float64
	position          float64
	maxspeed          float64
	stopmoving        func()
}

func newFakeMotor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	motor.Motor, error,
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
	newConf, err := resource.NativeConfig[Config](conf)
	if err != nil {
		return err
	}

	if f.positionReporting = newConf.PositionOn; f.positionReporting {
		f.position = 0
	}

	f.maxspeed = newConf.MaxSpeed
	if f.maxspeed == 0 {
		f.maxspeed = defaultRPM
	}

	return nil
}
func (f *fake) SetPower(ctx context.Context, power float64, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.power = power
	f.moving = true
	return nil
}

func (f *fake) SetRPM(ctx context.Context, rpm float64, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	goal := math.Inf(1) * common.Sign(rpm)
	return f.GoFor(ctx, rpm, goal, extra)
}

func (f *fake) GoFor(ctx context.Context, rpm, revolutions float64, extra map[string]interface{}) error {
	currPos, err := f.Position(ctx, nil)
	f.logger.Infof("currPos at start %v", currPos)
	if err != nil {
		return errors.Join(errors.New("GoFor cannot be executed"), err)
	}
	if rpm == 0 {
		return motor.NewZeroRPMError()
	}
	goal := currPos + revolutions

	// time in MilliSeconds
	timeIncrement := 100 * time.Millisecond
	revIncrement := rpm /*revolutions_per_minute*/ / 60e3 /*milliseconds_per_minute*/ *
		100 /*milliseeconds*/

	f.logger.Infof("timeIncrement %v", timeIncrement)
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.stopmoving != nil {
		f.stopmoving()
	}

	f.moving = true
	f.power = rpm / f.maxspeed
	if math.Abs(f.power) > 1 {
		f.power = 1
	}

	var moveCtx context.Context
	moveCtx, f.stopmoving = context.WithCancel(context.Background())

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond) // Update every 100ms
		defer ticker.Stop()
		for {
			select {
			case <-moveCtx.Done():
				// Exit the goroutine if the context is canceled
				return
			case <-ticker.C:
				// if goal is infinite, this should go on forever.
				// I should factor out the simulation into a simulate move at some point to take in currentPos, goal and rpm 
				// and make it stop itself on a new call.
				// but not today.
				if !utils.Float64AlmostEqual(currPos, goal, common.GoalWithinRange) {
					f.moving = true
					currPos += revIncrement
					f.position = currPos
					f.logger.Infof("pos %v", f.position)
				}
			}
			f.logger.Debugf("position %v", f.position)
		}

	}()

	return nil
}

func (f *fake) GoTo(ctx context.Context, rpm, targetPos float64, extra map[string]interface{}) error {
	currPos, err := f.Position(ctx, nil)
	if err != nil {
		return errors.Join(errors.New("GoTo cannot be executed"), err)
	}

	revs := targetPos - currPos
	return f.GoFor(ctx, rpm, revs, nil)
}

func (f *fake) Position(ctx context.Context, extra map[string]interface{}) (float64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !f.positionReporting {
		return math.NaN(), errors.New("motor is not PositionReporting cannot report position")

	}
	return f.position, nil
}

func (f *fake) ResetZeroPosition(ctx context.Context, offset float64, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.position = 0
	return nil
}

func (f *fake) Properties(ctx context.Context, extra map[string]interface{}) (motor.Properties, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return motor.Properties{PositionReporting: f.positionReporting}, nil
}

func (f *fake) IsMoving(ctx context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.moving, nil
}

func (f *fake) IsPowered(ctx context.Context, extra map[string]interface{}) (bool, float64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.moving, f.power, nil
}

func (f *fake) Stop(ctx context.Context, extra map[string]interface{}) error {
	currPos, err := f.Position(ctx, nil)
	if err != nil {
		currPos = math.NaN()
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	if f.stopmoving != nil {
		f.stopmoving()
	}
	f.moving = false
	f.power = 0
	f.position = currPos
	return nil
}

func (f *fake) Close(ctx context.Context) error {
	return f.Stop(ctx, nil)
}
