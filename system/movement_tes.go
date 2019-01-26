package system

import (
	"fmt"

	"github.com/peterhellberg/gfx"

	"github.com/kyeett/ecs/entity"
)

func checkCollisionY(e, t string, em *entity.Manager) (bool, bool) {
	fmt.Printf("checking %s vs %s\n", e, t)

	v := em.Velocity(e)
	if v.Y == 0 {
		return false, false
	}

	sourceHitbox, targetHitbox := movedHitbox(e, em), movedHitbox(t, em)

	var rays []gfx.Vec
	if v.Y < 0 {
		// Moving upwards
		rays = rayVectors(sourceHitbox.Min.AddXY(0, v.Y), sourceHitbox.W())
		fmt.Println("upwards")
	} else {
		// Moving downwards
		rays = rayVectors(sourceHitbox.Max.AddXY(-sourceHitbox.W(), v.Y), sourceHitbox.W())

		fmt.Println("downwards")
	}

	var collisions [3]bool
	for i, r := range rays {
		collisions[i] = targetHitbox.Contains(r)
	}

	var softCollision bool
	hardCollision := collisions[0] || collisions[1] || collisions[2]
	if hardCollision {
		fmt.Printf("%s collided hard with %s\n", e, t)
	}
	return hardCollision, softCollision
}
