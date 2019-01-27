package world

import (
	"fmt"
	"log"

	"github.com/kyeett/ecs/blocks"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/player"
	"github.com/kyeett/gomponents/components"
	"github.com/kyeett/gomponents/pathanimation"
	"github.com/peterhellberg/gfx"
)

func (w *World) populateWorld() {
	switch w.mapName {
	case "default":
		defaultEntities(w.em)
	case "single_block":
		singleBlock(w.em)
	case "horizontalplatform":
		horizontalPlatform(w.em)
	case "horizontalcollision":
		horizontalCollision(w.em)
	case "softcollision":
		softCollision(w.em)
	default:
		log.Fatal("not a valid world")
	}
}

func defaultEntities(em *entity.Manager) {
	// Add a player
	pathID := em.NewEntity("path")
	em.Add(pathID, components.Path{
		Label:  "line",
		Points: gfx.Polygon{gfx.V(10, 100), gfx.V(10, 150), gfx.V(70, 150), gfx.V(70, 100), gfx.V(10, 100)},
		Type:   pathanimation.Polygon,
	})

	blocks.NewDrawable(em, 0, 100, components.OnPath{
		Label:     pathID,
		Speed:     1,
		Target:    1,
		Mode:      pathanimation.LinearLoop,
		Direction: 1,
	})

	pathID = em.NewEntity("path")
	em.Add(pathID, components.Path{
		Label:  "ellipse",
		Points: gfx.Polygon{gfx.V(110, 60), gfx.V(110, 120)},
		Type:   pathanimation.Ellipse,
	})
	blocks.NewDrawable(em, 110, 120, components.OnPath{
		Label:     pathID,
		Speed:     1,
		Target:    1,
		Mode:      pathanimation.LinearLoop,
		Direction: 1,
	})

	player.NewDrawable(em, 100, 20)
}

func horizontalPlatform(em *entity.Manager) {
	// Add a player
	pathID := em.NewEntity("path")
	em.Add(pathID, components.Path{
		Label:  "line",
		Points: gfx.Polygon{gfx.V(10, 120), gfx.V(300, 120)},
		Type:   pathanimation.Polygon,
	})

	blocks.NewDrawable(em, 10, 110, components.OnPath{
		Label:     pathID,
		Speed:     1,
		Target:    1,
		Mode:      pathanimation.LinearLoop,
		Direction: 1,
	})

	blocks.NewDrawable(em, 110, 40)
	player.NewDrawable(em, 40, 40)
	fmt.Println(em.Velocity("player_1"))
}

func singleBlock(em *entity.Manager) {
	blocks.NewDrawable(em, 10, 120)
	player.NewDrawable(em, 70, 80)
}

func softCollision(em *entity.Manager) {
	blocks.NewDrawable(em, 10, 120)
	blocks.NewDrawable(em, 10-48-3+3, 120-64-30)
	blocks.NewDrawable(em, 10+48-3+3, 120-64-30)
	player.NewDrawable(em, 20+19+3, 93) //+19
	// player.NewDrawable(em, 20, 93)  // Left side
}

func horizontalCollision(em *entity.Manager) {
	blocks.NewDrawable(em, 10, 120)
	blocks.NewDrawable(em, 10+64, 80)
	blocks.NewDrawable(em, 10-64, 80)
	player.NewDrawable(em, 40, 80)
}
