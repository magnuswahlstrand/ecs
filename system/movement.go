package system

import (
	"fmt"
	"math"
	"sort"

	"github.com/SolarLune/resolv/resolv"
	"github.com/kyeett/ecs/constants"
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

	collided, softCollision, possibleMove, collidingID := checkCollisionY(playerID, m.em, dt)
	pos := m.em.Pos(playerID)
	v := m.em.Velocity(playerID)
	// fmt.Printf("hard:%t soft:%t\n", collided, softCollision)
	switch collided && !softCollision {
	case true:
		pos.Y += possibleMove
		movingDownward := v.Y > 0
		v.Y = 0

		// If landing on top, mark colliding entity as parent
		fmt.Println(collidingID, movingDownward)

		// Mark colliding as parent!
		m.em.Add(playerID, components.Parented{ID: collidingID})

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
	collided, possibleMove = checkCollisionX(playerID, m.em, dt)
	pos.X += possibleMove
}

type CollisionResult struct {
	ID   string
	move float64
}

// checkCollisionY uses raycasting in at left, center, and right of hitbox to determine collision
func checkCollisionY(e string, em *entity.Manager, dt float64) (bool, bool, float64, string) {
	parentVelocity := gfx.ZV
	if em.HasComponents(e, components.ParentedType) {
		parented := em.Parented(e)
		parentVelocity = em.Velocity(parented.ID).Vec
	}

	v := em.Velocity(e)
	totV := v.Add(parentVelocity).Scaled(dt)
	if totV.Y == 0 {
		return false, false, 0, ""
	}

	movingUp := totV.Y < 0
	sourceHitbox := movedHitbox(e, em)

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
	for _, t := range em.FilteredEntities(components.HitboxType) {
		if t == e {
			continue
		}
		targetHitbox := movedHitbox(t, em)

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
					hb := movedHitbox(collisions[0].ID, em)
					shift := hb.Max.X - rays[0].X

					// Todo, move this outside of the check itself
					if shift < constants.SoftCollisionLimit {
						pos := em.Pos(e)
						pos.X += shift + 0.1 // Add extra shift for robustness on borders
						fmt.Println("shifted!")
						softCollision = true
					}

					fmt.Printf("left is candidate %s %s %0.2f\n", hb, rays[0], shift)
				}

				// Check right side
				if collisions[2].move != math.MaxFloat64 {
					hb := movedHitbox(collisions[2].ID, em)
					shift := hb.Min.X - rays[2].X
					fmt.Printf("right is candidate %s %s %0.2f\n", hb, rays[2], shift)
					if shift > -constants.SoftCollisionLimit {
						pos := em.Pos(e)
						pos.X += shift - 0.1 // Add extra shift for robustness on borders
						fmt.Println("shifted!")
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

func checkCollisionX(e string, em *entity.Manager, dt float64) (bool, float64) {
	parentVelocity := gfx.ZV
	v := em.Velocity(e)
	if em.HasComponents(e, components.ParentedType) {
		parented := em.Parented(e)
		parentVelocity = em.Velocity(parented.ID).Vec
	}

	totV := v.Add(parentVelocity).Scaled(dt)
	if totV.X == 0 {
		return false, 0
	}
	sourceHitbox := movedHitbox(e, em).Moved(gfx.V(totV.X, 0))

	zeroRect := gfx.Rect{}
	for _, t := range em.FilteredEntities(components.HitboxType, components.PosType) {
		if t == e {
			continue
		}
		targetHitbox := movedHitbox(t, em)
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
