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

func resolvLine(v1, v2 gfx.Vec) *resolv.Line {
	return resolv.NewLine(int32(v1.X), int32(v1.Y), int32(v1.X), int32(v1.Y))
}

// Update the movement system
func (m *Movement) Update(dt float64) {
	if dt == 0.0 {
		return
	}
	playerID := "player_1"

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
	m.movePlayer(dt)
}

func (m *Movement) movePlayer(dt float64) {

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
		r := resolvRectangle(hbMoved)
		r.SetTags(e)
		space.AddShape(r)
	}

	pos := m.em.Pos(playerID)    //Todo, remove hardcoding?
	v := m.em.Velocity(playerID) //Todo, remove hardcoding?
	hb := m.em.Hitbox(playerID)

	parentVelocity := gfx.ZV
	if m.em.HasComponents(playerID, components.ParentedType) {
		parented := m.em.Parented(playerID)
		parentVelocity = m.em.Velocity(parented.ID).Vec
		if parentVelocity.Y < -0.5 {
			// time.Sleep(1 * time.Second)
		}
	}

	// Round to whole int steps
	tX, tY := pos.Add(v.Vec.Add(parentVelocity).Scaled(dt)).XY()
	rX, rY := int32(tX), int32(tY)

	pX, pY := pos.XY()
	rPX, rPY := int32(pX), int32(pY)

	r := resolvRectangle(hb.Moved(pos.Vec))

	// if rY-rPY != 0 {
	if res := space.Resolve(r, 0, rY-rPY); res.Colliding() && !res.Teleporting {
		// fmt.Println("Y Colliding", res.ShapeB.GetTags(), v.Y, res.ResolveY)

		collidingOnTop := v.Y > 0

		v.Y = 0

		// If landing on top, mark colliding entity as parent
		if collidingOnTop {
			collidingID := res.ShapeB.GetTags()[0]
			hbColl := m.em.Hitbox(collidingID).Moved(m.em.Pos(collidingID).Vec)
			cV := m.em.Velocity(collidingID)

			// Set pos to upper limit
			fmt.Printf("add colliding, resolv=%v, playerhb=%v, collidedHB=%v\n", res.ResolveY, hb.Moved(pos.Vec), hbColl)
			// time.Sleep(15 * time.Second)
			// pos.Y += float64(res.ResolveY)
			pos.Y += hbColl.Min.Sub(hb.Moved(pos.Vec).Max).Y

			// Mark colliding as parent!
			m.em.Add(playerID, components.Parented{ID: collidingID})
			v.X -= cV.X
		}

	} else {
		// fmt.Println("Y OK!", res.Colliding(), res.Teleporting)
		pos.Y += (v.Y + parentVelocity.Y) * dt
		r.Move(0, rY-rPY) //FIXME, is this correct?
	}

	if res := space.Resolve(r, rX-rPX, 0); res.Colliding() && !res.Teleporting {
		// fmt.Println("X Colliding with", res.ShapeB.GetTags())
		v.X = 0
	} else {
		pos.X += (v.X + parentVelocity.X) * dt
		// fmt.Println("X OK!", res.Colliding(), res.Teleporting)
	}
}

// Send is an empty method to implement the System interface
func (m *Movement) Send(ev events.Event) {}
