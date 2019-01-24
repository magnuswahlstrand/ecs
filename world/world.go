package world

import (
	"github.com/kyeett/ecs/camera"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/player"
	"github.com/kyeett/ecs/system"
	"github.com/peterhellberg/gfx"
	"github.com/sirupsen/logrus"
)

type World struct {
	eventCh       chan events.Event
	camera        *camera.Camera
	systems       []system.System
	renderSystems []system.RenderSystem
}

func New(width, height int) *World {
	em := entity.NewManager(logging.NewLogger())
	// Add a player
	player.New(em)

	eventCh := make(chan events.Event, 1000)
	return &World{
		eventCh: eventCh,
		systems: []system.System{
			system.NewInput(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewRandomInput(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewControls(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewFriction(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewGravity(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewMovement(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewFollow(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewShaking(em, logging.NewLogger(logrus.DebugLevel)),
			// system.NewCollision(em, eventCh, logging.NewLogger(logrus.DebugLevel)),
		},
		renderSystems: []system.RenderSystem{
			system.NewRenderImage("assets/images/background.png", logging.NewLogger()),
			system.NewRender(em, logging.NewLogger()),
			system.NewDebugRender(em, logging.NewLogger()),
			// system.NewCamera(em, logging.NewLogger()),
		},
		camera: camera.New(em, width, height),
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

const defaultTimeStep = 0.1

var timeStep = defaultTimeStep

func (w *World) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return gfx.ErrDone
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		timeStep = defaultTimeStep - timeStep
	}

	for _, s := range w.systems {
		s.Update(timeStep)
	}

	for _, s := range w.renderSystems {
		s.Update(w.camera.Image)
	}

	r, op := w.camera.View()
	screen.DrawImage(w.camera.SubImage(r).(*ebiten.Image), op)
	return nil
}
