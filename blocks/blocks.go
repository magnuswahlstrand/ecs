package blocks

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/peterhellberg/gfx"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/gomponents/components"
	"github.com/kyeett/gomponents/direction"
)

func New(em *entity.Manager, x, y int, cs ...interface{}) string {

	e := em.NewEntity("block")
	hitbox := gfx.R(0, 0, 64, 64)
	em.Add(e, components.Pos{Vec: gfx.IV(x, y)})
	em.Add(e, components.Velocity{Vec: gfx.V(0, 0)})
	em.Add(e, components.NewHitbox(hitbox))
	em.Add(e, cs...)
	return e
}

func NewDrawable(em *entity.Manager, x, y int, cs ...interface{}) {
	tmp, err := gfx.OpenPNG("assets/images/box.png")
	if err != nil {
		log.Fatal(err)
	}
	pImage, _ := ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)

	cs = append(cs, components.Drawable{pImage})
	New(em, x, y, cs...)
}

func NewOneWay(em *entity.Manager, x, y int, cs ...interface{}) {
	e := em.NewEntity("platform")
	hitbox := gfx.R(0, 0, 32, 10)
	em.Add(e, components.Pos{Vec: gfx.IV(x, y)})
	em.Add(e, components.Velocity{Vec: gfx.V(0, 0)})
	em.Add(e, components.Hitbox{Rect: hitbox, BlockedDirections: direction.Up})
	em.Add(e, cs...)
}

func NewOneWayDrawable(em *entity.Manager, x, y int, cs ...interface{}) {
	tmp, err := gfx.OpenPNG("assets/images/platform.png")
	if err != nil {
		log.Fatal(err)
	}
	pImage, _ := ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)
	cs = append(cs, components.Drawable{pImage})
	NewOneWay(em, x, y, cs...)
}
