package system

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

// Render is responsible for drawing entities to the screen
type Render struct {
	em  *entity.Manager
	log logging.Logger
}

// NewRender creates a new Render system
func NewRender(em *entity.Manager, logger logging.Logger) *Render {
	return &Render{
		em:  em,
		log: logger.WithField("s", "Render"),
	}
}

// Update the Render system
func (r *Render) Update(screen *ebiten.Image) {
	for _, e := range r.em.FilteredEntities(components.PosType, components.DrawableType) {
		fmt.Println(e)
		pos := r.em.Pos(e)

		sImg := r.em.Drawable(e).Image
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pos.X, pos.Y)
		screen.DrawImage(sImg, op)

	}
}

// Send is an empty method to implement the System interface
func (r *Render) Send(ev events.Event) {}

func drawRect(screen *ebiten.Image, r gfx.Rect, c color.Color) {
	pts := []gfx.Vec{r.Min, r.Min.Add(gfx.V(0, r.H())), r.Max, r.Min.Add(gfx.V(r.W(), 0))}
	for j := range pts {
		p1, p2 := pts[j], pts[(j+1)%len(pts)]
		ebitenutil.DrawLine(screen, p1.X, p1.Y, p2.X, p2.Y, c)
	}
}
