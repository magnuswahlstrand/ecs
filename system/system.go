package system

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/events"
)

// A System is updated every iteration, and can receive an Event
type System interface {
	Send(events.Event)
	Update(diff float64)
}

// A RenderSystem is updated every iteration, and draws to a screen
type RenderSystem interface {
	Update(*ebiten.Image)
}
