package system

import (
	"github.com/kyeett/ecs/constants"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/inputhandler"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

// Controls is responsible translating input to entity movement
type Controls struct {
	em  *entity.Manager
	log logging.Logger
}

// NewControls creates a new Controls system
func NewControls(em *entity.Manager, logger logging.Logger) *Controls {
	return &Controls{
		em:  em,
		log: logger.WithField("s", "controls"),
	}
}

// Update the Controls system
func (c *Controls) Update(dt float64) {

	for _, e := range c.em.FilteredEntities(components.PosType, components.VelocityType, components.JoystickType) {
		v := c.em.Velocity(e)

		if inputhandler.KeyJustPressed("up") && v.Y == 0 {

			if c.em.HasComponents(e, components.ParentedType) {
				c.log.Debugf("remove parenting from %s", e)
				removeParenting(e, c.em)
			}
			v.Vec = v.Add(gfx.V(0, -5))
		}

		if inputhandler.KeyPressed("left") {
			v.Vec = v.Add(gfx.V(-constants.AccelerationX*dt, 0))
		}

		if inputhandler.KeyPressed("right") {
			v.Vec = v.Add(gfx.V(constants.AccelerationX*dt, 0))

		}
	}
}

// Remove parenting, and add velocity of parent to entity
func removeParenting(e string, em *entity.Manager) {
	parented := em.Parented(e)
	v := em.Velocity(e)
	pV := em.Velocity(parented.ID)
	v.Vec = pV.Vec
	em.Remove(e, components.ParentedType)
}
