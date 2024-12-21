package fakebase

import (
	"context"
	"math"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/components/base"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

type nanBase struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	resource.Named

	logger logging.Logger
}

var nan = math.NaN()

func newNanBase(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (base.Base, error) {

	return &nanBase{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}, nil
}

func (n *nanBase) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return nil, nil
}

func (n *nanBase) Properties(context.Context, map[string]interface{}) (base.Properties, error) {
	return base.Properties{
		WidthMeters:              nan,
		TurningRadiusMeters:      nan,
		WheelCircumferenceMeters: nan,
	}, nil
}

func (n *nanBase) IsMoving(context.Context) (bool, error) {
	return false, nil
}

func (n *nanBase) SetPower(context.Context, r3.Vector, r3.Vector, map[string]interface{}) error {
	return nil
}

func (n *nanBase) SetVelocity(context.Context, r3.Vector, r3.Vector, map[string]interface{}) error {
	return nil
}

func (n *nanBase) MoveStraight(context.Context, int, float64, map[string]interface{}) error {
	return nil
}

func (n *nanBase) Spin(context.Context, float64, float64, map[string]interface{}) error {
	return nil
}

func (n *nanBase) Stop(context.Context, map[string]interface{}) error {
	return nil
}
