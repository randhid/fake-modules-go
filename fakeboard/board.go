package fakeboard

import (
	"context"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"time"

	"fake-modules-go/common"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	pb "go.viam.com/api/component/board/v1"
	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

const boardName = "fake-board"

var Model = common.FakesFamily.WithModel(boardName)

type Config struct {
	resource.TriviallyValidateConfig
	AnalogReaders     []board.AnalogReaderConfig     `json:"analogs,omitempty"`
	DigitalInterrupts []board.DigitalInterruptConfig `json:"digital_interrupts,omitempty"`
}

func init() {
	resource.RegisterComponent(board.API, Model, resource.Registration[board.Board, *Config]{
		Constructor: newFakeBoard,
	})
}

type fakeboard struct {
	resource.Named
	resource.AlwaysRebuild
	logger logging.Logger

	mu       sync.Mutex
	Analogs  map[string]*Analog
	Digitals map[string]*DigitalInterrupt
	GPIOPins map[string]*GPIOPin
}

func newFakeBoard(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	board.Board, error,
) {
	f := &fakeboard{
		Named:    conf.ResourceName().AsNamed(),
		logger:   logger,
		Analogs:  map[string]*Analog{},
		Digitals: map[string]*DigitalInterrupt{},
		GPIOPins: map[string]*GPIOPin{},
	}

	if err := f.Reconfigure(ctx, deps, conf); err != nil {
		return nil, err
	}

	return f, nil
}

// Validate ensures all parts of the config are valid.
func (conf *Config) Validate(path string) ([]string, error) {
	for idx, conf := range conf.AnalogReaders {
		if err := conf.Validate(fmt.Sprintf("%s.%s.%d", path, "analogs", idx)); err != nil {
			return nil, err
		}
	}
	for idx, conf := range conf.DigitalInterrupts {
		if err := conf.Validate(fmt.Sprintf("%s.%s.%d", path, "digital_interrupts", idx)); err != nil {
			return nil, err
		}
	}

	return []string{}, nil
}

func (f *fakeboard) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}

	// TODO(RSDK-2684): we dont configure pins so we just unset them here. not really great behavior.
	f.GPIOPins = map[string]*GPIOPin{}

	stillExists := map[string]struct{}{}

	for _, c := range newConf.AnalogReaders {
		stillExists[c.Name] = struct{}{}
		if curr, ok := f.Analogs[c.Name]; ok {
			if curr.pin != c.Pin {
				curr.reset(c.Pin)
			}
			continue
		}
		f.Analogs[c.Name] = newAnalogReader(c.Pin)
	}
	for name := range f.Analogs {
		if _, ok := stillExists[name]; ok {
			continue
		}
		delete(f.Analogs, name)
	}
	stillExists = map[string]struct{}{}

	var errs error
	for _, c := range newConf.DigitalInterrupts {
		stillExists[c.Name] = struct{}{}
		if curr, ok := f.Digitals[c.Name]; ok {
			if !reflect.DeepEqual(curr.conf, c) {
				curr.reset(c)
			}
			continue
		}
		var err error
		f.Digitals[c.Name], err = NewDigitalInterrupt(c)
		if err != nil {
			errs = multierr.Combine(errs, err)
		}
	}
	for name := range f.Digitals {
		if _, ok := stillExists[name]; ok {
			continue
		}
		delete(f.Digitals, name)
	}

	return nil
}

// AnalogByName returns the analog pin by the given name if it exists.
func (f *fakeboard) AnalogByName(name string) (board.Analog, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	a, ok := f.Analogs[name]
	if !ok {
		return nil, errors.Errorf("can't find AnalogReader (%s)", name)
	}
	return a, nil
}

// DigitalInterruptByName returns the interrupt by the given name if it exists.
func (f *fakeboard) DigitalInterruptByName(name string) (board.DigitalInterrupt, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	d, ok := f.Digitals[name]
	if !ok {
		return nil, fmt.Errorf("cant find DigitalInterrupt (%s)", name)
	}
	return d, nil
}

// GPIOPinByName returns the GPIO pin by the given name if it exists.
func (f *fakeboard) GPIOPinByName(name string) (board.GPIOPin, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	p, ok := f.GPIOPins[name]
	if !ok {
		pin := &GPIOPin{}
		f.GPIOPins[name] = pin
		return pin, nil
	}
	return p, nil
}

// AnalogNames returns the names of all known analog pins.
func (f *fakeboard) AnalogNames() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	names := []string{}
	for k := range f.Analogs {
		names = append(names, k)
	}
	return names
}

// DigitalInterruptNames returns the names of all known digital interrupts.
func (f *fakeboard) DigitalInterruptNames() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	names := []string{}
	for k := range f.Digitals {
		names = append(names, k)
	}
	return names
}

// SetPowerMode sets the board to the given power mode. If provided,
// the board will exit the given power mode after the specified
// duration.
func (f *fakeboard) SetPowerMode(ctx context.Context, mode pb.PowerMode, duration *time.Duration) error {
	return grpc.UnimplementedError
}

