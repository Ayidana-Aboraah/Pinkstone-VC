package Cam

import (
	"gocv.io/x/gocv"
)

func OpenCam() string{
	webcam, _ := gocv.OpenVideoCapture(0)
	window := gocv.NewWindow("Bronze Hermes")
	img := gocv.NewMat()

	for {
		webcam.Read(&img)
		window.IMShow(img)
		//Turning Image mat into a normal image
		imgObj,_ := img.ToImage()
		//[Note:] Maybe use the converted image as an image in the GUI instead(if better for performance)

		//Reading the new Image
		res := ReadImage(imgObj)

		if res != nil{
			window.Close()
			webcam.Close()
			return res.String()
		}

		window.WaitKey(1)
	}
}

func StartCamera(){
	webcam, _ := gocv.OpenVideoCapture(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()
	for  {
		webcam.Read(&img)
		window.IMShow(img)
		window.WaitKey(1)
	}
}

