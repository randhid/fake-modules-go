package fakebase

import (
	"context"

	"github.com/golang/geo/r3"
	"go.viam.com/rdk/components/base"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

type emptyBase struct {
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	resource.Named

	logger logging.Logger
}

func newEmptyBase(_ context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger) (base.Base, error) {

	return &emptyBase{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}, nil
}

func (e *emptyBase) Geometries(context.Context, map[string]interface{}) ([]spatialmath.Geometry, error) {
	return nil, nil
}

func (e *emptyBase) Properties(context.Context, map[string]interface{}) (base.Properties, error) {
	return base.Properties{}, nil
}

func (e *emptyBase) IsMoving(context.Context) (bool, error) {
	return false, nil
}

func (e *emptyBase) SetPower(context.Context, r3.Vector, r3.Vector, map[string]interface{}) error {
	return nil
}

func (e *emptyBase) SetVelocity(context.Context, r3.Vector, r3.Vector, map[string]interface{}) error {
	return nil
}

func (e *emptyBase) MoveStraight(context.Context, int, float64, map[string]interface{}) error {
	return nil
}

func (e *emptyBase) Spin(context.Context, float64, float64, map[string]interface{}) error {
	return nil
}

func (e *emptyBase) Stop(context.Context, map[string]interface{}) error {
	return nil
}
