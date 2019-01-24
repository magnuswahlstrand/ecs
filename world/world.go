package world

import (
	"errors"

	"github.com/peterhellberg/gfx"

	"github.com/kyeett/gomponents/components"

	"github.com/kyeett/ecs/blocks"
	"github.com/kyeett/ecs/camera"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/player"
	"github.com/kyeett/ecs/system"
	"github.com/sirupsen/logrus"
)

type World struct {
	eventCh       chan events.Event
	camera        *camera.Camera
	systems       []system.System
	renderSystems []system.RenderSystem
	em            *entity.Manager
}

func New(width, height int) *World {
	em := entity.NewManager(logging.NewLogger())
	// Add a player
	// player.New(em)
	blocks.New(em, 0, 100)
	e := em.NewEntity("path")
	em.Add(e, components.Path{"line", gfx.Polygon{gfx.V(0, 100), gfx.V(100, 100), gfx.V(110, 110)}})

	eventCh := make(chan events.Event, 1000)
	return &World{
		eventCh: eventCh,
		systems: []system.System{
			system.NewInput(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewRandomInput(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewControls(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewFriction(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewGravity(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewPath(em, eventCh, logging.NewLogger(logrus.DebugLevel)),
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
		em:     em,
	}
}

func (w *World) Reset() {
	w.em.Reset()
	w.camera.Reset()
	player.New(w.em)
	blocks.New(w.em, 0, 100)
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

const defaultTimeStep = 0.5

var timeStep = defaultTimeStep

func (w *World) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return errors.New("exit game")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		w.Reset()
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

	r, op := w.camera.View(timeStep)
	screen.DrawImage(w.camera.SubImage(r).(*ebiten.Image), op)
	return nil
}
