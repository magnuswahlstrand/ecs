package world

import (
	"errors"

	"github.com/peterhellberg/gfx"

	"github.com/kyeett/gomponents/components"
	"github.com/kyeett/gomponents/pathanimation"

	"github.com/kyeett/ecs/blocks"
	"github.com/kyeett/ecs/camera"
	"github.com/kyeett/ecs/player"
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
}

const defaultTimeStep = 1.0

func New(width, height int) *World {
	em := entity.NewManager(logging.NewLogger())

	defaultEntities(em)

	eventCh := make(chan events.Event, 100)
	return &World{
		eventCh: eventCh,
		systems: []system.System{
			system.NewInput(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewRandomInput(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewControls(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewFriction(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewGravity(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			system.NewPath(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewMovement(em, logging.NewLogger(logrus.InfoLevel)),
			system.NewFollow(em, eventCh, logging.NewLogger(logrus.InfoLevel)),
			// system.NewShaking(em, logging.NewLogger(logrus.DebugLevel)),
		},
		renderSystems: []rendersystem.System{
			rendersystem.NewRenderImage("assets/images/background.png", logging.NewLogger()),
			rendersystem.NewRender(em, logging.NewLogger()),
			// rendersystem.NewDebugRender(em, logging.NewLogger()),
		},
		camera: camera.New(em, width, height),
		em:     em,
	}
}

func defaultEntities(em *entity.Manager) {
	// Add a player

	pathID := em.NewEntity("path")
	em.Add(pathID, components.Path{
		Label:  "line",
		Points: gfx.Polygon{gfx.V(10, 100), gfx.V(10, 150), gfx.V(70, 150), gfx.V(70, 100), gfx.V(10, 100)},
		Type:   pathanimation.Polygon,
	})

	blocks.NewDrawable(em, 0, 100, components.OnPath{
		Label:     pathID,
		Speed:     1,
		Target:    1,
		Mode:      pathanimation.LinearLoop,
		Direction: 1,
	})

	pathID = em.NewEntity("path")
	em.Add(pathID, components.Path{
		Label:  "ellipse",
		Points: gfx.Polygon{gfx.V(110, 60), gfx.V(110, 120)},
		Type:   pathanimation.Ellipse,
	})
	blocks.NewDrawable(em, 110, 120, components.OnPath{
		Label:     pathID,
		Speed:     1,
		Target:    1,
		Mode:      pathanimation.LinearPingPong,
		Direction: 1,
	})

	pathID = em.NewEntity("path")
	em.Add(pathID, components.Path{
		Label:  "ellipse",
		Points: gfx.Polygon{gfx.V(170, 80), gfx.V(170, 120)},
		Type:   pathanimation.Ellipse,
	})
	blocks.NewDrawable(em, 110, 120, components.OnPath{
		Label:     pathID,
		Speed:     1,
		Target:    1,
		Mode:      pathanimation.LinearLoop,
		Direction: 1,
	})
	player.NewDrawable(em)
}

func (w *World) Reset() {
	w.em.Reset()
	w.camera.Reset()
	defaultEntities(w.em)
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
