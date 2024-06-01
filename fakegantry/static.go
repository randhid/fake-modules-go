package fakegantry

import (
	"context"
	"fmt"
	"sync"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/components/gantry"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
)

var (
	lengths = []float64{1., 2., 3.}
	homePos = []float64{0., 0., 0.}
)

type static struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable

	logger logging.Logger

	mu        sync.Mutex
	positions []float64
	model     referenceframe.Model
}

func newStaticGantry(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (
	gantry.Gantry, error,
) {

	m := referenceframe.NewSimpleModel("")
	for _, length := range lengths {
		f, err := referenceframe.NewTranslationalFrame(
			conf.ResourceName().ShortName(),
			r3.Vector{X: 1, Y: 0, Z: 0},
			referenceframe.Limit{Min: 0, Max: length})
		if err != nil {
			return nil, err
		}
		m.OrdTransforms = append(m.OrdTransforms, f)

	}

	s := &static{
		Named:     conf.ResourceName().AsNamed(),
		logger:    logger,
		model:     m,
		positions: homePos,
	}
	return s, nil
}

func (s *static) Position(ctx context.Context, extra map[string]interface{}) ([]float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.positions, nil
}

func (s *static) Lengths(ctx context.Context, extra map[string]interface{}) ([]float64, error) {
	return lengths, nil
}

func (s *static) Home(ctx context.Context, extra map[string]interface{}) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.positions = homePos
	return true, nil
}

func (s *static) MoveToPosition(ctx context.Context, target []float64, speeds []float64, extra map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(lengths) != len(target) {
		return fmt.Errorf(
			"gave array of %d targets for %d positions, their length must must match", len(target), len(s.positions),
		)
	}

	s.positions = target
	return nil
}

func (s *static) Stop(context.Context, map[string]interface{}) error {
	return nil
}

func (s *static) IsMoving(context.Context) (bool, error) {
	return false, nil
}

func (s *static) CurrentInputs(ctx context.Context) ([]referenceframe.Input, error) {
	pos, err := s.Position(ctx, nil)
	if err != nil {
		return []referenceframe.Input{}, err
	}
	inputs := referenceframe.FloatsToInputs(pos)
	return inputs, nil
}

func (s *static) GoToInputs(ctx context.Context, inputs ...[]referenceframe.Input) error {
	for _, input := range inputs {
		pos := referenceframe.InputsToFloats(input)
		if err := s.MoveToPosition(ctx, pos, []float64{}, nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *static) ModelFrame() referenceframe.Model {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.model
}
