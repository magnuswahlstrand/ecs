package system

import (
	"time"

	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Animation is responsible for updating animations
type Animation struct {
	em        *entity.Manager
	log       logging.Logger
	lastFrame time.Time
}

// NewAnimation creates a new Animation system
func NewAnimation(em *entity.Manager, logger logging.Logger) *Animation {

	return &Animation{
		em:        em,
		log:       logger.WithField("s", "animation"),
		lastFrame: time.Now(),
	}
}

// Update the animation  system
func (a *Animation) Update(dt float64) {
	for _, e := range a.em.FilteredEntities(components.AnimationType) {

		animation := a.em.Animation(e)
		if animation.IsLastFrame {
			a.em.RemoveEntity(e)
			continue
		}
		animation.Update(dt * float64(time.Since(a.lastFrame).Nanoseconds()) / 1000000)
	}
	a.lastFrame = time.Now()
}
