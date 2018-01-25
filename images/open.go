package images

import (
	"image"
	"image/png"
	"os"
)

func OpenPng(filename string) (image.Image, error) {
	if f, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		return png.Decode(f)
	}

}
