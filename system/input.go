package system

import (
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Input is responsible for handling user input and sending input events
type Input struct {
	em    *entity.Manager
	outCh chan events.Event
	log   logging.Logger
}

// NewInput creates a new input system
func NewInput(em *entity.Manager, ch chan events.Event, logger logging.Logger) *Input {
	return &Input{
		em:    em,
		outCh: ch,
		log:   logger.WithField("s", "input"),
	}
}

// Update the input system
func (i *Input) Update(dt float64) {

	if inpututil.IsKeyJustPressed(ebiten.KeyA) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		i.outCh <- events.LeftJustPressed{}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		i.outCh <- events.RightJustPressed{}
		i.em.Add("camera_1", components.Shaking{time.Now()})

	}
}

// Send is an empty method to implement the System interface
func (i *Input) Send(ev events.Event) {}
