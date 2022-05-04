package UI

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func HandleErrorWindow(err error, w fyne.Window) bool {
	if err != nil {
		dialog.ShowError(err, w)
		return true
	}
	return false
}

func HandleTestError(err error, t *testing.T) {
	if err != nil {
		t.Log(err)
	}
}

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
