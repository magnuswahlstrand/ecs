package system

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
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

	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		i.outCh <- events.KeyJustPressed{ebiten.KeyUp}
	}

	for _, k := range []ebiten.Key{ebiten.KeyZ, ebiten.KeyX, ebiten.KeyC} {
		if inpututil.IsKeyJustPressed(k) {
			i.outCh <- events.KeyJustPressed{k}
		}
	}

	for _, k := range []ebiten.Key{ebiten.KeyZ, ebiten.KeyX, ebiten.KeyC} {
		if ebiten.IsKeyPressed(k) {
			i.outCh <- events.KeyPressed{k}
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		i.outCh <- events.KeyPressed{ebiten.KeyLeft}
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		i.outCh <- events.KeyPressed{ebiten.KeyRight}
	}
}

// Send is an empty method to implement the System interface
func (i *Input) Send(ev events.Event) {}
