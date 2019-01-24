package system

import (
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
			switch ev.Type() {
			case events.LeftJustPressedType:
				v.Vec = v.Add(gfx.V(-2, 0))

			case events.RightJustPressedType:
				v.Vec = v.Add(gfx.V(2, 0))
			}
		}
	}
	c.events = []events.Event{}
}

// Send listens for key, mouse and pad-events
func (c *Controls) Send(ev events.Event) {
	switch ev.Type() {
	case events.LeftJustPressedType, events.RightJustPressedType:
		c.events = append(c.events, ev)
		c.log.Debugf("recieved %q event", ev.Type())
	default:
		c.log.Debugf("discard %q event", ev.Type())
	}
}
