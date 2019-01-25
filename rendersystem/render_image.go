package rendersystem

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/peterhellberg/gfx"
)

// RenderImage is responsible for drawing entities to the screen
type RenderImage struct {
	log logging.Logger
	img *ebiten.Image
}

// NewRenderImage creates a new RenderImage system
func NewRenderImage(path string, logger logging.Logger) *RenderImage {
	img, err := gfx.OpenPNG(path)
	if err != nil {
		logger.Fatalf("open image", err)
	}
	eImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		logger.Fatalf("create image", err)
	}

	return &RenderImage{
		log: logger.WithField("s", "RenderImage"),
		img: eImg,
	}
}

// Update the RenderImage system
func (ri *RenderImage) Update(screen *ebiten.Image) {
	screen.DrawImage(ri.img, &ebiten.DrawImageOptions{})
	// drawRect(screen, gfx.R(10, 10, 50, 50))
	// drawRect(screen, gfx.R(50, 50, 300+10, 280+10))
	// drawRect(screen, gfx.R(60, 60, 300+20, 280+20))
}

// Send is an empty method to implement the System interface
func (ri *RenderImage) Send(ev events.Event) {}
