package entity

import "github.com/kyeett/gomponents/components"

type Manager struct {
	entityList []string
	entities   *components.Map
}

func NewManager() *Manager {
	return &Manager{
		entityList: []string{},
		entities:   components.NewMap(),
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

func (em *Manager) Pos(e string) *components.Pos {
	return em.entities.GetUnsafe(e, components.PosType).(*components.Pos)
}
