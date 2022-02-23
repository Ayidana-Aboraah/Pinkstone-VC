package Cam

import (
	"fmt"
	"image"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/oned"
)

func ReadImage(img image.Image) *gozxing.Result {
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	reader := oned.NewUPCAReader()

	result, _ := reader.Decode(bmp, nil)
	if result == nil {
		europeReader := oned.NewEAN13Reader()
		result, _ = europeReader.Decode(bmp, nil)
	}

	fmt.Println(result)
	return result
}
