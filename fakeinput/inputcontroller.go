package fakeinput

import (
	"context"
	"errors"
	"fake-modules-go/common"
	"math/rand"
	"sync"
	"time"

	"go.viam.com/rdk/components/input"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils"
)

const inputName = "fake-input"

var Model = common.FakesFamily.WithModel(inputName)

func init() {
	resource.RegisterComponent(input.API, Model, resource.Registration[input.Controller, *Config]{
		Constructor: newFakeInput,
	})
}

func newFakeInput(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	input.Controller, error,
) {
	return setUpInputController(ctx, conf, logger)
}

// Config can list input structs (with their states), define event values and callback delays.
type Config struct {
	resource.TriviallyValidateConfig
	controls []input.Control

	// EventValue will dictate the value of the events returned. Random between -1 to 1 if unset.
	EventValue *float64 `json:"event_value,omitempty"`

	// CallbackDelaySec is the amount of time between callbacks getting triggered. Random between (1-2] sec if unset.
	// 0 is not valid and will be overwritten by a random delay.
	CallbackDelaySec float64 `json:"callback_delay_sec"`
}

type callback struct {
	control  input.Control
	triggers []input.EventType
	ctrlFunc input.ControlFunction
}

// setUpInputController returns a fake input.Controller.
func setUpInputController(ctx context.Context, conf resource.Config, logger logging.Logger) (input.Controller, error) {
	closeCtx, cancelFunc := context.WithCancel(context.Background())

	f := &fake{
		Named:      conf.ResourceName().AsNamed(),
		closeCtx:   closeCtx,
		cancelFunc: cancelFunc,
		callbacks:  make([]callback, 0),
		logger:     logger,
	}

	if err := f.Reconfigure(ctx, nil, conf); err != nil {
		return nil, err
	}

	// start callback thread
	f.activeBackgroundWorkers.Add(1)
	utils.ManagedGo(func() {
		f.startCallbackLoop()
	}, f.activeBackgroundWorkers.Done)

	return f, nil
}

// An InputController fakes an input.Controller.
type fake struct {
	resource.Named

	closeCtx                context.Context
	cancelFunc              func()
	activeBackgroundWorkers sync.WaitGroup

	mu            sync.Mutex
	controls      []input.Control
	eventValue    *float64
	callbackDelay *time.Duration
	callbacks     []callback
	logger        logging.Logger
}

// Reconfigure updates the config of the controller.
func (f *fake) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	f.controls = newConf.controls
	f.eventValue = newConf.EventValue
	if newConf.CallbackDelaySec != 0 {
		// convert to milliseconds to avoid any issues with float to int conversions
		delay := time.Duration(newConf.CallbackDelaySec*1000) * time.Millisecond
		f.callbackDelay = &delay
	}
	return nil
}

// Controls lists the inputs of the gamepad.
func (f *fake) Controls(ctx context.Context, extra map[string]interface{}) ([]input.Control, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.controls) == 0 {
		return []input.Control{input.AbsoluteX, input.ButtonStart}, nil
	}
	return f.controls, nil
}

func (f *fake) eventVal() float64 {
	if f.eventValue != nil {
		return *f.eventValue
	}
	//nolint:gosec
	return rand.Float64()
}

// Events returns the a specified or random input.Event (the current state) for AbsoluteX.
func (f *fake) Events(ctx context.Context, extra map[string]interface{}) (map[input.Control]input.Event, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	eventsOut := make(map[input.Control]input.Event)

	eventsOut[input.AbsoluteX] = input.Event{Time: time.Now(), Event: input.PositionChangeAbs, Control: input.AbsoluteX, Value: f.eventVal()}
	return eventsOut, nil
}

// RegisterControlCallback registers a callback function to be executed on the specified trigger Event. The fake implementation will
// trigger the callback at a random or user-specified interval with a random or user-specified value.
func (f *fake) RegisterControlCallback(
	ctx context.Context,
	control input.Control,
	triggers []input.EventType,
	ctrlFunc input.ControlFunction,
	extra map[string]interface{},
) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.callbacks = append(f.callbacks, callback{control: control, triggers: triggers, ctrlFunc: ctrlFunc})
	return nil
}

func (f *fake) startCallbackLoop() {
	for {
		var callbackDelay time.Duration

		if f.closeCtx.Err() != nil {
			return
		}

		f.mu.Lock()
		if f.callbackDelay != nil {
			callbackDelay = *f.callbackDelay
		} else {
			//nolint:gosec
			callbackDelay = 1*time.Second + time.Duration(rand.Float64()*1000)*time.Millisecond
		}
		f.mu.Unlock()

		if !utils.SelectContextOrWait(f.closeCtx, callbackDelay) {
			return
		}

		select {
		case <-f.closeCtx.Done():
			return
		default:
			f.mu.Lock()
			evValue := f.eventVal()
			for _, callback := range f.callbacks {
				for _, t := range callback.triggers {
					event := input.Event{Time: time.Now(), Event: t, Control: callback.control, Value: evValue}
					callback.ctrlFunc(f.closeCtx, event)
				}
			}
			f.mu.Unlock()
		}
	}
}

// TriggerEvent allows directly sending an Event (such as a button press) from external code.
func (f *fake) TriggerEvent(ctx context.Context, event input.Event, extra map[string]interface{}) error {
	return errors.New("unsupported")
}

// Close attempts to cleanly close the input controller.
func (f *fake) Close(ctx context.Context) error {
	f.mu.Lock()
	var err error
	if f.cancelFunc != nil {
		f.cancelFunc()
		f.cancelFunc = nil
	}
	f.mu.Unlock()

	f.activeBackgroundWorkers.Wait()
	return err
}
