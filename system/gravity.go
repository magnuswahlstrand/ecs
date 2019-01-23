package system

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Gravity is responsible for adding gravity to all entitites with a velocity type. It doesn't send or receive any events
type Gravity struct {
	em  *entity.Manager
	log logging.Logger
}

// NewGravity creates a new gravity system
func NewGravity(em *entity.Manager, ch chan events.Event, logger logging.Logger) *Gravity {

	return &Gravity{
		em:  em,
		log: logger.WithField("s", "gravity"),
	}
}

// Update the input system
func (g *Gravity) Update() {
	for _, e := range g.em.FilteredEntities(components.VelocityType) {
		g.log.WithField("id", e).Debugf("update")
		pos := g.em.Pos(e)
		pos.Y += 0.25
	}
}

// Send is an empty method to implement the System interface
func (g *Gravity) Send(ev events.Event) {
	g.log.Debugf("send, do nothing")
}
