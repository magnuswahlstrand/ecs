package events

type Event interface {
	Type() EventType
}

type Collision struct {
	entityA, entityB string
}

type LeftJustPressed struct{}
type RightJustPressed struct{}
type UpJustPressed struct{}

func (l LeftJustPressed) Type() EventType  { return LeftJustPressedType }
func (l UpJustPressed) Type() EventType    { return UpJustPressedType }
func (l RightJustPressed) Type() EventType { return RightJustPressedType }

type EventType string

const (
	UpJustPressedType    = "upJustPressed"
	LeftJustPressedType  = "leftJustPressed"
	RightJustPressedType = "UpJustPressed"
)
