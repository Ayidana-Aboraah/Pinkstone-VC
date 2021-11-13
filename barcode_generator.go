package main

import (
	"image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

func main() {
	// Generate a barcode image (*BitMatrix)
	enc := oned.NewCode128Writer()
	img, _ := enc.Encode("1", gozxing.BarcodeFormat_CODE_128, 250, 50, nil)

	file, _ := os.Create("barcodes [test]/barcode.png")
	defer file.Close()

	// *BitMatrix implements the image.Image interface,
	// so it is able to be passed to png.Encode directly.
	_ = png.Encode(file, img)
}