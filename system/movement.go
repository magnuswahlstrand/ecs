package system

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Movement is responsible adding velocity to entity positions
type Movement struct {
	em    *entity.Manager
	outCh chan events.Event
	log   logging.Logger
}

// NewMovement creates a new movement system
func NewMovement(em *entity.Manager, ch chan events.Event, logger logging.Logger) *Movement {
	return &Movement{
		em:    em,
		outCh: ch,
		log:   logger.WithField("s", "movement"),
	}
}

// Update the movement system
func (m *Movement) Update(dt float64) {
	for _, e := range m.em.FilteredEntities(components.PosType, components.VelocityType) {
		pos := m.em.Pos(e)
		before := pos
		v := m.em.Velocity(e)
		pos.Vec = pos.Add(v.Vec.Scaled(dt))
		m.log.WithField("id", e).Debugf("%q moving from %s to %s", e, before.Vec, pos.Vec)
	}
}

// Send is an empty method to implement the System interface
func (m *Movement) Send(ev events.Event) {}
