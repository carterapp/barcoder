package zpl

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"io"
)

func Start(output io.Writer) {
	fmt.Fprint(output, "^XA,^FS\n")
}
func End(output io.Writer) {
	fmt.Fprint(output, "^FS,^XZ\n")
}
func MoveCursor(x int, y int, output io.Writer) {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	fmt.Fprintf(output, "^FO%d,%d\n", x, y)
}
func PutImage(source image.Image, darkness uint16, output io.Writer) {
	size := source.Bounds().Size()
	width := size.X / 8
	height := size.Y
	if size.Y%8 != 0 {
		width = width + 1
	}
	fmt.Fprintf(output, " ^GFA, %d, %d, %d,\n", width*height, width*height, width)

	for y := 0; y < size.Y; y++ {
		line := make([]uint8, width)
		lineIndex := 0
		index := uint8(0)
		currentByte := line[lineIndex]
		for x := 0; x < size.X; x++ {
			index = index + 1
			p := source.At(x, y)
			lum := color.Gray16Model.Convert(p).(color.Gray16)
			if lum.Y < darkness {
				currentByte = currentByte | (1 << (8 - index))
			}
			if index >= 8 {
				line[lineIndex] = currentByte
				lineIndex++
				if lineIndex < len(line) {
					currentByte = line[lineIndex]
				}
				index = 0
			}
		}
		fmt.Fprintln(output, hex.EncodeToString(line))
	}

}
