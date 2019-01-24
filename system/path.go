package system

import (
	"fmt"

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
		path := p.em.Path("path_1")
		onPath := p.em.OnPath(e)

		// Find direction towards point
		pos := p.em.Pos(e)

		// Handle endpoints
		atEnd := pos.Eq(path.Points[len(path.Points)-1])
		// atStart := pos.Eq(path.Points[0])
		if atEnd {
			p.log.Debugf("%q reached end of path", e)
			switch onPath.Mode {
			case pathanimation.LinearLoop:
				// Loop to beginning
				pos.Vec = path.Points[0]
				onPath.Target = 1
				p.log.Debugf("loop to beginning")
			case pathanimation.LinearPingPong:
				// Turn around
				onPath.Target = len(path.Points) - 2
				onPath.Direction = -1
				p.log.Debugf("turn around")
			}
		}

		v := p.em.Velocity(e)
		target := path.Points[onPath.Target]

		// Next point is closer than speed, set velocity to reach point exactly
		to := target.Sub(pos.Vec)
		fmt.Printf("real diff %0.2f, speed: %0.2f\n, target: %d", to.Len(), onPath.Speed, onPath.Direction)
		if to.Len() < onPath.Speed {
			v.Vec = to
			onPath.Target += onPath.Direction
		}

		p.log.Debugf("%q pathing %q", e, path.Label)

	}
}

// Send is an empty method to implement the System interface
func (p *Path) Send(ev events.Event) {
	p.log.Debugf("send, do nothing")
}
