package Cam

import (
	"github.com/makiuchi-d/gozxing"
	"gocv.io/x/gocv"
)

func OpenCam() (*gozxing.Result, error, string){
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil{
		return nil, err, "You're camera May not be connected"
	}

	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		imageObject, _ := img.ToImage()
		results := ReadImage(imageObject)

		if results != nil{
			e := webcam.Close()
			if e != nil{
				return results, e, "Your Camera seems to not want to close."
			}
			er := window.Close()
			if er != nil{
				return results, er, "You're camera window doesn't want to close"
			}
			return results, err, "All good"
		}

		window.WaitKey(3)
	}
}