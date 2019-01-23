package system

import (
	"math/rand"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
)

// Input is responsible for handling user input and sending input events
type RandomInput struct {
	em    *entity.Manager
	outCh chan events.Event
	log   logging.Logger
}

// NewRandomInput creates a new RandomInput system
func NewRandomInput(em *entity.Manager, ch chan events.Event, logger logging.Logger) *RandomInput {
	return &RandomInput{
		em:    em,
		outCh: ch,
		log:   logger.WithField("s", "RandomInput"),
	}
}

// Update the RandomInput system
func (i *RandomInput) Update() {

	switch rand.Intn(100) {
	case 1:
		i.outCh <- events.LeftJustPressed{}
	case 2:
		i.outCh <- events.RightJustPressed{}
	default:

	}
}

// Send is an empty method to implement the System interface
func (i *RandomInput) Send(ev events.Event) {}
