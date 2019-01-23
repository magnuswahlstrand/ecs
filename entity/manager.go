package entity

import (
	"fmt"
	"math/rand"

	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

type Manager struct {
	entityList  []string
	entities    *components.Map
	log         logging.Logger
	typeCounter map[string]int
}

func NewManager(logger logging.Logger) *Manager {
	return &Manager{
		entityList:  []string{},
		entities:    components.NewMap(),
		log:         logger.WithField("s", "entity-mgr"),
		typeCounter: make(map[string]int),
	}
}

func (em *Manager) Reset() {
	for _, e := range em.entityList {
		em.entities.RemoveAll(e)
	}
}

func (em *Manager) FilteredEntities(types ...components.Type) []string {
	var IDs []string
	for _, ID := range em.entityList {
		if em.entities.HasComponents(ID, types...) {
			IDs = append(IDs, ID)
		}
	}
	return IDs
}

func (em *Manager) NewEntity(typ ...string) string {
	id := uuid()
	if len(typ) > 0 {
		em.typeCounter[typ[0]] += 1
		id = fmt.Sprintf("%s_%d", typ[0], em.typeCounter[typ[0]])
	}
	em.entityList = append(em.entityList, id)
	em.log.Debugf("new entity '%s' created", id)
	return id
}

func (em *Manager) Add(e string, cs ...interface{}) {
	em.entities.Add(e, cs...)
}

func uuid() string {
	return fmt.Sprintf("%d", rand.Intn(10000))
}

func (em *Manager) Pos(e string) *components.Pos {
	return em.entities.GetUnsafe(e, components.PosType).(*components.Pos)
}

func (em *Manager) Velocity(e string) *components.Velocity {
	return em.entities.GetUnsafe(e, components.VelocityType).(*components.Velocity)
}

func (em *Manager) Follow(e string) *components.Following {
	return em.entities.GetUnsafe(e, components.FollowingType).(*components.Following)
}
