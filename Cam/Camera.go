package Cam

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
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

	var complete, evacuate bool
	w.SetOnClosed(func() {
		if !complete {
			evacuate = true
			done <- true
			return
		}
	})

	label.SetText(StartCamera(&CamOutput, done))
	complete = true

	if evacuate {
		fmt.Println("Evacuating...")
		return "V"
	}

	w.Close()
	return label.Text
}
