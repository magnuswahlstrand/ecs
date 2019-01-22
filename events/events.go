package events

type Event interface {
	Type() EventType
}

type Collision struct {
	entityA, entityB string
}

type LeftPressed struct{}

func (l LeftPressed) Type() EventType {
	return LeftPressedType
}

func (l RightPressed) Type() EventType {
	return RightPressedType
}

type RightPressed struct{}

type EventType string

const (
	LeftPressedType  = "leftPressed"
	RightPressedType = "rightPressed"
)
