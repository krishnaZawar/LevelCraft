package movement

import (
	"encoding/json"
	"math"

	"github.com/krishnaZawar/LevelCraft/utils/component/base"
	"github.com/krishnaZawar/LevelCraft/utils/component/internal/models"
)

const (
	// default speed value the component holds on init
	defaultPatrolSpeed = 100
)

// Point is a position in the world the object can be asked to move towards.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Patrol walks the object through a fixed route of points, with no input.
//
// A route of fewer than two points is inert rather than invalid: a component
// has to be addable in the editor before it has been configured, so the empty
// state is a valid one.
type Patrol struct {
	speed  float64 // the speed with which the object should move
	points []Point // the route, visited in order
	// whether the object reverses along the route on reaching the end, or
	// returns to the first point and goes round again
	pingPong bool

	targetIndex int // the point currently being moved towards
	step        int // +1 or -1, the direction being walked through points
}

// intermediary structure of Patrol used for marshalling and unmarshalling component details
type patrolJSON struct {
	Speed    float64 `json:"speed"`
	PingPong bool    `json:"pingPong"`
	Points   []Point `json:"points"`
}

// internal function used to register the base component copy with the componentRegistry
func NewBasePatrol() *Patrol {
	return &Patrol{
		speed: defaultPatrolSpeed,
		step:  1,
	}
}

func NewPatrol(speed float64, points []Point, pingPong bool) *Patrol {
	patrol := NewBasePatrol()
	patrol.speed = speed
	patrol.points = points
	patrol.pingPong = pingPong
	return patrol
}

// get the configured speed for the component
func (p *Patrol) GetSpeed() float64 {
	return p.speed
}

// get the route the object walks
func (p *Patrol) GetPoints() []Point {
	return p.points
}

// index of the point currently being moved towards
func (p *Patrol) GetTargetIndex() int {
	return p.targetIndex
}

// advances to the next point on the route, reversing or wrapping at the end
func (p *Patrol) advance() {
	if len(p.points) < 2 {
		return
	}

	if !p.pingPong {
		p.targetIndex = (p.targetIndex + 1) % len(p.points)
		return
	}

	// Turn around at either end rather than wrapping, so the object retraces
	// the route instead of jumping from the last point back to the first.
	if p.targetIndex+p.step < 0 || p.targetIndex+p.step >= len(p.points) {
		p.step = -p.step
	}
	p.targetIndex += p.step
}

// Returns the velocity that moves the object towards the point it is heading
// for, and advances to the next point once this frame's step would reach it.
//
// dt is needed to land exactly on a point: without it the object steps past
// and spends every frame correcting, which reads as a jitter around the point.
func (p *Patrol) NextVelocity(x float64, y float64, dt float64) (float64, float64) {
	if len(p.points) < 2 || dt <= 0 {
		return 0, 0
	}

	target := p.points[p.targetIndex]
	dx, dy := target.X-x, target.Y-y
	distance := math.Hypot(dx, dy)

	if distance == 0 {
		p.advance()
		return 0, 0
	}

	// This frame would reach or overshoot the point, so travel exactly the
	// remaining distance and aim at the next one from here on.
	if p.speed*dt >= distance {
		p.advance()
		return dx / dt, dy / dt
	}

	return (dx / distance) * p.speed, (dy / distance) * p.speed
}

// Returns the name of the component
func (p *Patrol) GetComponentName() string {
	return base.ComponentName_Patrol
}

// Returns a snapshot of the complete data stored in the component
func (p *Patrol) GetComponentDetails() models.ComponentDetails {
	data := patrolJSON{
		Speed:    p.speed,
		PingPong: p.pingPong,
		Points:   p.points,
	}
	byteData, _ := json.Marshal(data)
	return models.ComponentDetails{
		Name: p.GetComponentName(),
		Data: json.RawMessage(byteData),
	}
}

// Build component from provided details
func (p *Patrol) BuildFromDetails(data json.RawMessage) error {
	var componentData patrolJSON
	if err := json.Unmarshal(data, &componentData); err != nil {
		return err
	}

	if componentData.Speed < 0 {
		return base.ErrNegativeSpeed
	}

	p.speed = componentData.Speed
	p.pingPong = componentData.PingPong
	p.points = componentData.Points

	// The route may have shrunk since the last save, so start it again rather
	// than keeping an index that no longer addresses a point.
	p.targetIndex = 0
	p.step = 1

	return nil
}

// the route is fixed, so no input affects it
func (p *Patrol) HandlesInput(_ string) bool {
	return false
}

var _ models.Component = &Patrol{}
