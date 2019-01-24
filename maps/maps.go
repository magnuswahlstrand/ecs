package maps

import (
	"fmt"
	"log"

	tiled "github.com/kyeett/go-tiled"
)

func New(path string) {
	m, err := tiled.LoadFromFile("assets/maps/polygon.tmx")
	if err != nil {
		log.Fatal("load map ", err)
	}
	fmt.Printf("%#v\n", m.ObjectGroups[0])
}
