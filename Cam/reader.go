package Cam

import (
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"image"
	_ "image/png"
)


func ReadImage(img image.Image) *gozxing.Result {
	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	barReader := oned.NewUPCAReader()
	result, _ := barReader.Decode(bmp, nil)


	fmt.Println(result)
	return result
}
