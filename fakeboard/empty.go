package fakeboard

import (
	"context"
	"time"

	pb "go.viam.com/api/component/board/v1"
	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type emptyBoard struct {
	resource.TriviallyCloseable
	resource.AlwaysRebuild
	resource.Named

	logger logging.Logger
}

func newEmptyBoard(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (board.Board, error) {
	return &emptyBoard{
		logger: logger,
		Named:  conf.ResourceName().AsNamed(),
	}, nil
}

// AnalogByName implements board.Board.
func (e *emptyBoard) AnalogByName(string) (board.Analog, error) {
	return nil, nil
}

// AnalogNames implements board.Board.
func (e *emptyBoard) AnalogNames() []string {
	return nil
}

// DigitalInterruptByName implements board.Board.
func (e *emptyBoard) DigitalInterruptByName(string) (board.DigitalInterrupt, error) {
	return nil, nil
}

// DigitalInterruptNames implements board.Board.
func (e *emptyBoard) DigitalInterruptNames() []string {
	return nil
}

// GPIOPinByName implements board.Board.
func (e *emptyBoard) GPIOPinByName(string) (board.GPIOPin, error) {
	return nil, nil
}

// StreamTicks implements board.Board.
func (e *emptyBoard) StreamTicks(context.Context, []board.DigitalInterrupt, chan board.Tick, map[string]interface{}) error {
	return nil
}

func (e *emptyBoard) SetPowerMode(context.Context, pb.PowerMode, *time.Duration) error {
	return nil
}
