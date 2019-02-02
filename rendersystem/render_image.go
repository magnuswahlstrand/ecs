package rendersystem

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/kyeett/ecs/logging"
	"github.com/peterhellberg/gfx"
)

// RenderImage is responsible for drawing entities to the screen
type RenderImage struct {
	log logging.Logger
	img *ebiten.Image
}

// NewRenderImageFromPath creates a new RenderImage system from a file path
func NewRenderImageFromPath(path string, logger logging.Logger) *RenderImage {
	img, err := gfx.OpenPNG(path)
	if err != nil {
		logger.Fatalf("open image", err)
	}
	eImg, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		logger.Fatalf("create image", err)
	}

	return NewRenderImage(eImg, logger)
}

// NewRenderImage creates a new RenderImage system
func NewRenderImage(img *ebiten.Image, logger logging.Logger) *RenderImage {
	return &RenderImage{
		log: logger.WithField("s", "renderimage"),
		img: img,
	}
}

// Update the RenderImage system
func (ri *RenderImage) Update(screen *ebiten.Image) {
	screen.DrawImage(ri.img, &ebiten.DrawImageOptions{})
}
