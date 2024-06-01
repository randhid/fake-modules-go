package fakegantry

import (
	"context"
	"errors"
	"fake-modules-go/common"
	"fmt"
	"sync"
	"time"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/components/gantry"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/utils"
)

const (
	defaultspeed    = 2
	goalWithinRange = 0.2
)

type Config struct {
	Lengths []float64 `json:"lengths_array_mm"`
}

func (c Config) Validate(path string) ([]string, error) {
	if len(c.Lengths) == 0 {
		return nil, errors.New("need at least one non-zero element in lengths_array_mm in config attributes")
	}

	for i, val := range c.Lengths {
		if val <= 0 {
			return nil, fmt.Errorf("index %d of gantry array has negative or zero length, must be positive", i)
		}

	}

	return []string{}, nil
}

type fake struct {
	resource.Named
	logger logging.Logger

	mu         sync.Mutex
	lengths    []float64
	positions  []float64
	moving     bool
	stopmoving func()
	modelframe referenceframe.Model
}

func newFakeGantry(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	gantry.Gantry, error,
) {
	f := &fake{
		Named:     conf.ResourceName().AsNamed(),
		logger:    logger,
		positions: []float64{},
	}

	f.Reconfigure(ctx, deps, conf)

	return f, nil
}

func (f *fake) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	newConf, err := resource.NativeConfig[Config](conf)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.lengths = newConf.Lengths
	f.positions = make([]float64, len(newConf.Lengths))

	makeNewModel := func(lengths []float64) (referenceframe.Model, error) {
		m := referenceframe.NewSimpleModel(conf.ResourceName().ShortName())
		for _, length := range lengths {
			frame, err := referenceframe.NewTranslationalFrame(
				conf.ResourceName().ShortName(), r3.Vector{X: 1}, referenceframe.Limit{Min: 0, Max: length})
			if err != nil {
				return m, fmt.Errorf("error creating frame: %w", err)
			}
			m.OrdTransforms = append(m.OrdTransforms, frame)
		}
		return m, nil
	}

	f.modelframe, err = makeNewModel(f.lengths)
	if err != nil {
		return err
	}

	return nil
}

func (f *fake) MoveToPosition(ctx context.Context, target []float64, speeds []float64, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(target) != len(f.lengths) {
		return fmt.Errorf("received unequal length of positions request %v for gantry axes %v", target, f.lengths)
	}

	for i, val := range target {
		if val > f.lengths[i] || val < 0 {
			return fmt.Errorf(
				"trying to go to a negative or position outside gantry length, fake-gantry cannot do this"+
					"target position %.4f, length %.4f", val, f.lengths[i])
		}

		// check if speeds are not sent, they are optional
		if len(speeds) == 0 {
			speeds = make([]float64, len(f.lengths))
		}
		// check if any speeds are negative
		if speeds[i] < 0 {
			return fmt.Errorf(
				"trying to go to a negative speed %.4f, fake-gantry cannot do this", speeds[i])
		}

	}

	current := append([]float64{}, f.positions...) // Make a copy of the current positions
	increments := make([]float64, len(f.lengths))
	for i := range speeds {
		speed := speeds[i]
		if speeds[i] == 0 {
			speed = defaultspeed
		}
		increments[i] = speed * (target[i] - current[i])
	}

	var moveCtx context.Context
	moveCtx, f.stopmoving = context.WithCancel(context.Background())

	// simulate motion
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100) // Update every 100ms
		defer ticker.Stop()

		for {
			select {
			case <-moveCtx.Done():
				// Exit the goroutine if the context is canceled
				return
			case <-ticker.C:
				allReached := true
				// Loop through the lengths and update positions
				for i := range f.lengths {
					if !utils.Float64AlmostEqual(target[i], current[i], common.GoalWithinRange) {
						f.moving = true
						current[i] += increments[i]
						f.positions = current
						allReached = false
					}
				}
				if allReached {
					// All positions are within the goal range
					f.moving = false
					return
				}
			}
			f.logger.Debugf("joint positions %v", f.positions)

		}
	}()

	return nil
}

func (f *fake) Stop(ctx context.Context, extra map[string]interface{}) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.stopmoving != nil {
		f.stopmoving()
	}
	return nil
}

func (f *fake) Close(ctx context.Context) error {
	if err := f.Stop(ctx, nil); err != nil {
		return err
	}
	return nil
}

func (f *fake) Position(ctx context.Context, extra map[string]interface{}) ([]float64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.positions, nil
}

func (f *fake) Lengths(ctx context.Context, extra map[string]interface{}) ([]float64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.lengths, nil
}

func (f *fake) Home(ctx context.Context, extra map[string]interface{}) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	zeroes := make([]float64, len(f.lengths))
	f.positions = zeroes
	return true, nil
}

func (f *fake) IsMoving(ctx context.Context) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.moving, nil
}

func (f *fake) ModelFrame() referenceframe.Model {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.modelframe
	// return nil
}

func (f *fake) CurrentInputs(ctx context.Context) ([]referenceframe.Input, error) {
	pos, err := f.Position(ctx, nil)
	if err != nil {
		return []referenceframe.Input{}, err
	}
	inputs := referenceframe.FloatsToInputs(pos)
	return inputs, nil
}

func (f *fake) GoToInputs(ctx context.Context, inputs ...[]referenceframe.Input) error {
	for _, input := range inputs {
		pos := referenceframe.InputsToFloats(input)
		if err := f.MoveToPosition(ctx, pos, []float64{}, nil); err != nil {
			return err
		}
	}
	return nil
}
