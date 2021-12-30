package Cam

import (
	"github.com/makiuchi-d/gozxing"
	"gocv.io/x/gocv"
)

func OpenCam() *gozxing.Result{
	webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		imageObject, _ := img.ToImage()
		results := ReadImage(imageObject)

		if results != nil{
			webcam.Close()
			window.Close()
			return results
		}

		window.WaitKey(1)
	}
}
