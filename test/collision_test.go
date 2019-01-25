package test

import (
	"fmt"
	"testing"

	"github.com/kyeett/ecs/blocks"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/ecs/player"
	"github.com/kyeett/ecs/system"
	"github.com/sirupsen/logrus"
)

func Test_Collision(t *testing.T) {
	em := entity.NewManager(logging.NewLogger(logrus.DebugLevel))
	playerID := player.New(em)
	v := em.Velocity(playerID)
	v.Y = 9.0

	blockID := blocks.New(em, 20, 50)

	movSystem := system.NewMovement(em, logging.NewLogger(logrus.DebugLevel))
	movSystem.Update(1)
	movSystem.Update(1)
	movSystem.Update(1)
	movSystem.Update(1)
	fmt.Println(blockID, playerID)
}
