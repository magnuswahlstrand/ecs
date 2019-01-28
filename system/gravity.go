package system

import (
	"github.com/kyeett/ecs/constants"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Gravity is responsible for adding gravity to all entitites with a velocity type. It doesn't send or receive any events
type Gravity struct {
	em  *entity.Manager
	log logging.Logger
}

// NewGravity creates a new gravity system
func NewGravity(em *entity.Manager, logger logging.Logger) *Gravity {

	return &Gravity{
		em:  em,
		log: logger.WithField("s", "gravity"),
	}
}

// Update the input system
func (g *Gravity) Update(dt float64) {
	for _, e := range g.em.FilteredEntities(components.VelocityType, components.JoystickType) {
		if g.em.HasComponents(e, components.ParentedType) {
			continue
		}
		g.log.WithField("id", e).Debugf("update")
		v := g.em.Velocity(e)
		v.Y += constants.Gravity * dt
	}
}
