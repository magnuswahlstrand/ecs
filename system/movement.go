package system

import (
	"math"
	"sort"

	"github.com/SolarLune/resolv/resolv"
	"github.com/kyeett/gomponents/direction"
	"github.com/peterhellberg/gfx"

	"github.com/kyeett/ecs/constants"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
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
	if !m.em.HasComponents(playerID, components.VelocityType) {
		return
	}

	collided, softCollision, possibleMove, collidingID := m.checkCollisionY(playerID, dt)
	pos := m.em.Pos(playerID)
	v := m.em.Velocity(playerID)
	switch collided && !softCollision {
	case true:
		pos.Y += possibleMove
		v.Y = 0

		// Mark colliding as parent!
		if m.em.HasComponents(collidingID, components.VelocityType) {
			m.em.Add(playerID, components.Parented{ID: collidingID})
		}

	default:
		// Todo, handle this in a nicer way
		parentVelocity := gfx.ZV
		if m.em.HasComponents(playerID, components.ParentedType) {
			parented := m.em.Parented(playerID)
			parentVelocity = m.em.Velocity(parented.ID).Vec
		}
		pos.Y += v.Add(parentVelocity).Y * dt
	}

	// Horizontal check
	collided, possibleMove = m.checkCollisionX(playerID, dt)
	if abs(possibleMove) < 0.1 {
		v.X = 0
	}
	pos.X += possibleMove
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

type CollisionResult struct {
	ID   string
	move float64
}

// checkCollisionY uses raycasting in at left, center, and right of hitbox to determine collision
func (m *Movement) checkCollisionY(e string, dt float64) (bool, bool, float64, string) {
	parentVelocity := gfx.ZV
	if m.em.HasComponents(e, components.ParentedType) {
		parented := m.em.Parented(e)
		if m.em.HasComponents(parented.ID, components.VelocityType) {
			parentVelocity = m.em.Velocity(parented.ID).Vec
		}
	}

	v := m.em.Velocity(e)
	totV := v.Add(parentVelocity).Scaled(dt)
	if totV.Y == 0 {
		return false, false, 0, ""
	}

	movingUp := totV.Y < 0
	sourceHitbox := movedHitbox(e, m.em)

	var collisions []CollisionResult
	var rays []gfx.Vec
	if movingUp {
		rays = rayVectors(sourceHitbox.Min.AddXY(0, totV.Y), sourceHitbox.W())
		collisions = []CollisionResult{{"", math.MaxFloat64}, {"", math.MaxFloat64}, {"", math.MaxFloat64}}
	} else {
		rays = rayVectors(sourceHitbox.Max.AddXY(-sourceHitbox.W(), totV.Y), sourceHitbox.W())
		collisions = []CollisionResult{{"", -math.MaxFloat64}, {"", -math.MaxFloat64}, {"", -math.MaxFloat64}}
	}

	var hardCollision, softCollision bool
	var collisionCount int
	for _, t := range m.em.FilteredEntities(components.HitboxType) {
		if t == e {
			continue
		}

		targetHitbox := movedHitbox(t, m.em)
		hitboxBlocksDown := (m.em.Hitbox(t).BlockedDirections & direction.Down) != 0
		if movingUp && !hitboxBlocksDown {
			m.log.Debugf("ignoring %s, not blocking %s", t, direction.Down)
			continue
		}

		for i, r := range rays {
			if targetHitbox.Contains(r) {
				collisionCount++
				if movingUp {
					y := targetHitbox.Max.Y - sourceHitbox.Min.Y + 0.1
					if y < collisions[i].move {
						collisions[i].move = y
						collisions[i].ID = t
					}
					hardCollision = true
				} else {
					y := targetHitbox.Min.Y - sourceHitbox.Max.Y - 0.1
					if y > collisions[i].move {
						collisions[i].move = y
						collisions[i].ID = t
					}
					hardCollision = true
				}
			}
		}
	}

	if hardCollision {
		if movingUp {

			// Check for collision being soft
			// Inspired by this article https://forums.tigsource.com/index.php?topic=46289.msg1387138#msg1387138
			if collisionCount == 1 {
				// Check left side
				if collisions[0].move != math.MaxFloat64 {
					hb := movedHitbox(collisions[0].ID, m.em)
					shift := hb.Max.X - rays[0].X
					if shift < constants.SoftCollisionLimit {
						pos := m.em.Pos(e)
						pos.X += shift + 0.1 // Add extra shift for robustness on borders
						softCollision = true
					}
				}

				// Check right side
				if collisions[2].move != math.MaxFloat64 {
					hb := movedHitbox(collisions[2].ID, m.em)
					shift := hb.Min.X - rays[2].X
					if shift > -constants.SoftCollisionLimit {
						pos := m.em.Pos(e)
						pos.X += shift - 0.1 // Add extra shift for robustness on borders
						softCollision = true
					}
				}
			}

			sort.Slice(collisions, func(i int, j int) bool {
				return collisions[i].move < collisions[j].move
			})

		} else {
			sort.Slice(collisions, func(i int, j int) bool {
				return collisions[i].move > collisions[j].move
			})
		}
	}
	return hardCollision, softCollision, collisions[0].move, collisions[0].ID
}

func (m *Movement) checkCollisionX(e string, dt float64) (bool, float64) {
	parentVelocity := gfx.ZV
	v := m.em.Velocity(e)
	if m.em.HasComponents(e, components.ParentedType) {
		parented := m.em.Parented(e)
		if m.em.HasComponents(parented.ID, components.VelocityType) {
			parentVelocity = m.em.Velocity(parented.ID).Vec
		}
	}

	totV := v.Add(parentVelocity).Scaled(dt)
	if totV.X == 0 {
		return false, 0
	}
	sourceHitbox := movedHitbox(e, m.em).Moved(gfx.V(totV.X, 0))

	zeroRect := gfx.Rect{}
	for _, t := range m.em.FilteredEntities(components.HitboxType, components.PosType) {
		if t == e {
			continue
		}
		movingLeft := v.X < 0
		movingRight := v.X > 0

		hitboxBlocksLeft := (m.em.Hitbox(t).BlockedDirections & direction.Left) != 0
		hitboxBlocksRight := (m.em.Hitbox(t).BlockedDirections & direction.Right) != 0
		if movingLeft && !hitboxBlocksLeft {
			m.log.Debugf("ignoring %s, not directed %s", t, direction.Left)
			continue
		}

		if movingRight && !hitboxBlocksRight {
			m.log.Debugf("ignoring %s, not directed %s", t, direction.Right)
			continue
		}

		targetHitbox := movedHitbox(t, m.em)
		intersection := sourceHitbox.Intersect(targetHitbox)
		if intersection != zeroRect {
			return true, 0
		}
	}

	return false, totV.X * 0.99
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
