package events

import (
	"github.com/hajimehoshi/ebiten"
)

type Event interface {
	Type() EventType
}

type Collision struct {
	entityA, entityB string
}

type KeyPressed struct{ ebiten.Key }
type KeyJustPressed struct{ ebiten.Key }

func (kp KeyPressed) Type() EventType     { return KeyPressedType }
func (kp KeyJustPressed) Type() EventType { return KeyJustPressedType }

type EventType string

const (
	KeyJustPressedType = "keyJustPressed"
	KeyPressedType     = "keyPressed"
)
