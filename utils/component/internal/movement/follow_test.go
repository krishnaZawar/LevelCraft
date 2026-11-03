package movement

import (
	"encoding/json"
	"testing"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/stretchr/testify/assert"
)

func Test_MarshalAndUnmarshalFollow(t *testing.T) {
	expected := NewFollow(250, FollowTarget_Point, 10)

	found := NewBaseFollow()

	details := expected.GetComponentDetails()
	err := found.BuildFromDetails(details.Data)

	assert.Nil(t, err)
	assert.Equal(t, base.ComponentName_Follow, details.Name)
	assert.Equal(t, expected.GetSpeed(), found.GetSpeed())
	assert.Equal(t, expected.GetTargetMode(), found.GetTargetMode())
	assert.Equal(t, expected.GetStopDistance(), found.GetStopDistance())
}

func Test_FollowRejectsNegativeValues(t *testing.T) {
	tests := []struct {
		name    string
		payload followJSON
		wantErr error
	}{
		{"negative speed", followJSON{Speed: -1}, base.ErrNegativeSpeed},
		{"negative stop distance", followJSON{StopDistance: -1}, base.ErrNegativeStopDistance},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.payload)
			assert.Nil(t, err)

			assert.ErrorIs(t, NewBaseFollow().BuildFromDetails(data), tt.wantErr)
		})
	}
}

func Test_FollowRejectsAnUnknownTargetMode(t *testing.T) {
	data, err := json.Marshal(followJSON{Speed: 100, TargetMode: "whatever"})
	assert.Nil(t, err)

	assert.Error(t, NewBaseFollow().BuildFromDetails(data))
}

// Written before the mode existed, or by a client that left it out.
func Test_FollowDefaultsToFollowingAnObject(t *testing.T) {
	data, err := json.Marshal(followJSON{Speed: 100})
	assert.Nil(t, err)

	follow := NewBaseFollow()
	assert.Nil(t, follow.BuildFromDetails(data))
	assert.Equal(t, FollowTarget_Object, follow.GetTargetMode())
}

// Defaulting to the origin would drag every newly added follower to the
// corner of the scene before it has been given a target.
func Test_FollowWithoutATargetIsInert(t *testing.T) {
	follow := NewBaseFollow()

	assert.False(t, follow.HasTarget())
	assert.Equal(t, "", follow.GetTargetID())

	_, _, isPoint := follow.GetTargetPoint()
	assert.False(t, isPoint)
}

func Test_FollowMovesTowardsItsTarget(t *testing.T) {
	follow := NewFollow(100, FollowTarget_Object, 0)

	x, y := follow.VelocityTowards(0, 0, 100, 0, 0.1)

	assert.Equal(t, float64(100), x)
	assert.Equal(t, float64(0), y)
}

// Without a stop distance the follower overshoots its target every frame and
// turns back, which reads as a jitter on top of it.
func Test_FollowStopsWithinItsStopDistance(t *testing.T) {
	follow := NewFollow(100, FollowTarget_Object, 20)

	x, y := follow.VelocityTowards(90, 0, 100, 0, 0.1)

	assert.Equal(t, float64(0), x, "already inside the stop distance")
	assert.Equal(t, float64(0), y)
}

func Test_FollowLandsExactlyOnItsStopDistance(t *testing.T) {
	follow := NewFollow(100, FollowTarget_Object, 10)

	// 90px separate them and 10 should remain, but a 1s step would cover 100
	x, y := follow.VelocityTowards(0, 0, 100, 0, 1)

	assert.Equal(t, float64(90), x, "velocity times dt should stop exactly at the stop distance")
	assert.Equal(t, float64(0), y)
}

// The id is meaningless in point mode, and keeping it would resume a chase
// nobody asked for if the mode were switched back.
func Test_FollowDropsTheObjectIDInPointMode(t *testing.T) {
	data, err := json.Marshal(followJSON{
		Speed:      100,
		TargetMode: FollowTarget_Point,
		TargetID:   "some-object",
		TargetX:    40,
		TargetY:    60,
	})
	assert.Nil(t, err)

	follow := NewBaseFollow()
	assert.Nil(t, follow.BuildFromDetails(data))

	assert.Equal(t, "", follow.GetTargetID())

	x, y, isPoint := follow.GetTargetPoint()
	assert.True(t, isPoint)
	assert.Equal(t, float64(40), x)
	assert.Equal(t, float64(60), y)
}

// A follower holds its ground rather than backing away, so a target that
// moves can close the gap below the stop distance on its own. The stop
// distance bounds the follower's approach, not the gap itself.
func Test_FollowDoesNotRetreatFromAnApproachingTarget(t *testing.T) {
	follow := NewFollow(100, FollowTarget_Object, 50)

	// the target has come closer than the follower would ever approach
	x, y := follow.VelocityTowards(0, 0, 10, 0, 0.1)

	assert.Equal(t, float64(0), x, "should stand still, not reverse")
	assert.Equal(t, float64(0), y)
}

func Test_FollowHandlesNoInput(t *testing.T) {
	assert.False(t, NewBaseFollow().HandlesInput("KEY_A"))
}
