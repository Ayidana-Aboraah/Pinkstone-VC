package Cam

import (
	"image"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

func ReadImage(img image.Image) *gozxing.Result {
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	reader := oned.NewUPCAReader()

	result, _ := reader.Decode(bmp, nil)
	return result
}
