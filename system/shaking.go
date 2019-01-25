package system

import (
	"time"

	"github.com/hajimehoshi/ebiten"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Shaking is responsible for triggering shaking
type Shaking struct {
	em     *entity.Manager
	events []events.Event
	log    logging.Logger
}

// NewShaking creates a new Shaking system
func NewShaking(em *entity.Manager, logger logging.Logger) *Shaking {
	return &Shaking{
		em:     em,
		events: []events.Event{},
		log:    logger.WithField("s", "shaking"),
	}
}

// Update the Shaking system
func (s *Shaking) Update(dt float64) {
	e := "camera_1"

	// Check for player collision
	for _, ev := range s.events {

		switch val := ev.(type) {
		case events.KeyPressed:
			switch val.Key {
			case ebiten.KeyRight:
				s.em.Add(e, components.Shaking{Started: time.Now()})
			}
		}
	}
	s.events = []events.Event{}
}

// Send listens collision events with player
func (s *Shaking) Send(ev events.Event) {
	switch ev.Type() {
	case events.KeyPressedType:
		s.events = append(s.events, ev)
	default:
	}

}
