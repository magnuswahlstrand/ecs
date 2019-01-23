package system

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
)

// Camera is responsible for keeping track of the camera modifiers
type Camera struct {
	em     *entity.Manager
	events []events.Event
	log    logging.Logger
}

// NewCamera creates a new Camera system
func NewCamera(em *entity.Manager, logger logging.Logger) *Camera {
	return &Camera{
		em:     em,
		events: []events.Event{},
		log:    logger.WithField("s", "camera"),
	}
}

// Update the Camera system
func (c *Camera) Update(screen *ebiten.Image) {
	// e := "camera_1"

	// c.events = []events.Event{}

	// shakeDuration := 2 * time.Second
	// shakeMagnitude := 15.0

	// if c.em.HasComponents(e, components.ShakingType) {
	// 	s := c.em.Shaking(e)
	// 	since := time.Since(s.Started)
	// 	if since < shakeDuration {
	// 		decay := math.Pow(math.E, -float64(10*float64(since)/float64(shakeDuration)))
	// 		shakeMagnitude *= decay
	// 		offset := gfx.V(rand.Float64(), rand.Float64()).Unit().Sub(gfx.V(0.5, 0.5)).Scaled(shakeMagnitude)
	// 		pos := c.em.Pos(e)
	// 		pos.Vec = pos.Add(offset)
	// 		fmt.Println(pos)
	// 	} else {
	// 		//Todo remove shake after
	// 	}

	// 	// screen.DrawImage(screen, &ebiten.DrawImageOptions{})
	// }
}

// Send listens collision events with player
func (c *Camera) Send(ev events.Event) {}