// StreamTicks starts a stream of digital interrupt ticks.
func (f *fakeboard) StreamTicks(ctx context.Context, interrupts []board.DigitalInterrupt, ch chan board.Tick,
	extra map[string]interface{},
) error {
	for _, i := range interrupts {
		name := i.Name()
		d, ok := f.Digitals[name]
		if !ok {
			return fmt.Errorf("could not find digital interrupt: %s", name)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Keep going
		}
		// Get a random bool for the high tick value.
		// linter complans about security but we don't care if someone
		// can predict if the fake interrupts will be high or low.
		//nolint:gosec
		randBool := rand.Int()%2 == 0
		select {
		case ch <- board.Tick{Name: d.conf.Name, High: randBool, TimestampNanosec: uint64(time.Now().Unix())}:
		default:
			// if nothing is listening to the channel just do nothing.
		}
	}
	return nil
}

// Close attempts to cleanly close each part of the board.
func (f *fakeboard) Close(ctx context.Context) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	return nil
}

// An Analog reads back the same set value.
type Analog struct {
	pin        string
	Value      int
	CloseCount int
	Mu         sync.RWMutex
	fakeValue  int
}

func newAnalogReader(pin string) *Analog {
	return &Analog{pin: pin}
}

func (a *Analog) reset(pin string) {
	a.Mu.Lock()
	a.pin = pin
	a.Value = 0
	a.Mu.Unlock()
}

func (a *Analog) Read(ctx context.Context, extra map[string]interface{}) (board.AnalogValue, error) {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	a.fakeValue++
	a.fakeValue %= 1001
	a.Value = a.fakeValue

	return board.AnalogValue{Value: a.Value, Min: 0, Max: 1000, StepSize: 1}, nil
}

func (a *Analog) Write(ctx context.Context, value int, extra map[string]interface{}) error {
	a.Set(value)
	return nil
}

// Set is used to set the value of an Analog.
func (a *Analog) Set(value int) {
	a.Mu.Lock()
	defer a.Mu.Unlock()
	a.Value = value
}

// A GPIOPin reads back the same set values.
type GPIOPin struct {
	high    bool
	pwm     float64
	pwmFreq uint

	mu sync.Mutex
}

// Set sets the pin to either low or high.
func (gp *GPIOPin) Set(ctx context.Context, high bool, extra map[string]interface{}) error {
	gp.mu.Lock()
	defer gp.mu.Unlock()

	gp.high = high
	gp.pwm = 0
	gp.pwmFreq = 0
	return nil
}

// Get gets the high/low state of the pin.
func (gp *GPIOPin) Get(ctx context.Context, extra map[string]interface{}) (bool, error) {
	gp.mu.Lock()
	defer gp.mu.Unlock()

	return gp.high, nil
}

// PWM gets the pin's given duty cycle.
func (gp *GPIOPin) PWM(ctx context.Context, extra map[string]interface{}) (float64, error) {
	gp.mu.Lock()
	defer gp.mu.Unlock()

	return gp.pwm, nil
}

// SetPWM sets the pin to the given duty cycle.
func (gp *GPIOPin) SetPWM(ctx context.Context, dutyCyclePct float64, extra map[string]interface{}) error {
	gp.mu.Lock()
	defer gp.mu.Unlock()

	gp.pwm = dutyCyclePct
	return nil
}

// PWMFreq gets the PWM frequency of the pin.
func (gp *GPIOPin) PWMFreq(ctx context.Context, extra map[string]interface{}) (uint, error) {
	gp.mu.Lock()
	defer gp.mu.Unlock()

	return gp.pwmFreq, nil
}

// SetPWMFreq sets the given pin to the given PWM frequency.
func (gp *GPIOPin) SetPWMFreq(ctx context.Context, freqHz uint, extra map[string]interface{}) error {
	gp.mu.Lock()
	defer gp.mu.Unlock()

	gp.pwmFreq = freqHz
	return nil
}

// DigitalInterrupt is a fake digital interrupt.
type DigitalInterrupt struct {
	mu    sync.Mutex
	conf  board.DigitalInterruptConfig
	value int64
}

// NewDigitalInterrupt returns a new fake digital interrupt.
func NewDigitalInterrupt(conf board.DigitalInterruptConfig) (*DigitalInterrupt, error) {
	return &DigitalInterrupt{
		conf: conf,
	}, nil
}

func (s *DigitalInterrupt) reset(conf board.DigitalInterruptConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conf = conf
}

// Value returns the current value of the interrupt which is
// based on the type of interrupt.
func (s *DigitalInterrupt) Value(ctx context.Context, extra map[string]interface{}) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.value++
	return s.value, nil
}

// Name returns the name of the digital interrupt.
func (s *DigitalInterrupt) Name() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conf.Name
}
