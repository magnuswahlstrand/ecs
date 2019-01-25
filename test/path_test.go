package test

import (
	"fmt"
	"testing"

	"github.com/kyeett/gomponents/pathanimation"

	"github.com/kyeett/ecs/blocks"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/system"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
	"github.com/sirupsen/logrus"
)

func Test_PathPingPongFast(t *testing.T) {
	em := entity.NewManager(logging.NewLogger(logrus.DebugLevel))

	pathID := em.NewEntity("path")
	em.Add(pathID, components.Path{"line", gfx.Polygon{gfx.V(0, 100), gfx.V(100, 100), gfx.V(100, 0)}})
	path := em.Path(pathID)

	// Set up a block that follows a path
	blockID := blocks.New(em, 50, 100, components.OnPath{
		Label:     pathID,
		Speed:     125,
		Target:    1,
		Mode:      pathanimation.LinearPingPong,
		Direction: 1,
	})

	pathSystem := system.NewPath(em, logging.NewLogger(logrus.InfoLevel))
	movSystem := system.NewMovement(em, logging.NewLogger(logrus.InfoLevel))
	pos := em.Pos(blockID)
	onPath := em.OnPath(blockID)

	// Move along this path, speed is high enough, that we should hit every corner
	//    |2
	//    |
	//4___|1,3,5
	//  0

	steps := []struct {
		expectedPos    gfx.Vec
		expectedTarget int
	}{
		{path.Points[1], 2},
		{path.Points[2], 3}, // This is a bit weird, the three will be changed at the beginning of path system update
		{path.Points[1], 0},
		{path.Points[0], -1}, // Same as above
		{path.Points[1], 2},
	}

	for _, s := range steps {
		pathSystem.Update(1)
		movSystem.Update(1)
		if !pos.Vec.Eq(s.expectedPos) {
			t.Fatalf("expected pos=%s, got %s", s.expectedPos, pos.Vec)
		}
		if onPath.Target != s.expectedTarget {
			t.Fatalf("expected target to be %d, got %d", s.expectedTarget, onPath.Target)
		}

	}
}

func Test_PathLoopFast(t *testing.T) {
	em := entity.NewManager(logging.NewLogger(logrus.DebugLevel))

	pathID := em.NewEntity("path")
	em.Add(pathID, components.Path{"line", gfx.Polygon{gfx.V(0, 100), gfx.V(100, 100), gfx.V(100, 0)}})
	path := em.Path(pathID)

	// Set up a block that follows a path
	blockID := blocks.New(em, 50, 100, components.OnPath{
		Label:     pathID,
		Speed:     125,
		Target:    1,
		Mode:      pathanimation.LinearLoop,
		Direction: 1,
	})

	pathSystem := system.NewPath(em, logging.NewLogger(logrus.DebugLevel))
	movSystem := system.NewMovement(em, logging.NewLogger(logrus.DebugLevel))
	pos := em.Pos(blockID)
	onPath := em.OnPath(blockID)

	// Move along this path, speed is high enough, that we should hit every corner
	//    |2,4
	//    |
	// ___|1,3,5
	// 0

	steps := []struct {
		expectedPos    gfx.Vec
		expectedTarget int
	}{
		{path.Points[1], 2},
		{path.Points[2], 3}, // This is a bit weird, the three will be changed at the beginning of path system update
		{path.Points[1], 2},
		{path.Points[2], 3},
		{path.Points[1], 2},
	}

	for _, s := range steps {
		pathSystem.Update(1)
		movSystem.Update(1)
		if !pos.Vec.Eq(s.expectedPos) {
			t.Fatalf("expected pos=%s, got %s", s.expectedPos, pos.Vec)
		}
		if onPath.Target != s.expectedTarget {
			t.Fatalf("expected target to be %d, got %d", s.expectedTarget, onPath.Target)
		}

	}
}

func Test_PathLoopSlow(t *testing.T) {
	em := entity.NewManager(logging.NewLogger(logrus.DebugLevel))

	pathID := em.NewEntity("path")
	em.Add(pathID, components.Path{"line", gfx.Polygon{gfx.V(0, 100), gfx.V(100, 100), gfx.V(100, 0)}})
	path := em.Path(pathID)

	// Set up a block that follows a path
	blockID := blocks.New(em, 0, 100, components.OnPath{
		Label:     pathID,
		Speed:     75,
		Target:    1,
		Mode:      pathanimation.LinearLoop,
		Direction: 1,
	})

	pathSystem := system.NewPath(em, logging.NewLogger(logrus.DebugLevel))
	movSystem := system.NewMovement(em, logging.NewLogger(logrus.DebugLevel))
	pos := em.Pos(blockID)
	onPath := em.OnPath(blockID)

	// Move along this path
	//      |4
	//      |3
	//0,5___|2
	//     1,6
	steps := []struct {
		expectedPos    gfx.Vec
		expectedTarget int
	}{
		{gfx.V(75, 100), 1},
		{path.Points[1], 2}, // This is a bit weird, the three will be changed at the beginning of path system update
		{gfx.V(100, 25), 2},
		{path.Points[2], 3},
		{gfx.V(75, 100), 1},
	}

	for i, s := range steps {
		fmt.Println("TURN!", i)
		pathSystem.Update(1)
		movSystem.Update(1)
		if !pos.Vec.Eq(s.expectedPos) {
			t.Fatalf("expected pos=%s, got %s", s.expectedPos, pos.Vec)
		}
		if onPath.Target != s.expectedTarget {
			t.Fatalf("expected target to be %d, got %d", s.expectedTarget, onPath.Target)
		}

	}
}

func Test_PathSmallSteps(t *testing.T) {
	em := entity.NewManager(logging.NewLogger(logrus.DebugLevel))

	pathID := em.NewEntity("path")
	em.Add(pathID, components.Path{"line", gfx.Polygon{gfx.V(0, 100), gfx.V(100, 100), gfx.V(100, 0)}})
	path := em.Path(pathID)

	// Set up a block that follows a path
	blockID := blocks.New(em, 20, 100, components.OnPath{
		Label:     pathID,
		Speed:     120,
		Target:    1,
		Mode:      pathanimation.LinearLoop,
		Direction: 1,
	})

	pathSystem := system.NewPath(em, logging.NewLogger(logrus.DebugLevel))
	movSystem := system.NewMovement(em, logging.NewLogger(logrus.DebugLevel))
	pos := em.Pos(blockID)
	onPath := em.OnPath(blockID)

	// Move along this path
	//      |4
	//      |3
	//0,5___|2
	//     1,6

	steps := []struct {
		expectedPos    gfx.Vec
		expectedTarget int
	}{
		{gfx.V(80, 100), 1},
		{path.Points[1], 2}, // This is a bit weird, the three will be changed at the beginning of path system update
		{gfx.V(100, 40), 2},
		{path.Points[2], 3},
		{gfx.V(60, 100), 1},
	}

	timeStep := 0.5
	for _, s := range steps {
		pathSystem.Update(timeStep)
		movSystem.Update(timeStep)
		if !pos.Vec.Eq(s.expectedPos) {
			t.Fatalf("expected pos=%s, got %s", s.expectedPos, pos.Vec)
		}
		if onPath.Target != s.expectedTarget {
			t.Fatalf("expected target to be %d, got %d", s.expectedTarget, onPath.Target)
		}
	}
}
