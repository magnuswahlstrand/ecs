package system

import (
	"fmt"

	"github.com/SolarLune/resolv/resolv"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

// Movement is responsible adding velocity to entity positions
type Movement struct {
	em  *entity.Manager
	log logging.Logger
}

// NewMovement creates a new movement system
func NewMovement(em *entity.Manager, logger logging.Logger) *Movement {
	return &Movement{
		em:  em,
		log: logger.WithField("s", "movement"),
	}
}

func resolvRectangle(r gfx.Rect) *resolv.Rectangle {
	return resolv.NewRectangle(int32(r.Min.X), int32(r.Min.Y), int32(r.W()), int32(r.H()))
}

// Update the movement system
func (m *Movement) Update(dt float64) {
	playerID := "player_1"

	// Create space
	var space resolv.Space
	for _, e := range m.em.FilteredEntities(components.HitboxType) {
		if e == playerID {
			continue
		}
		pos := m.em.Pos(e)
		hb := m.em.Hitbox(e)
		hbMoved := hb.Moved(pos.Vec)
		fmt.Println(hbMoved)
		r := resolvRectangle(hbMoved)
		r.SetTags(e)
		space.AddShape(r)
	}

	pos := m.em.Pos(playerID)    //Todo, remove hardcoding?
	v := m.em.Velocity(playerID) //Todo, remove hardcoding?
	hb := m.em.Hitbox(playerID)

	// Round to whole int steps
	tX, tY := pos.Add(v.Vec.Scaled(dt)).XY()
	rX, rY := int32(tX), int32(tY)

	pX, pY := pos.XY()
	rPX, rPY := int32(pX), int32(pY)
	fmt.Println(pX, pY, rX-rPX, rY-rPY, v.Y)

	r := resolvRectangle(hb.Moved(pos.Vec))

	// if rY-rPY != 0 {
	if res := space.Resolve(r, 0, rY-rPY); res.Colliding() && !res.Teleporting {
		fmt.Println("Colliding", res.ShapeB.GetTags())
		v.Y = 0
	} else {
		fmt.Println("OK!", res.Colliding(), res.Teleporting)
		pos.Y += v.Y * dt
		r.Move(0, rY-rPY) //FIXME, is this correct?
	}
	// }

	if rX-rPX != 0 {
		if res := space.Resolve(r, rX-rPX, 0); res.Colliding() && !res.Teleporting {
			fmt.Println("Colliding with", res.ShapeB.GetTags())
			v.X = 0
		} else {
			fmt.Println("OK!", res.Colliding(), res.Teleporting)
			pos.X += v.X * dt
		}
	}
	// r.Move(x int32, y int32)
	// if res := space.Resolve(r, 0, int32(v.Y)); res.Colliding() && !res.Teleporting {
	// 	v.Y = 0
	// } else {

	// }

	for _, e := range m.em.FilteredEntities(components.PosType, components.VelocityType) {
		if e == playerID {
			continue
		}
		pos := m.em.Pos(e)
		before := *pos
		v := m.em.Velocity(e)
		pos.Vec = pos.Add(v.Vec.Scaled(dt))
		m.log.WithField("id", e).Debugf("%q moving from %s to %s", e, before.Vec, pos.Vec)
	}
}

// Send is an empty method to implement the System interface
func (m *Movement) Send(ev events.Event) {}
