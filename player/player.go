package player

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

func New(em *entity.Manager, x, y float64, cs ...interface{}) string {
	e := em.NewEntity("player")
	hitbox := gfx.R(0, 0, 18, 22)
	em.Add(e, components.Pos{Vec: gfx.V(x, y)})
	em.Add(e, components.Velocity{Vec: gfx.V(0, 0)})
	em.Add(e, components.Joystick{})
	em.Add(e, components.NewHitbox(hitbox))
	em.Add(e, cs...)
	return e
}

func NewDrawable(em *entity.Manager, x, y float64, cs ...interface{}) {
	tmp, err := gfx.OpenPNG("assets/images/platformer2.png")
	if err != nil {
		log.Fatal(err)
	}
	pImage, _ := ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)

	cs = append(cs, components.Drawable{pImage.SubImage(image.Rect(5, 10, 27, 32)).(*ebiten.Image)})
	New(em, x, y, cs...)
}
