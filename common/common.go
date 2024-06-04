package common

import (
	"math"
	"math/rand"
	"time"

	"go.viam.com/rdk/resource"
)

const (
	GoalWithinRange   = 0.2
	DefaultWaitTimeMs = 500 * time.Millisecond // default wait time in milllisecods
)

var (
	FakesFamily = resource.NewModelFamily("rand", "fake-modules-go")
)

func Randomsign() float64 {
	negative := rand.Int()%2 == 0
	if negative {
		return -1
	}
	return 1

}

func Sign(x float64) float64 {
	if x == 0 {
		return 0.
	}

	if math.Signbit(x) {
		return -1.0
	}
	return 1.0
}
