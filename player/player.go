package player

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

func New(em *entity.Manager) {
	tmp, err := gfx.OpenPNG("assets/images/platformer2.png")
	if err != nil {
		log.Fatal(err)
	}
	pImage, _ := ebiten.NewImageFromImage(tmp, ebiten.FilterDefault)

	e := em.NewEntity("player")
	hitbox := gfx.R(0, 6, 18, 22)
	em.Add(e, components.Pos{Vec: gfx.V(20, 20)})
	em.Add(e, components.Velocity{Vec: gfx.V(0, 0)})
	em.Add(e, components.Joystick{})
	em.Add(e, components.NewHitbox(hitbox))
	em.Add(e, components.Drawable{pImage.SubImage(image.Rect(5, 10, 27, 32)).(*ebiten.Image)})
}
