package system

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/eventsystem"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
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

func (t *Trigger) conditonsMet(trigger *components.Trigger) bool {
	for _, condition := range trigger.Conditions {
		if !condition.IsMet() {
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

	for _, e := range t.em.FilteredEntities(components.TriggerType) {
		trig := t.em.Trigger(e)

		if t.conditonsMet(trig) {
			t.log.Debugf("conditions met for %q", trig.Name)
			t.updateConditionalDrawable(true, trig.Name)
			eventsystem.UpdateEventChains(e)
		} else {
			t.log.Debugf("conditions not met for %q", trig.Name)
			t.updateConditionalDrawable(false, trig.Name)
		}
	}
}
