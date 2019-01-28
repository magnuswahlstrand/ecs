package world

import (
	"errors"

	"github.com/kyeett/ecs/camera"
	"github.com/kyeett/ecs/constants"
	"github.com/kyeett/ecs/rendersystem"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/system"
	"github.com/sirupsen/logrus"
)

type World struct {
	camera        *camera.Camera
	systems       []system.System
	renderSystems []rendersystem.System
	em            *entity.Manager
	mapName       string
}

func New(m string, width, height int) *World {
	em := entity.NewManager(logging.NewLogger())
	w := World{
		systems: []system.System{
			system.NewFriction(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewControls(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewGravity(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewParenting(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewPath(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewMovement(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewFollow(em, logging.NewLogger(logrus.InfoLevel)),
		},
		renderSystems: []rendersystem.System{
			rendersystem.NewRenderImageFromPath("assets/images/background.png", logging.NewLogger()),
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
