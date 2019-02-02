package entity

import (
	"fmt"
	"math/rand"
	"os"
	"text/tabwriter"

	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/animation"
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
	em.entityList = []string{}
	em.typeCounter = make(map[string]int)
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

func (em *Manager) DumpEntity(e string) {
	fmt.Printf("Entity: %s\n", e)

	components, err := em.entities.GetAll(e)
	if err != nil {
		fmt.Println("has no components")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight)
	for k, v := range components {
		fmt.Fprintf(w, "%s\t%v\t\n", k, v)
	}
	w.Flush()
}

func (em *Manager) Add(e string, cs ...interface{}) {
	em.entities.Add(e, cs...)
}

func (em *Manager) Remove(e string, typ components.Type) {
	em.entities.Remove(e, typ)
}

func (em *Manager) RemoveEntity(e string) {
	em.entities.RemoveAll(e)
	removed := []string{}
	for _, f := range em.entityList {
		if e == f {
			continue
		}
		removed = append(removed, f)
	}
	em.entityList = removed
}

func (em *Manager) HasComponents(e string, types ...components.Type) bool {
	return em.entities.HasComponents(e, types...)
}

func uuid() string {
	return fmt.Sprintf("%d", rand.Intn(100000))
}

func (em *Manager) Area(e string) *components.Area {
	return em.entities.GetUnsafe(e, components.AreaType).(*components.Area)
}

func (em *Manager) Get(e string, typ components.Type) (interface{}, error) {
	return em.entities.Get(e, typ)
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

func (em *Manager) Shaking(e string) *components.Shaking {
	return em.entities.GetUnsafe(e, components.ShakingType).(*components.Shaking)
}

func (em *Manager) Drawable(e string) *components.Drawable {
	return em.entities.GetUnsafe(e, components.DrawableType).(*components.Drawable)
}

func (em *Manager) Hitbox(e string) *components.Hitbox {
	return em.entities.GetUnsafe(e, components.HitboxType).(*components.Hitbox)
}

func (em *Manager) Path(e string) *components.Path {
	return em.entities.GetUnsafe(e, components.PathType).(*components.Path)
}

func (em *Manager) OnPath(e string) *components.OnPath {
	return em.entities.GetUnsafe(e, components.OnPathType).(*components.OnPath)
}

func (em *Manager) Parented(e string) *components.Parented {
	return em.entities.GetUnsafe(e, components.ParentedType).(*components.Parented)
}

func (em *Manager) Trigger(e string) *components.Trigger {
	return em.entities.GetUnsafe(e, components.TriggerType).(*components.Trigger)
}

func (em *Manager) ConditionalDrawable(e string) *components.ConditionalDrawable {
	return em.entities.GetUnsafe(e, components.ConditionalDrawableType).(*components.ConditionalDrawable)
}

func (em *Manager) Animation(e string) *animation.Animation {
	return em.entities.GetUnsafe(e, components.AnimationType).(*animation.Animation)
}
