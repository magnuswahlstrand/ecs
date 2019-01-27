package world

import (
	"errors"

	"github.com/kyeett/ecs/camera"
	"github.com/kyeett/ecs/constants"
	"github.com/kyeett/ecs/rendersystem"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/system"
	"github.com/sirupsen/logrus"
)

type World struct {
	eventCh       chan events.Event
	camera        *camera.Camera
	systems       []system.System
	renderSystems []rendersystem.System
	em            *entity.Manager
	mapName       string
}

func New(m string, width, height int) *World {
	em := entity.NewManager(logging.NewLogger())
	eventCh := make(chan events.Event, 100)
	w := World{
		eventCh: eventCh,
		systems: []system.System{
			system.NewInput(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewRandomInput(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewFriction(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewControls(em, eventCh, logging.NewLogger(logrus.DebugLevel)),
			system.NewGravity(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewPath(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewParenting(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewMovement(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewFollow(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewShaking(em, logging.NewLogger(logrus.DebugLevel)),
		},
		renderSystems: []rendersystem.System{
			rendersystem.NewRenderImage("assets/images/background.png", logging.NewLogger()),
			rendersystem.NewRender(em, logging.NewLogger()),
			// rendersystem.NewDebugRender(em, logging.NewLogger()),
		},
		camera:  camera.New(em, width, height),
		em:      em,
		mapName: m,
	}
	w.populateWorld()
	return &w
}

func (w *World) Reset() {
	w.em.Reset()
	w.camera.Reset()
	w.populateWorld()
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

var timeStep = constants.DefaultTimeStep

func (w *World) Update(screen *ebiten.Image) error {
	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return errors.New("exit game")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		w.Reset()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		timeStep = constants.DefaultTimeStep - timeStep
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
