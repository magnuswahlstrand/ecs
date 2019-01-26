package system

import (
	"fmt"
	"math"

	"github.com/peterhellberg/gfx"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/gomponents/components"
)

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
		fmt.Printf("checking %s vs %s\n", e, t)
		targetHitbox := movedHitbox(t, em)

		for i, r := range rays {
			if targetHitbox.Contains(r) {
				// Collision
				if movingUp {
					h := em.Hitbox(e).H()
					collisions[i] = min(targetHitbox.Max.Y-sourceHitbox.Min.Y+0.1, collisions[i])
					hardCollision = true
					fmt.Println(v.Y, sourceHitbox, sourceHitbox.Size(), collisions[i], targetHitbox)
					fmt.Println(targetHitbox.Max.Y + h + 0.1)
					// time.Sleep(10 * time.Second)
				} else {
					// Remove player hitbox
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
	fmt.Println(collisions)

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
