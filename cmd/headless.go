package main

import "github.com/kyeett/ecs/world"

const (
	screenWidth  = 200
	screenHeight = 200
)

func main() {

	w := world.New(screenWidth, screenHeight)
	w.StartEventQueue()
}
