package rendersystem

import (
	"fmt"
	"image/color"

	"github.com/kyeett/gomponents/pathanimation"

	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/peterhellberg/gfx"
	"golang.org/x/image/colornames"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// DebugRender is responsible for drawing debug info on screen
type DebugRender struct {
	em          *entity.Manager
	log         logging.Logger
	blackscreen *ebiten.Image
}

// NewDebugRender creates a new DebugRender system
func NewDebugRender(em *entity.Manager, logger logging.Logger) *DebugRender {
	img, _ := ebiten.NewImage(200, 200, ebiten.FilterDefault)
	img.Fill(color.Black)
	return &DebugRender{
		em:          em,
		log:         logger.WithField("s", "DebugRender"),
		blackscreen: img,
	}
}

const (
	mask        = true
	hitboxes    = true
	spriteboxes = false
	paths       = false
	rays        = false
	text        = false
)

// Update the DebugRender system
func (r *DebugRender) Update(screen *ebiten.Image) {

	if mask {
		r.drawBlackmask(screen)
	}

	if spriteboxes {
		r.drawSpriteboxes(screen)
	}

	if hitboxes {
		r.drawHitboxes(screen)
	}

	if text {
		r.drawText(screen)
	}

	if paths {
		r.drawPaths(screen)
	}

	if rays {
		r.drawRays(screen)
	}
}

func (r *DebugRender) drawBlackmask(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.5)
	screen.DrawImage(r.blackscreen, op)
}

func (r *DebugRender) drawHitboxes(screen *ebiten.Image) {
	for _, e := range r.em.FilteredEntities(components.PosType, components.DrawableType) {
		pos := r.em.Pos(e)
		if r.em.HasComponents(e, components.HitboxType) {
			hb := r.em.Hitbox(e)
			drawRect(screen, hb.Rect.Moved(pos.Vec), colornames.Red)
		}
	}
}

func (r *DebugRender) drawSpriteboxes(screen *ebiten.Image) {
	for _, e := range r.em.FilteredEntities(components.PosType, components.DrawableType) {
		pos := r.em.Pos(e)
		s := r.em.Drawable(e)

		imgRect := gfx.IR(0, 0, s.Bounds().Dx(), s.Bounds().Dy())
		drawRect(screen, gfx.BoundsToRect(imgRect).Moved(pos.Vec), colornames.Greenyellow)
	}
}

func (r *DebugRender) drawText(screen *ebiten.Image) {
	for _, e := range r.em.FilteredEntities(components.PosType, components.VelocityType) {
		pos := r.em.Pos(e)
		v := r.em.Velocity(e)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("xV: %0.2f", v.X), pos.Pt().X, pos.Pt().Y-15)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("yV: %0.2f", v.Y), pos.Pt().X, pos.Pt().Y)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("p: (%0.0f,%0.0f)", pos.X, pos.Y), pos.Pt().X, pos.Pt().Y+15)
	}
}

func (r *DebugRender) drawPaths(screen *ebiten.Image) {
	for _, e := range r.em.FilteredEntities(components.PathType) {
		path := r.em.Path(e)
		switch path.Type {
		case pathanimation.Polygon:
			for i := range path.Points[:len(path.Points)-1] {
				ebitenutil.DrawLine(screen, path.Points[i].X, path.Points[i].Y, path.Points[i+1].X, path.Points[i+1].Y, colornames.Red)
			}
		case pathanimation.Ellipse:
			gfx.DrawCircle(screen, path.Points[0], path.Points[0].Sub(path.Points[1]).Len(), 1.1, colornames.Purple)
			gfx.DrawLine(screen, path.Points[0], path.Points[1], 1.1, colornames.Yellow)
		}
	}
}

func (r *DebugRender) drawRays(screen *ebiten.Image) {
	for _, e := range r.em.FilteredEntities(components.JoystickType, components.HitboxType, components.PosType) {
		pos := r.em.Pos(e)
		hb := r.em.Hitbox(e)

		upperLeft := hb.Rect.Moved(pos.Vec).Min
		bottomRight := hb.Rect.Moved(pos.Vec).Max

		for _, v := range []gfx.Vec{upperLeft, upperLeft.Add(gfx.V(hb.W()/2, 0)), upperLeft.Add(gfx.V(hb.W(), 0))} {
			gfx.DrawLine(screen, v, v.Add(gfx.V(0, -10)), 1.01, colornames.Yellow)
		}
		for _, v := range []gfx.Vec{bottomRight.Sub(gfx.V(hb.W(), 0)), bottomRight.Sub(gfx.V(hb.W()/2, 0)), bottomRight} {
			gfx.DrawLine(screen, v, v.Add(gfx.V(0, 10)), 1.01, colornames.Yellow)
		}
	}
}

// Send is an empty method to implement the System interface
func (r *DebugRender) Send(ev events.Event) {}
