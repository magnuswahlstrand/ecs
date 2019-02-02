package eventsystem

import (
	"log"

	"github.com/kyeett/ecs/entity"
)

// Todo: integrate package in a better way

var eventChains []EventChain
var entityManager *entity.Manager
var eventMapper EventMapper

type EventMapper interface {
	Handle(*entity.Manager, interface{})
}

type EventChain interface {
	Next(*entity.Manager)
	WaitingForID() string
	Done() bool
}

func InitializeEventQueue(em *entity.Manager, m EventMapper) {
	entityManager = em
	eventMapper = m
}

func clean() {
	var tmp []EventChain
	for _, ec := range eventChains {
		if ec.Done() {
			continue
		}
		tmp = append(tmp, ec)
	}
	eventChains = tmp
}

func Reset() {
	eventChains = []EventChain{}
}

func AddEventChain(ec EventChain) {
	eventChains = append(eventChains, ec)
}

func UpdateEventChains(conditionID string) {
	var stoppedEventChains bool
	for _, ec := range eventChains {
		if ec.WaitingForID() == conditionID {
			ec.Next(entityManager)
			if ec.Done() {
				stoppedEventChains = true
			}
		}
	}

	// Remove stopped chains
	if stoppedEventChains {
		clean()
	}
}

type Basic struct {
	CurrentStep  int
	waitingForID string
	done         bool
	ComponentStack
}

type ComponentStack struct {
	stack []interface{}
}

func (cs *ComponentStack) Push(c interface{}) {
	cs.stack = append(cs.stack, c)
}

func (cs *ComponentStack) Pop() interface{} {
	c := cs.stack[len(cs.stack)-1]
	cs.stack = cs.stack[:len(cs.stack)-1]
	return c
}

func (pd *Basic) Next(em *entity.Manager) bool {
	log.Fatal("not implemented")
	return true
}

func (pd *Basic) Done() bool {
	return pd.done
}

func (pd *Basic) SetDone(b bool) {
	pd.done = b
}

func (pd *Basic) WaitingForID() string {
	return pd.waitingForID
}

func (pd *Basic) SetWaitingForID(ID string) {
	pd.waitingForID = ID
}
