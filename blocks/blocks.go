package blocks

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

func New(em *entity.Manager, x, y int, cs ...interface{}) string {
	// tmp, err := gfx.OpenPNG("assets/images/box.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pImage, _ := ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)

	e := em.NewEntity("block")
	hitbox := gfx.R(0, 0, 64, 64)
	em.Add(e, components.Pos{Vec: gfx.IV(x, y)})
	em.Add(e, components.Velocity{Vec: gfx.V(0, 0)})
	em.Add(e, components.NewHitbox(hitbox))
	// em.Add(e, components.Drawable{pImage})
	em.Add(e, cs...)
	return e
}
