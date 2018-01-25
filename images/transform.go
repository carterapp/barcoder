package images

import (
	"image"

	"github.com/disintegration/imaging"
)

func Resize(img image.Image, width, height int) *image.NRGBA {
	return imaging.Resize(img, width, height, imaging.Gaussian)
}
