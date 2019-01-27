package system

import (
	"strings"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

// Trigger is responsible checking if certain triggers are fullfilled
type Trigger struct {
	em     *entity.Manager
	outCh  chan events.Event
	events []events.Event
	log    logging.Logger
}

// NewTrigger creates a new Trigger system
func NewTrigger(em *entity.Manager, ch chan events.Event, logger logging.Logger) *Trigger {
	return &Trigger{
		em:     em,
		outCh:  ch,
		events: []events.Event{},
		log:    logger.WithField("s", "trigger"),
	}
}

// Check if key pressed is in event list
func (t *Trigger) keyPressed(key string) bool {
	for _, ev := range t.events {
		switch val := ev.(type) {
		case events.KeyJustPressed:
			if val.Key.String() == key {
				return true
			}
		case events.KeyPressed:
			if val.Key.String() == key {
				return true
			}
		}
	}
	return false
}

func (t *Trigger) conditonsMet(cond *components.Condition) bool {
	for i, c := range cond.Conditions {
		switch c[0] {
		case "key_pressed":
			wantedKey := strings.ToUpper(c[1])
			if !t.keyPressed(wantedKey) {
				t.log.Debugf("condition %d not met: %s", i, c)
				return false
			}

		case "in_area":
			params := c[1:]
			r1 := movedHitbox(params[0], t.em)
			r2 := t.em.Area(params[1])
			zeroRect := gfx.Rect{}
			if r1.Intersect(r2.Rect) != zeroRect {
				t.log.Debugf("condition %d met: %s", i, cond)
			} else {
				t.log.Debugf("condition %d not met: %s", i, c)
				return false
			}
		default:
			t.log.Errorf("unknown condition type %s, mark as condition not met", c[0])
			return false
		}
	}

	return true
}

// Update the Trigger system
func (t *Trigger) Update(dt float64) {

	for _, e := range t.em.FilteredEntities(components.ConditionType) {
		cond := t.em.Condition(e)

		if t.conditonsMet(cond) {
			t.log.Debugf("conditions met for %s", cond.Name)
		} else {
			t.log.Debugf("conditions not met for %s", cond.Name)
		}
	}
	t.events = []events.Event{}
}

// Send listens for key, mouse and pad-events
func (t *Trigger) Send(ev events.Event) {
	switch ev.Type() {
	case events.KeyJustPressedType, events.KeyPressedType:
		t.events = append(t.events, ev)
		t.log.Debugf("recieved %q event", ev.Type())
	default:
		t.log.Debugf("discard %q event", ev.Type())
	}
}
