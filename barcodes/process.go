package barcodes

import (
	"fmt"
	"io"

	"github.com/codenaut/barcoder/images"
	"github.com/codenaut/barcoder/zpl"
	"gopkg.in/urfave/cli.v1"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type BarcodeConfig struct {
	Offset []int
	Unit   string
	Size   []int
	Dpi    int

	Image []ImageConfig
	Qr    []QrConfig
}
type ImageConfig struct {
	File     string
	Size     []int
	Position []int
	Darkness uint16
}
type QrConfig struct {
	Input    int
	Value    string
	Size     []int
	Position []int
}

func scale(factor float64, size []int) (int, int) {
	if len(size) > 1 {
		x := size[0]
		y := size[1]
		return int(factor * float64(x)), int(factor * float64(y))
	}

	return 0, 0
}

func moveCursor(output io.Writer, xOffset, yOffset int, position []int, factor float64) {
	x, y := scale(factor, position)
	zpl.MoveCursor(x+xOffset, y+yOffset, output)
}

func Process(config BarcodeConfig, output io.Writer, args cli.Args) error {
	xOffset := 0
	yOffset := 0
	scaleFactor := 1.0
	if len(config.Offset) > 1 {
		xOffset = config.Offset[0]
		yOffset = config.Offset[1]
	}
	zpl.Start(output)
	for _, img := range config.Image {
		if i, err := images.OpenPng(img.File); err != nil {
			return err
		} else {
			darknessLimit := img.Darkness
			fmt.Println("%#v", img)
			if darknessLimit == 0 {
				darknessLimit = 0xafff
			}
			if len(img.Size) > 1 {
				i = images.Resize(i, img.Size[0], img.Size[1])
			}
			moveCursor(output, xOffset, yOffset, img.Position, scaleFactor)
			flat := images.FlattenImage(i)
			zpl.PutImage(flat, uint16(darknessLimit), output)
		}
	}
	for _, qrConfig := range config.Qr {
		str := qrConfig.Value
		qrCode, err := qr.Encode(str, qr.M, qr.Auto)
		if err != nil {
			return err
		}
		moveCursor(output, xOffset, yOffset, qrConfig.Position, scaleFactor)
		if len(qrConfig.Size) > 1 {
			qrCode, err = barcode.Scale(qrCode, qrConfig.Size[0], qrConfig.Size[1])
		}
		zpl.PutImage(qrCode, 0xfff, output)

	}
	zpl.End(output)
	return nil
}
