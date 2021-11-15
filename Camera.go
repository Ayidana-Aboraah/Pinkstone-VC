package main

import (
	"github.com/makiuchi-d/gozxing"
	"gocv.io/x/gocv"
)

var res *gozxing.Result

func OpenCam() {
	webcam, _ := gocv.OpenVideoCapture(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)

		//Turning Image mat into a normal image
		imgObj,_ := img.ToImage()
		//[Note:] Maybe use the converted image as an image in the GUI instead(if better for performance)

		//Reading the new Image
		res = ReadImage(imgObj)

		window.WaitKey(1)
	}
}