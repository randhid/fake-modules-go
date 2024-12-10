package fakeinput

import (
	"context"

	"go.viam.com/rdk/components/input"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

type emptyInput struct {
	resource.TriviallyCloseable
	resource.AlwaysRebuild
	resource.Named
	logger logging.Logger
}

func newEmptyInput(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (input.Controller, error) {
	return &emptyInput{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}, nil
}

// Controls implements input.Controller.
func (e *emptyInput) Controls(ctx context.Context, extra map[string]interface{}) ([]input.Control, error) {
	return nil, nil
}

// Events implements input.Controller.
func (e *emptyInput) Events(ctx context.Context, extra map[string]interface{}) (map[input.Control]input.Event, error) {
	return nil, nil
}

// RegisterControlCallback implements input.Controller.
func (e *emptyInput) RegisterControlCallback(ctx context.Context, control input.Control, triggers []input.EventType, ctrlFunc input.ControlFunction, extra map[string]interface{}) error {
	return nil
}
