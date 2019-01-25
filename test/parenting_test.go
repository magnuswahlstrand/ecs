package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/peterhellberg/gfx"

	"github.com/kyeett/ecs/blocks"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/player"
	"github.com/kyeett/ecs/system"
	"github.com/kyeett/gomponents/components"
	"github.com/sirupsen/logrus"
)

func Test_MovingPlatform(t *testing.T) {
	em := entity.NewManager(logging.NewLogger(logrus.DebugLevel))
	playerID := player.New(em)
	v := em.Velocity(playerID)
	v.Y = 9.0

	blockID := blocks.New(em, 20, 50)
	blockV := em.Velocity(blockID)
	blockV.X = 5.0

	movSystem := system.NewMovement(em, logging.NewLogger(logrus.DebugLevel))
	movSystem.Update(1)
	movSystem.Update(1)

	// Player should have collided, and have block as parent
	if !em.HasComponents(playerID, components.ParentedType) {
		t.Fatal("player should have a parent after collision")
	}

	parented := em.Parented(playerID)
	if parented.ID != blockID {
		t.Fatalf("player should have %q as parent after collision, got %q", blockID, parented.ID)
	}

	playerPos := em.Pos(playerID)
	posBefore := *playerPos

	fmt.Println(playerPos, posBefore)
	movSystem.Update(1)

	// Player velocity should be 0, but position should have changed by same as parent speed
	expectedV := gfx.V(0, 0)
	if !v.Eq(expectedV) {
		log.Fatalf("expected player velocity=%s (still), got %s", expectedV, v)
	}

	expectedPos := posBefore.Add(blockV.Vec)
	if !playerPos.Eq(expectedPos) {
		log.Fatalf("expected player to have moved to %s, got got %s", expectedPos, playerPos)
	}

	fmt.Println(playerPos, posBefore)

	// Player should now have same velocity as parent

	movSystem.Update(1)
	fmt.Println(blockID, playerID)
}
