package Cam

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	_ "github.com/pion/mediadevices/pkg/driver/camera"
)

func OpenCam() string {
	CamOutput := canvas.Image{}
	CamOutput.FillMode = canvas.ImageFillOriginal
	CamOutput.Refresh()

	w := fyne.CurrentApp().NewWindow("Camera")
	w.SetContent(&CamOutput)
	w.Show()

	label := widget.NewLabel(" ")

	done := make(chan bool)
	defer close(done)

	evacuate := false
	w.SetOnClosed(func() {
		done <- true
		evacuate = true
		label.SetText("V")
	})

	label.SetText(StartCamera(&CamOutput, done))

	if evacuate {
		return "V"
	}

	w.Close()
	return label.Text
}
