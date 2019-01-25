package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/world"
)

const (
	screenWidth  = 400
	screenHeight = 400
)

func main() {

	w := world.New(screenWidth, screenHeight)
	w.StartEventQueue()

	if err := ebiten.Run(w.Update, screenWidth, screenHeight, 1, "ECS demo"); err != nil {
		log.Fatal(err)
	}
}
