package Cam

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func OpenCam(origin *fyne.Window) int {
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

	w.Close()

	if evacuate {
		fmt.Println("Evacuating...")
		return 0
	}

	if label.Text == "X" {
		dialog.ShowInformation("Time Up!", "The camera has been open for too long, but you can open it again.", *origin)
		return 0
	}

	conID, _ := strconv.Atoi(label.Text)
	return conID
}
