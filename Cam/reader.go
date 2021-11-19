package Cam

import (
	"fmt"
	"github.com/makiuchi-d/gozxing/oned"
	"image"
	_ "image/png"
	"os"

	"github.com/makiuchi-d/gozxing"
)

var Path = "barcodes [test]/"

func main() {
	// open and decode image file
	//path := "barcode.png"
	file, _ := os.Open(Path + "barcode.png")
	img, _, _ := image.Decode(file)
	ReadImage(img)
}

func ReadImage(img image.Image) *gozxing.Result{
	// prepare BinaryBitmap
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)

	// decode image
	barReader := oned.NewUPCAReader()
	result, _ := barReader.Decode(bmp, nil)

	fmt.Println(result)
	return result
}

func AddNewStoreItem(){
	
}