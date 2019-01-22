package world

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/system"
)

type World struct {
	eventCh chan events.Event
	systems []system.System
}

func New() *World {
	em := entity.NewManager()
	eventCh := make(chan events.Event)
	return &World{
		eventCh: eventCh,
		systems: []system.System{
			system.NewInput(em, eventCh, logging.NewSystemLogger()),
			system.NewGravity(em, eventCh, logging.NewSystemLogger()),
			system.NewMovement(em, eventCh, logging.NewSystemLogger()),
		},
	}
}

func (w *World) StartEventQueue() {
	go func() {
		for {
			ev := <-w.eventCh
			for _, s := range w.systems {
				s.Send(ev)
			}
		}
	}()
}

func (w *World) Update() {
	for _, s := range w.systems {
		s.Update()
	}
}
