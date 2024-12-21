package main

import (
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
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
)

var ModuleFamily = resource.NewModelFamily("rand", "go-fakes")

func main() {
	module.ModularMain(
		// Arms
		resource.APIModel{arm.API, fakearm.Model},
		resource.APIModel{arm.API, fakearm.EmptyModel},
		resource.APIModel{arm.API, fakearm.StaticModel},
		resource.APIModel{arm.API, fakearm.ErroringModel},
		resource.APIModel{arm.API, fakearm.NaNModel},

		// Bases
		resource.APIModel{base.API, fakebase.Model},
		resource.APIModel{base.API, fakebase.EmptyModel},
		resource.APIModel{base.API, fakebase.NanModel},

		// Boards
		resource.APIModel{board.API, fakeboard.Model},
		resource.APIModel{board.API, fakeboard.EmptyModel},

		// Cameras
		resource.APIModel{camera.API, fakecamera.Model},

		// Encoders
		resource.APIModel{encoder.API, fakeencoder.Model},

		// Gantries
		resource.APIModel{gantry.API, fakegantry.Model},
		resource.APIModel{gantry.API, fakegantry.EmptyModel},
		resource.APIModel{gantry.API, fakegantry.StaticModel},
		resource.APIModel{gantry.API, fakegantry.ErroringModel},

		// Grippers
		resource.APIModel{gripper.API, fakegripper.Model},

		// Input Controllers
		resource.APIModel{input.API, fakeinput.Model},
		resource.APIModel{input.API, fakeinput.EmptyModel},

		// Motors
		resource.APIModel{motor.API, fakemotor.Model},
		resource.APIModel{motor.API, fakemotor.EmptyModel},
		resource.APIModel{motor.API, fakemotor.StaticModel},

		// MovementSensors
		resource.APIModel{movementsensor.API, fakemovementsensor.Model},
		resource.APIModel{movementsensor.API, fakemovementsensor.EmptyModel},
		resource.APIModel{movementsensor.API, fakemovementsensor.StaticModel},
		resource.APIModel{movementsensor.API, fakemovementsensor.WaitingModel},
		resource.APIModel{movementsensor.API, fakemovementsensor.ErroringModel},
		resource.APIModel{movementsensor.API, fakemovementsensor.NaNModel},

		// PowerSensors
		resource.APIModel{powersensor.API, fakepowersensor.Model},
		resource.APIModel{powersensor.API, fakepowersensor.EmptyModel},
		resource.APIModel{powersensor.API, fakepowersensor.StaticModel},
		resource.APIModel{powersensor.API, fakepowersensor.ErroringModel},
		resource.APIModel{powersensor.API, fakepowersensor.WaitingModel},

		// Servos
		resource.APIModel{servo.API, fakeservo.Model},
		resource.APIModel{servo.API, fakeservo.StaticModel},

		// Sensors
		resource.APIModel{sensor.API, fakesensor.Model},
		resource.APIModel{sensor.API, fakesensor.EmptyModel},
		resource.APIModel{sensor.API, fakesensor.StaticModel},
		resource.APIModel{sensor.API, fakesensor.WaitingModel},
		resource.APIModel{sensor.API, fakesensor.ErrroingModel},
	)
}
