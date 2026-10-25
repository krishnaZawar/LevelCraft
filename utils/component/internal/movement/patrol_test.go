package movement

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

var route = []Point{{X: 0, Y: 0}, {X: 100, Y: 0}, {X: 100, Y: 100}}

func Test_MarshalAndUnmarshalPatrol(t *testing.T) {
	expected := NewPatrol(250, route, true)

	found := NewBasePatrol()

	details := expected.GetComponentDetails()
	err := found.BuildFromDetails(details.Data)

	assert.Nil(t, err)
	assert.Equal(t, base.ComponentName_Patrol, details.Name)
	assert.Equal(t, expected.GetSpeed(), found.GetSpeed())
	assert.Equal(t, expected.GetPoints(), found.GetPoints())
}

func Test_PatrolRejectsNegativeSpeed(t *testing.T) {
	data, err := json.Marshal(patrolJSON{Speed: -1, Points: route})
	assert.Nil(t, err)

	err = NewBasePatrol().BuildFromDetails(data)
	assert.ErrorIs(t, err, base.ErrNegativeSpeed)
}

// A component has to be addable in the editor before it has been configured,
// so an empty route is a valid state rather than a rejected one.
func Test_PatrolWithoutARouteIsInert(t *testing.T) {
	patrol := NewBasePatrol()

	x, y := patrol.NextVelocity(0, 0, 0.016)

	assert.Equal(t, float64(0), x)
	assert.Equal(t, float64(0), y)
}

// The route begins at its first point rather than assuming the object was
// placed on one, so an object dropped anywhere walks onto the route first.
func Test_PatrolMovesTowardsItsFirstPoint(t *testing.T) {
	patrol := NewPatrol(100, []Point{{X: 0, Y: 0}, {X: 100, Y: 0}}, false)

	x, y := patrol.NextVelocity(50, 0, 0.1)

	assert.Equal(t, float64(-100), x, "should head back towards the first point at the configured speed")
	assert.Equal(t, float64(0), y)
}

// Without clamping the final step the object steps past its point and spends
// every frame correcting, which reads as a jitter around that point.
func Test_PatrolLandsExactlyOnAPointItWouldOvershoot(t *testing.T) {
	patrol := NewPatrol(100, []Point{{X: 100, Y: 0}, {X: 0, Y: 0}}, false)

	// a 1 second step at speed 100 would carry it 100px, but only 10px remain
	x, y := patrol.NextVelocity(90, 0, 1)

	assert.Equal(t, float64(10), x, "velocity times dt should land exactly on the point")
	assert.Equal(t, float64(0), y)
	assert.Equal(t, 1, patrol.GetTargetIndex(), "and it should now head for the next point")
}

func Test_PatrolCyclesThroughTheRoute(t *testing.T) {
	patrol := NewPatrol(100, route, false)
	assert.Equal(t, 0, patrol.GetTargetIndex())

	// arrive at each point in turn; the index should walk forward and wrap
	for _, want := range []int{1, 2, 0, 1} {
		target := patrol.GetPoints()[patrol.GetTargetIndex()]
		patrol.NextVelocity(target.X, target.Y, 1)
		assert.Equal(t, want, patrol.GetTargetIndex())
	}
}

// Ping-pong retraces the route instead of jumping from the last point back to
// the first, which would look like a teleport.
func Test_PatrolPingPongsBackAlongTheRoute(t *testing.T) {
	patrol := NewPatrol(100, route, true)

	for _, want := range []int{1, 2, 1, 0, 1} {
		target := patrol.GetPoints()[patrol.GetTargetIndex()]
		patrol.NextVelocity(target.X, target.Y, 1)
		assert.Equal(t, want, patrol.GetTargetIndex())
	}
}

// The route can shrink between saves, and an index kept from the old one
// would address a point that no longer exists.
func Test_PatrolRestartsItsRouteOnRebuild(t *testing.T) {
	patrol := NewPatrol(100, route, false)
	patrol.NextVelocity(0, 0, 1)
	assert.Equal(t, 1, patrol.GetTargetIndex())

	data, err := json.Marshal(patrolJSON{Speed: 100, Points: []Point{{X: 5, Y: 5}}})
	assert.Nil(t, err)
	assert.Nil(t, patrol.BuildFromDetails(data))

	assert.Equal(t, 0, patrol.GetTargetIndex())
}

func Test_PatrolHandlesNoInput(t *testing.T) {
	assert.False(t, NewBasePatrol().HandlesInput("KEY_A"))
}
