package test

import (
	"testing"

	"github.com/kyeett/gomponents/direction"

	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/system"
	"github.com/sirupsen/logrus"
)

func Test_GfxRectVecCollision(t *testing.T) {

	r1 := gfx.R(10, 10, 20, 20)
	r2 := gfx.R(5, 5, 8, 8)
	r3 := gfx.R(19, 19, 30, 30)
	v1 := gfx.V(0, 5)
	v2 := gfx.V(15, 15)

	if r1.Contains(v1) {
		t.Fatalf("%s should not contain %s", r1, v1)
	}

	if !r1.Contains(v2) {
		t.Fatalf("%s should contain %s", r1, v2)
	}

	res := r1.Intersect(r2)
	zeroRect := gfx.Rect{}
	if res != zeroRect {
		t.Fatalf("%s and %s should not intersect, got %s", r1, r2, res)
	}

	res = r1.Intersect(r3)
	if res == zeroRect {
		t.Fatalf("%s and %s should not intersect, got %s", r1, r3, res)
	}
}

func Test_Collision(t *testing.T) {
	//  0 1  2  3  4
	//  P       X  P
	//    P  P  P
	//       X  X  X
	//

	em := entity.NewManager(logging.NewLogger(logrus.DebugLevel))
	playerID := em.NewEntity("player")
	initial := gfx.V(0.0, 0.0)
	speed := 10.0
	em.Add(playerID, components.Pos{Vec: initial})
	em.Add(playerID, components.Velocity{Vec: gfx.V(0, speed)})
	em.Add(playerID, components.Hitbox{Rect: gfx.R(0, 0, 10, 10)})

	movSystem := system.NewMovement(em, logging.NewLogger(logrus.DebugLevel))
	movSystem.Update(1)

	pos := em.Pos(playerID)
	expected := initial.AddXY(0, speed)
	if !pos.Eq(expected) {
		t.Fatalf("expected %s, got %s", expected, pos.Vec)
	}

	// Add a box below
	blockID := em.NewEntity("block")
	em.Add(blockID, components.Pos{Vec: gfx.V(0, 20)})
	em.Add(blockID, components.Hitbox{Rect: gfx.R(5, 0, 25, 20)})
	movSystem.Update(1)

	// Should collide with box and not move
	if !pos.Eq(expected) {
		t.Fatalf("expected %s, got %s", expected, pos.Vec)
	}

	// Change direction
	v := em.Velocity(playerID)
	v.Y = -v.Y

	//Add a box above and try to move
	blockID = em.NewEntity("block")
	em.Add(blockID, components.Pos{Vec: gfx.V(0, 0)})
	em.Add(blockID, components.Hitbox{Rect: gfx.R(0, 0, 20, 20).Moved(gfx.V(0, -5))})
	movSystem.Update(1)

	if !pos.Eq(expected) {
		t.Fatalf("expected %s, got %s", expected, pos.Vec)
	}

	// Remove and try again
	em.Remove(blockID, components.HitboxType)
	movSystem.Update(1)
	if !pos.Eq(initial) {
		t.Fatalf("expected %s, got %s", expected, pos.Vec)
	}
}

func Test_OneWayCollision(t *testing.T) {
	// 1. Up, OK,
	// 2. Down Blocked
	//  0 1 2
	//    P P
	//  O O X
	//  P

	em := entity.NewManager(logging.NewLogger(logrus.DebugLevel))
	playerID := em.NewEntity("player")
	initial := gfx.V(0.0, 25.0)
	speed := -5.0
	em.Add(playerID, components.Pos{Vec: initial})
	em.Add(playerID, components.Velocity{Vec: gfx.V(0, speed)})
	em.Add(playerID, components.Hitbox{Rect: gfx.R(0, 0, 10, 5)})

	// Add a box below
	blockID := em.NewEntity("platform")
	em.Add(blockID, components.Pos{Vec: gfx.V(0, 10)})
	em.Add(blockID, components.Hitbox{Rect: gfx.R(0, 0, 10, 10)}, direction.Down)

	movSystem := system.NewMovement(em, logging.NewLogger(logrus.DebugLevel))
	movSystem.Update(3 * 1)

	pos := em.Pos(playerID)
	expected := initial.AddXY(0, 3*speed)
	if !pos.Eq(expected) {
		t.Fatalf("expected %s, got %s", expected, pos.Vec)
	}
}
