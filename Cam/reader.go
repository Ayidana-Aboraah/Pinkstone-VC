package Cam

import (
	"fmt"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
	"image"
)

func ReadImage(img image.Image) *gozxing.Result{
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	reader := oned.NewUPCAReader()

	result, _ := reader.Decode(bmp, nil)

	fmt.Println(result)
	return result
}
