package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/world"
)

const (
	screenWidth  = 200
	screenHeight = 200
)

func main() {

	w := world.New(screenWidth, screenHeight)
	w.StartEventQueue()

	if err := ebiten.Run(w.Update, screenWidth, screenHeight, 2, "ECS demo"); err != nil {
		log.Fatal(err)
	}
}
