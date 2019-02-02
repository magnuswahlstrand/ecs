package eventsystem

func Add(event interface{}) {
	eventMapper.Handle(entityManager, event)
}

type Collision struct {
	ID1, ID2 string
}
