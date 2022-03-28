package UI

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func HandleErrorWindow(err error, w fyne.Window) bool {
	if err != nil {
		dialog.ShowError(err, w)
		return true
	}
	return false
}
