package images

import (
	"image"
	"image/color"
)

func FlattenImage(source image.Image) *image.NRGBA {
	size := source.Bounds().Size()
	background := color.White
	target := image.NewNRGBA(source.Bounds())
	for y := 0; y < size.Y; y++ {
		for x := 0; x < size.X; x++ {
			p := source.At(x, y)
			flat := flatten(p, background)
			target.Set(x, y, flat)
		}
	}
	return target
}

func flatten(input color.Color, background color.Color) color.Color {
	source := color.NRGBA64Model.Convert(input).(color.NRGBA64)
	r, g, b, a := source.RGBA()
	bg_r, bg_g, bg_b, _ := background.RGBA()
	alpha := float32(a) / 0xffff
	conv := func(c uint32, bg uint32) uint8 {
		val := 0xffff - uint32((float32(bg) * alpha))
		val = val | uint32(float32(c)*alpha)
		return uint8(val >> 8)
	}
	c := color.NRGBA{
		conv(r, bg_r),
		conv(g, bg_g),
		conv(b, bg_b),
		uint8(0xff),
	}
	return c
}
