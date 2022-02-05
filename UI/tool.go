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

func HandleErrorWithMessage(err error, msg string, w fyne.Window) {
	dialog.ShowError(err, w)
	dialog.ShowInformation("Error", msg, w)
}
