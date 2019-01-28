package system

import (
	"strings"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/inputhandler"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

// Trigger is responsible checking if certain triggers are fullfilled
type Trigger struct {
	em  *entity.Manager
	log logging.Logger
}

// NewTrigger creates a new Trigger system
func NewTrigger(em *entity.Manager, logger logging.Logger) *Trigger {
	return &Trigger{
		em:  em,
		log: logger.WithField("s", "trigger"),
	}
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
			zeroRect := gfx.Rect{}
			if r1.Intersect(r2.Rect) != zeroRect {
			} else {
				return false
			}
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

// Update the Trigger system
func (t *Trigger) Update(dt float64) {

	for _, e := range t.em.FilteredEntities(components.ConditionType) {
		cond := t.em.Condition(e)

		if t.conditonsMet(cond) {
			t.log.Debugf("conditions met for %s", cond.Name)
			t.updateConditionalDrawable(true, cond.Name)

		} else {
			t.log.Debugf("conditions not met for %s", cond.Name)
			t.updateConditionalDrawable(false, cond.Name)
		}
	}
}
