package main

import (
	"flag"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/world"
)

const (
	screenWidth  = 200
	screenHeight = 200
)

func main() {

	var mapName string
	flag.StringVar(&mapName, "map", "default", "which map to play")
	flag.Parse()

	w := world.New(mapName, screenWidth, screenHeight)
	w.StartEventQueue()

	if err := ebiten.Run(w.Update, screenWidth, screenHeight, 2, "ECS demo"); err != nil {
		log.Fatal(err)
	}
}
