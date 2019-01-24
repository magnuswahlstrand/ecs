package system

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

// Friction is responsible for adding Friction to all entitites with a velocity type. It doesn't send or receive any events
type Friction struct {
	em  *entity.Manager
	log logging.Logger
}

// NewFriction creates a new Friction system
func NewFriction(em *entity.Manager, ch chan events.Event, logger logging.Logger) *Friction {
	return &Friction{
		em:  em,
		log: logger.WithField("s", "Friction"),
	}
}

// Update the input system
func (f *Friction) Update(dt float64) {
	for _, e := range f.em.FilteredEntities(components.VelocityType) {
		v := f.em.Velocity(e)
		v.Vec = v.ScaledXY(gfx.V(0.90, 0.90))
	}
}

// Send is an empty method to implement the System interface
func (f *Friction) Send(ev events.Event) {}
