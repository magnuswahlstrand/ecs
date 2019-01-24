package system

import (
	"time"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Collision is handling collisions
type Collision struct {
	em    *entity.Manager
	outCh chan events.Event
	log   logging.Logger
}

// NewCollision creates a new Collision system
func NewCollision(em *entity.Manager, ch chan events.Event, logger logging.Logger) *Collision {
	return &Collision{
		em:    em,
		outCh: ch,
		log:   logger.WithField("s", "Collision"),
	}
}

// Update the Collision system
func (c *Collision) Update() {
	for _, e := range c.em.FilteredEntities(components.PosType, components.VelocityType) {
		pos := c.em.Pos(e)
		v := c.em.Velocity(e)
		// fmt.Println(pos.Vec)
		if pos.Y > 100 {
			v.Vec.Y = -5 * v.Vec.Y
			pos.Y = 100
			c.em.Add("camera_1", components.Shaking{Started: time.Now()})
		} else {
			pos.Vec = pos.Add(v.Vec)
		}
	}
}

// Send is an empty method to implement the System interface
func (c *Collision) Send(ev events.Event) {}
