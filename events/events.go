package events

type Event interface {
	Type() EventType
}

type Collision struct {
	entityA, entityB string
}

type LeftJustPressed struct{}
type RightJustPressed struct{}

func (l LeftJustPressed) Type() EventType  { return LeftJustPressedType }
func (l RightJustPressed) Type() EventType { return RightJustPressedType }

type EventType string

const (
	LeftJustPressedType  = "leftJustPressed"
	RightJustPressedType = "rightJustPressed"
)
