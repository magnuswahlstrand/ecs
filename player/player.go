package player

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

func New(em *entity.Manager) {
	e := em.NewEntity("player")
	em.Add(e, components.Pos{Vec: gfx.V(0, 0)})
	em.Add(e, components.Velocity{Vec: gfx.V(0, 0)})
	em.Add(e, components.Following{ID: e})
	em.Add(e, components.Joystick{})
}
