package fakearm

import (
	_ "embed"

	"go.viam.com/rdk/referenceframe"
)

//go:embed kinematics.json
var kinematics []byte

func makeModelFrame(name string) (referenceframe.Model, error) {
	return referenceframe.UnmarshalModelJSON(kinematics, name)
}
