package world

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/player"
	"github.com/kyeett/ecs/system"
)

type World struct {
	eventCh chan events.Event
	systems []system.System
}

func New() *World {
	em := entity.NewManager(logging.NewLogger())
	player.New(em)
	eventCh := make(chan events.Event)
	return &World{
		eventCh: eventCh,
		systems: []system.System{
			system.NewInput(em, eventCh, logging.NewLogger()),
			system.NewGravity(em, eventCh, logging.NewLogger()),
			system.NewMovement(em, eventCh, logging.NewLogger()),
			system.NewFollow(em, eventCh, logging.NewLogger()),
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
