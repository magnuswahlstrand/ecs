package camera

import (
	"image"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

type Camera struct {
	*ebiten.Image
	em   *entity.Manager
	rect gfx.Rect
}

const padding = 0.0

func New(em *entity.Manager, width, height int) *Camera {
	initializeEntity(em)
	// em.Add(e, components.Following{ID: "player_1"})

	img, err := ebiten.NewImage(width+int(padding*2), height+int(padding*2), ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	return &Camera{
		Image: img,
		em:    em,
		rect:  gfx.R(padding, padding, padding+float64(img.Bounds().Dx()), padding+float64(img.Bounds().Dy())),
	}
}

func initializeEntity(em *entity.Manager) {
	e := em.NewEntity("camera")
	em.Add(e, components.Pos{Vec: gfx.V(0, 0)})
}

func (c *Camera) Reset() {
	initializeEntity(c.em)
}

func (c *Camera) View(dt float64) (image.Rectangle, *ebiten.DrawImageOptions) {
	op := &ebiten.DrawImageOptions{}
	e := "camera_1"
	// pos := c.em.Pos("camera_1")

	shakeDuration := 250 * time.Millisecond
	speed := 1.0
	shakeMagnitude := 10.0

	var offset gfx.Vec
	pos := c.em.Pos(e)
	pos.Y = 0
	if c.em.HasComponents(e, components.ShakingType) {
		s := c.em.Shaking(e)
		since := float64(time.Since(s.Started))
		normalizeShakeDuration := float64(shakeDuration)
		if since < normalizeShakeDuration {
			decay := math.Pow(math.E, -float64(speed*since/normalizeShakeDuration))
			shakeMagnitude *= decay
			offset = gfx.V(rand.Float64(), rand.Float64()).Unit().Sub(gfx.V(0.5, 0.5)).Scaled(shakeMagnitude)
			// pos.Vec = pos.Add(offset)
		} else {
			c.em.Remove(e, components.ShakingType)
			//Todo remove shake after
		}

		// screen.DrawImage(screen, &ebiten.DrawImageOptions{})
	}
	// op.GeoM.Translate(+padding, +padding)
	// op.ColorM.Scale(0.5, 1, 0.5, 0.8)

	// return c.Image.Bounds(), op
	return c.rect.Moved(pos.Vec).Moved(offset).Bounds(), op
}
