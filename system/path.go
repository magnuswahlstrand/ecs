package system

import (
	"fmt"
	"math"

	"github.com/kyeett/gomponents/pathanimation"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Path is responsible for moving an entity on a path
type Path struct {
	em  *entity.Manager
	log logging.Logger
}

// NewPath creates a new Path system
func NewPath(em *entity.Manager, logger logging.Logger) *Path {
	return &Path{
		em:  em,
		log: logger.WithField("s", "path"),
	}
}

// Update the Path system
func (p *Path) Update(dt float64) {
	for _, e := range p.em.FilteredEntities(components.PosType, components.OnPathType) {
		// Find target path
		onPath := p.em.OnPath(e)

		path := p.em.Path(onPath.Label)
		switch path.Type {
		case pathanimation.Polygon:
			p.followPolygonPath(e, path, onPath, dt)
		case pathanimation.Ellipse:
			p.followCirclePath(e, path, onPath, dt)
		}

		p.log.Debugf("%q pathing %q", e, path.Label)

	}
}

func (p *Path) followPolygonPath(e string, path *components.Path, onPath *components.OnPath, dt float64) {
	// Find direction towards point
	pos := p.em.Pos(e)

	// Handle endpoints
	atEnd := onPath.Target == len(path.Points) //pos.Eq(path.Points[len(path.Points)-1])
	atStart := onPath.Target == -1             //pos.Eq(path.Points[0])

	switch onPath.Mode {
	case pathanimation.LinearLoop:
		if atEnd {
			pos.Vec = path.Points[0]
			onPath.Target = 1
			p.log.Debugf("loop to beginning")
		}
	case pathanimation.LinearPingPong:
		if atEnd {
			// Turn around
			onPath.Target = len(path.Points) - 2
			onPath.Direction = -1
			p.log.Debugf("at end, turn around")
		}

		if atStart {
			// Turn around
			onPath.Target = 1
			onPath.Direction = 1
			p.log.Debugf("at start, turn around")
		}
	}

	v := p.em.Velocity(e)
	target := path.Points[onPath.Target]

	// Next point is closer than speed, set velocity to reach point exactly
	to := target.Sub(pos.Vec)
	if to.Len() < onPath.Speed*dt {
		v.Vec = to.Scaled(1 / dt)
		onPath.Target += onPath.Direction
	} else {
		v.Vec = to.Unit().Scaled(onPath.Speed)
	}
}

func (p *Path) followCirclePath(e string, path *components.Path, onPath *components.OnPath, dt float64) {

	center := path.Points[0]

	// Find direction towards point
	pos := p.em.Pos(e)

	r := path.Points[1].Sub(center).Len()
	angleV := onPath.Speed * dt / (r)

	fmt.Println(pos, pos.Sub(center), pos.Sub(center).Angle(), r, onPath.Speed*dt, onPath.Speed*dt/(r))
	norm := pos.Sub(center).Rotated(angleV)
	target := norm.Add(center)
	fmt.Println(target.Angle())
	fmt.Printf("r:%f, angleV:%0.2f, target:%s, nAngle: %0.2f. norm:%s\n", r, angleV/(2*math.Pi), target, norm.Angle(), norm)

	// v := p.em.Velocity(e)
	// target := path.Points[onPath.Target]
	v := p.em.Velocity(e)
	to := target.Sub(pos.Vec)
	v.Vec = to.Scaled(1 / dt)
	fmt.Println(to)

	// // Next point is closer than speed, set velocity to reach point exactly
	// to := target.Sub(pos.Vec)
	// if to.Len() < onPath.Speed*dt {
	// 	v.Vec = to.Scaled(1 / dt)
	// 	onPath.Target += onPath.Direction
	// } else {
	// 	v.Vec = to.Unit().Scaled(onPath.Speed)
	// }
}

// Send is an empty method to implement the System interface
func (p *Path) Send(ev events.Event) {
	p.log.Debugf("send, do nothing")
}
