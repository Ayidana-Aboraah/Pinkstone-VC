package Cam

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	_ "github.com/pion/mediadevices/pkg/driver/camera"
	"image"
)

var camRaw image.Image

func OpenCam() string{
	camRaw = canvas.NewImageFromFile("Assets/icon02.png").Image
	CamOutput = canvas.NewImageFromImage(camRaw)
	CamOutput.FillMode = canvas.ImageFillOriginal
	CamOutput.Refresh()

	w := fyne.CurrentApp().NewWindow("Camera")
	w.SetContent(CamOutput)
	w.Show()

	label := widget.NewLabel(" ")
	label.SetText(StartCamera())
	w.Close()
	return label.Text
}

func MakeCamMenu() string{
return ""
}