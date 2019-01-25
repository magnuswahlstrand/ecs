package rendersystem

import (
	"fmt"
	"image/color"

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

// Update the DebugRender system
func (r *DebugRender) Update(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.5)
	screen.DrawImage(r.blackscreen, op)

	for _, e := range r.em.FilteredEntities(components.PosType, components.DrawableType) {
		pos := r.em.Pos(e)
		s := r.em.Drawable(e)

		// gfx.DrawCircle(screen, pos.Vec, 5, 3, colornames.Green)
		// gfx.DrawPointCircle(screen, pos.Pt(), 2, 0.5, colornames.Purple)
		imgRect := gfx.IR(0, 0, s.Bounds().Dx(), s.Bounds().Dy())

		drawRect(screen, gfx.BoundsToRect(imgRect).Moved(pos.Vec), colornames.Greenyellow)
		if r.em.HasComponents(e, components.HitboxType) {
			hb := r.em.Hitbox(e)
			drawRect(screen, hb.Rect.Moved(pos.Vec), colornames.Orange)
		}

		if r.em.HasComponents(e, components.VelocityType) {
			v := r.em.Velocity(e)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("xV: %0.2f", v.X), pos.Pt().X, pos.Pt().Y-15)
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("yV: %0.2f", v.Y), pos.Pt().X, pos.Pt().Y)
		}
	}

	for _, e := range r.em.FilteredEntities(components.PathType) {
		path := r.em.Path(e)
		for i := range path.Points[:len(path.Points)-1] {
			ebitenutil.DrawLine(screen, path.Points[i].X, path.Points[i].Y, path.Points[i+1].X, path.Points[i+1].Y, colornames.Red)
		}
	}
}

// Send is an empty method to implement the System interface
func (r *DebugRender) Send(ev events.Event) {}
