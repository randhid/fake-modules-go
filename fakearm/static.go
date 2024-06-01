package fakearm

import (
	// for embedding model file.
	"context"
	_ "embed"
	"sync"

	pb "go.viam.com/api/component/arm/v1"
	"go.viam.com/rdk/components/arm"
	"go.viam.com/rdk/grpc"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/referenceframe"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/spatialmath"
)

type static struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable
	resource.Actuator

	logger logging.Logger

	// nil interfaces
	referenceframe.InputEnabled

	mu          sync.Mutex
	model       referenceframe.Model
	jointValues []float64
}

func newStaticArm(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (
	arm.Arm, error,
) {
	model, err := makeModelFrame(conf.Name)
	if err != nil {
		return nil, err
	}

	dof := len(model.DoF())
	jointValues := make([]float64, dof)

	s := &static{
		Named:       conf.ResourceName().AsNamed(),
		logger:      logger,
		model:       model,
		jointValues: jointValues,
	}

	return s, nil
}

// MoveToJointPositions sets the joints.
func (s *static) MoveToJointPositions(ctx context.Context, joints *pb.JointPositions, extra map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(joints.Values) != len(s.jointValues) {
		return nil
	}

	s.jointValues = joints.Values
	return nil
}

// JointPositions returns joints.
func (s *static) JointPositions(ctx context.Context, extra map[string]interface{}) (*pb.JointPositions, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	retJoint := &pb.JointPositions{Values: s.jointValues}
	return retJoint, nil
}

// ModelFrame returns the dynamic frame of the model.
func (s *static) ModelFrame() referenceframe.Model {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.model
}

// Geometries returns the list of geometries associated with the resource, in any order. The poses of the geometries reflect their
// current location relative to the frame of the resource.
func (s *static) Geometries(ctx context.Context, extra map[string]interface{}) ([]spatialmath.Geometry, error) {
	res, err := s.JointPositions(ctx, nil)
	if err != nil {
		return nil, err
	}
	inputs := s.model.InputFromProtobuf(res)
	gif, err := s.model.Geometries(inputs)
	if err != nil {
		return nil, err
	}
	return gif.Geometries(), nil
}

func (s *static) EndPosition(ctx context.Context, extra map[string]interface{}) (spatialmath.Pose, error) {
	return spatialmath.NewZeroPose(), grpc.UnimplementedError
}

func (s *static) MoveToPosition(ctx context.Context, pos spatialmath.Pose, extra map[string]interface{}) error {
	return grpc.UnimplementedError
}
