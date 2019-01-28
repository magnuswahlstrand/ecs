package system

// A System is updated every iteration, and can receive an Event
type System interface {
	Update(diff float64)
}
