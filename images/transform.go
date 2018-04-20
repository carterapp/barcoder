package images

import (
	"image"

	"github.com/disintegration/imaging"
)

func Resize(img image.Image, width, height int) *image.NRGBA {
	return imaging.Resize(img, width, height, imaging.Gaussian)
}

func Rotate(img image.Image, angle float64) *image.NRGBA {
	return imaging.Rotate(img, angle, nil)
}
