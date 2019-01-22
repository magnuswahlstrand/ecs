package system

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
)

// Movement is responsible adding velocity to entity positions
type Movement struct {
	em    *entity.Manager
	outCh chan events.Event
	log   logging.Logger
}

// NewMovement creates a new movement system
func NewMovement(em *entity.Manager, ch chan events.Event, logger logging.Logger) *Input {
	return &Input{
		em:    em,
		outCh: ch,
		log:   logger.WithField("system", "movement"),
	}
}

// Update the movement system
func (m *Movement) Update() {}

// Send is an empty method to implement the System interface
func (m *Movement) Send(ev events.Event) {}
