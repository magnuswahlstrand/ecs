package system

import (
	"fmt"
	"math"

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

	collided, possibleMove := checkCollisionY(playerID, m.em)
	pos := m.em.Pos(playerID)
	v := m.em.Velocity(playerID)
	hb := m.em.Hitbox(playerID)
	switch collided {
	case true:
		pos.Y += possibleMove
		v.Y = 0

	default:
		pos.Y += v.Y
		// fmt.Println("no collision for", playerID)
	}

	collided, possibleMove = checkCollisionX(playerID, m.em)
	switch collided {
	case true:
		fmt.Println("hard collision for", playerID)
		fmt.Println(possibleMove, hb.Moved(pos.Vec), pos)
		pos.X += possibleMove
		v.X = 0
	default:
		pos.X += v.X
	}

}

// checkCollisionY uses raycasting in at left, center, and right of hitbox to determine collision
func checkCollisionY(e string, em *entity.Manager) (bool, float64) {
	v := em.Velocity(e)
	if v.Y == 0 {
		return false, -1
	}

	movingUp := v.Y < 0
	sourceHitbox := movedHitbox(e, em)

	var collisions [3]float64
	var rays []gfx.Vec
	if movingUp {
		rays = rayVectors(sourceHitbox.Min.AddXY(0, v.Y), sourceHitbox.W())
		collisions = [3]float64{math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	} else {
		rays = rayVectors(sourceHitbox.Max.AddXY(-sourceHitbox.W(), v.Y), sourceHitbox.W())
		collisions = [3]float64{-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64}
	}

	var hardCollision bool
	for _, t := range em.FilteredEntities(components.HitboxType) {
		if t == e {
			continue
		}
		targetHitbox := movedHitbox(t, em)

		for i, r := range rays {
			if targetHitbox.Contains(r) {
				if movingUp {
					collisions[i] = min(targetHitbox.Max.Y-sourceHitbox.Min.Y+0.1, collisions[i])
					hardCollision = true
				} else {
					collisions[i] = max(targetHitbox.Min.Y-sourceHitbox.Max.Y-0.1, collisions[i])
					hardCollision = true
				}
			}
		}
	}

	var possibleMove float64
	if hardCollision {
		if movingUp {
			possibleMove = min(min(collisions[2], collisions[1]), collisions[0])
		} else {
			possibleMove = max(max(collisions[2], collisions[1]), collisions[0])
		}
	}
	return hardCollision, possibleMove
}

func checkCollisionX(e string, em *entity.Manager) (bool, float64) {
	v := em.Velocity(e)
	if v.X == 0 {
		return false, -1
	}
	sourceHitbox := movedHitbox(e, em).Moved(gfx.V(v.X, 0))

	var hardCollision bool
	zeroRect := gfx.Rect{}
	for _, t := range em.FilteredEntities(components.HitboxType, components.PosType) {
		if t == e {
			continue
		}
		targetHitbox := movedHitbox(t, em)
		intersection := sourceHitbox.Intersect(targetHitbox)
		if intersection != zeroRect {
			hardCollision = true
		}
	}

	var possibleMove float64
	return hardCollision, possibleMove
}

func min(a, b float64) float64 {
	if b < a {
		return b
	}
	return a
}

func max(a, b float64) float64 {
	if b > a {
		return b
	}
	return a
}

// Send is an empty method to implement the System interface
func (m *Movement) Send(ev events.Event) {}
