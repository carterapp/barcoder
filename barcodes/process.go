package barcodes

import (
	"fmt"
	"io"
	"os"

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
	Dpmm   int

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
	Center   bool
}

func scale(factor float64, size []int) (int, int) {
	if len(size) > 1 {
		x := float64(size[0])
		y := float64(size[1])
		return int(factor * x), int(factor * y)
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

	scaleFactor := float64(1.0)
	if config.Dpmm > 0 {
		scaleFactor = float64(config.Dpmm)
	}
	if len(config.Offset) > 1 {
		xOffset = config.Offset[0]
		yOffset = config.Offset[1]
	}
	xCenter := 0
	yCenter := 0
	if len(config.Size) > 1 {
		xCenter = int(float64(config.Size[0]) * scaleFactor / 2)
		yCenter = int(float64(config.Size[1]) * scaleFactor / 2)

	}
	zpl.Start(output)
	for _, img := range config.Image {
		if i, err := images.OpenPng(img.File); err != nil {
			return err
		} else {
			darknessLimit := img.Darkness
			if darknessLimit == 0 {
				darknessLimit = 0xafff
			}
			if len(img.Size) > 1 {
				x, y := scale(scaleFactor, img.Size)
				i = images.Resize(i, x, y)
			}
			moveCursor(output, xOffset, yOffset, img.Position, scaleFactor)
			flat := images.FlattenImage(i)
			zpl.PutImage(flat, uint16(darknessLimit), output)
		}
	}
	for _, qrConfig := range config.Qr {
		str := qrConfig.Value
		if str == "" {
			str = args.Get(qrConfig.Input)
		}
		qrCode, err := qr.Encode(str, qr.M, qr.Auto)
		if err != nil {
			return err
		}
		if len(qrConfig.Size) > 1 {
			x, y := scale(scaleFactor, qrConfig.Size)
			if x > 0 && y > 0 {
				qrCode, err = barcode.Scale(qrCode, x, y)
				if err != nil {
					return err
				}
			}
		}
		if qrConfig.Center {
			imgWidth := float64(qrCode.Bounds().Size().X)
			imgHeight := float64(qrCode.Bounds().Size().Y)
			x := xCenter - int(imgWidth/2) + xOffset
			y := yCenter - int(imgHeight/2) + yOffset
			fmt.Fprint(os.Stderr, "%d, %d", xCenter, yCenter)
			posx, posy := scale(scaleFactor, qrConfig.Position)
			zpl.MoveCursor(posx+x, posy+y, output)
		} else {
			moveCursor(output, xOffset, yOffset, qrConfig.Position, scaleFactor)
		}

		zpl.PutImage(qrCode, 0xfff, output)

	}
	zpl.End(output)
	return nil
}
