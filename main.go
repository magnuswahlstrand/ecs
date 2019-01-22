package main

import (
	"github.com/kyeett/ecs/world"
)

func main() {
	w := world.New()
	w.StartEventQueue()
	for i := 0; i < 3; i++ {
		w.Update()
	}
}
