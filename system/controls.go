package system

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/constants"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

// Controls is responsible translating input to entity movement
type Controls struct {
	em     *entity.Manager
	outCh  chan events.Event
	events []events.Event
	log    logging.Logger
}

// NewControls creates a new Controls system
func NewControls(em *entity.Manager, ch chan events.Event, logger logging.Logger) *Controls {
	return &Controls{
		em:     em,
		outCh:  ch,
		events: []events.Event{},
		log:    logger.WithField("s", "controls"),
	}
}

// Update the Controls system
func (c *Controls) Update(dt float64) {

	for _, ev := range c.events {
		for _, e := range c.em.FilteredEntities(components.PosType, components.VelocityType, components.JoystickType) {
			v := c.em.Velocity(e)

			switch val := ev.(type) {
			case events.KeyJustPressed:
				if c.em.HasComponents(e, components.ParentedType) {
					c.log.Debugf("remove parenting from %s", e)
					removeParenting(e, c.em)
				}
				v.Vec = v.Add(gfx.V(0, -5))

			case events.KeyPressed:

				switch val.Key {
				case ebiten.KeyLeft:
					v.Vec = v.Add(gfx.V(-constants.AccelerationX*dt, 0))

				case ebiten.KeyRight:
					v.Vec = v.Add(gfx.V(constants.AccelerationX*dt, 0))
				}
			}
		}
	}
	c.events = []events.Event{}

}

// Send listens for key, mouse and pad-events
func (c *Controls) Send(ev events.Event) {
	switch ev.Type() {
	case events.KeyJustPressedType, events.KeyPressedType:
		c.events = append(c.events, ev)
		c.log.Debugf("recieved %q event", ev.Type())
	default:
		c.log.Debugf("discard %q event", ev.Type())
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
