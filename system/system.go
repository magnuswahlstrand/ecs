package system

import (
	"github.com/kyeett/ecs/events"
)

// A System is updated every iteration, and can receive an Event
type System interface {
	Send(events.Event)
	Update(diff float64)
}
