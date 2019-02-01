package system

import (
	"log"
	"strings"
	"time"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/inputhandler"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Trigger is responsible checking if certain triggers are fullfilled
type Trigger struct {
	em          *entity.Manager
	eventChains []EventChain
	log         logging.Logger
}

// NewTrigger creates a new Trigger system
func NewTrigger(em *entity.Manager, logger logging.Logger, eventChains ...EventChain) *Trigger {
	return &Trigger{
		em:          em,
		log:         logger.WithField("s", "trigger"),
		eventChains: eventChains,
	}
}

type EventChain interface {
	Next(*entity.Manager)
	WaitingForID() string
	Done() bool
}

func (t *Trigger) conditonsMet(cond *components.Condition) bool {
	for _, c := range cond.Conditions {
		switch c[0] {
		case "key_pressed":
			wantedKey := strings.ToUpper(c[1])
			if !inputhandler.KeyPressed(wantedKey) {
				return false
			}

		case "in_area":
			params := c[1:]
			r1 := movedHitbox(params[0], t.em)
			r2 := t.em.Area(params[1])

			if !r1.Overlaps(r2.Rect) {
				return false
			}
		case "wait_until":
			t, err := time.Parse(time.RFC3339Nano, c[1])
			if err != nil {
				log.Fatal("incorrect time format", c[1])
			}
			return time.Now().After(t)
		case "animation_complete":
			for _, e := range t.em.FilteredEntities(components.AnimationType) {
				if e == c[1] {
					return false
				}
			}
			// No such animation
			return true

		default:
			t.log.Errorf("unknown condition type %s, mark as condition not met", c[0])
			return false
		}
	}

	return true
}

func (t *Trigger) updateConditionalDrawable(b bool, conditionName string) {
	// Find entities that depend on this condition
	for _, f := range t.em.FilteredEntities(components.ConditionalDrawableType) {
		if t.em.ConditionalDrawable(f).ConditionName == conditionName {

			cd := t.em.ConditionalDrawable(f)
			noLimit := cd.MaxTransitions == 0
			belowLimit := cd.Transitions < cd.MaxTransitions
			if cd.ConditionMet != b && (noLimit || belowLimit) {
				cd.ConditionMet = b
				cd.Transitions++
			}
		}
	}
}

func (t *Trigger) updateEventChains(conditionID string) {

	var stoppedEventChains bool
	for _, ec := range t.eventChains {
		if ec.WaitingForID() == conditionID {
			ec.Next(t.em)
			if ec.Done() {
				stoppedEventChains = true
			}
		}
	}

	// Remove stopped chains
	if stoppedEventChains {
		var tmp []EventChain
		for _, ec := range t.eventChains {
			if ec.Done() {
				continue
			}
			tmp = append(tmp, ec)
		}
		t.eventChains = tmp
	}
}

// Update the Trigger system
func (t *Trigger) Update(dt float64) {

	for _, e := range t.em.FilteredEntities(components.ConditionType) {
		cond := t.em.Condition(e)

		if t.conditonsMet(cond) {
			t.log.Debugf("conditions met for %s", cond.Name)
			t.updateConditionalDrawable(true, cond.Name)
			t.updateEventChains(e)
		} else {
			t.log.Debugf("conditions not met for %s", cond.Name)
			t.updateConditionalDrawable(false, cond.Name)
		}
	}
}
