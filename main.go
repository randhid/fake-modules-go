package main

import (
	"context"
	"fake-modules-go/fakearm"
	"fake-modules-go/fakebase"
	"fake-modules-go/fakeboard"
	"fake-modules-go/fakecamera"
	"fake-modules-go/fakeencoder"
	"fake-modules-go/fakegantry"
	"fake-modules-go/fakegripper"
	"fake-modules-go/fakeinput"
	"fake-modules-go/fakemotor"
	"fake-modules-go/fakemovementsensor"
	"fake-modules-go/fakepowersensor"
	"fake-modules-go/fakesensor"
	"fake-modules-go/fakeservo"

	"go.viam.com/rdk/components/arm"
	"go.viam.com/rdk/components/base"
	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/camera"
	"go.viam.com/rdk/components/encoder"
	"go.viam.com/rdk/components/gantry"
	"go.viam.com/rdk/components/gripper"
	"go.viam.com/rdk/components/input"
	"go.viam.com/rdk/components/motor"
	"go.viam.com/rdk/components/movementsensor"
	"go.viam.com/rdk/components/powersensor"
	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils"
)

var ModuleFamily = resource.NewModelFamily("rand", "go-fakes")

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("Fake Go Modules"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	// instantiates the module itself
	fakes, err := module.NewModuleFromArgs(ctx)
	if err != nil {
		return err
	}

	// Arms
	if err = fakes.AddModelFromRegistry(ctx, arm.API, fakearm.Model); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, arm.API, fakearm.EmptyModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, arm.API, fakearm.StaticModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, arm.API, fakearm.ErroringModel); err != nil {
		return err
	}

	// Bases
	if err = fakes.AddModelFromRegistry(ctx, base.API, fakebase.Model); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, base.API, fakebase.EmptyModel); err != nil {
		return err
	}

	// Boards
	if err = fakes.AddModelFromRegistry(ctx, board.API, fakeboard.Model); err != nil {
		return err
	}

	// Cameras
	if err = fakes.AddModelFromRegistry(ctx, camera.API, fakecamera.Model); err != nil {
		return err
	}

	// Encoders
	if err = fakes.AddModelFromRegistry(ctx, encoder.API, fakeencoder.Model); err != nil {
		return err
	}

	// Gantries
	if err = fakes.AddModelFromRegistry(ctx, gantry.API, fakegantry.Model); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, gantry.API, fakegantry.EmptyModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, gantry.API, fakegantry.StaticModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, gantry.API, fakegantry.ErroringModel); err != nil {
		return err
	}

	// Grippers
	if err = fakes.AddModelFromRegistry(ctx, gripper.API, fakegripper.Model); err != nil {
		return err
	}

	// Input Controllers
	if err = fakes.AddModelFromRegistry(ctx, input.API, fakeinput.Model); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, input.API, fakeinput.EmptyModel); err != nil {
		return err
	}

	// Motors
	if err = fakes.AddModelFromRegistry(ctx, motor.API, fakemotor.Model); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, motor.API, fakemotor.EmptyModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, motor.API, fakemotor.StaticModel); err != nil {
		return err
	}

	// MovementSensors
	if err = fakes.AddModelFromRegistry(ctx, movementsensor.API, fakemovementsensor.Model); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, movementsensor.API, fakemovementsensor.EmptyModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, movementsensor.API, fakemovementsensor.StaticModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, movementsensor.API, fakemovementsensor.WaitingModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, movementsensor.API, fakemovementsensor.ErroringModel); err != nil {
		return err
	}

	// PowerSensors
	if err = fakes.AddModelFromRegistry(ctx, powersensor.API, fakepowersensor.Model); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, powersensor.API, fakepowersensor.EmptyModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, powersensor.API, fakepowersensor.StaticModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, powersensor.API, fakepowersensor.ErroringModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, powersensor.API, fakepowersensor.WaitingModel); err != nil {
		return err
	}

	// Servos
	if err = fakes.AddModelFromRegistry(ctx, servo.API, fakeservo.Model); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, servo.API, fakeservo.StaticModel); err != nil {
		return err
	}

	// Sensors
	if err = fakes.AddModelFromRegistry(ctx, sensor.API, fakesensor.Model); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, sensor.API, fakesensor.EmptyModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, sensor.API, fakesensor.StaticModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, sensor.API, fakesensor.WaitingModel); err != nil {
		return err
	}
	if err = fakes.AddModelFromRegistry(ctx, sensor.API, fakesensor.ErrroingModel); err != nil {
		return err
	}

	// Each module runs as its own process
	err = fakes.Start(ctx)
	logger.Warn("starting module")
	defer fakes.Close(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
