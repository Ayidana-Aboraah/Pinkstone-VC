package Debug

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var knownErrors = []string{
	"Invalid Input\nCheck your input",
	"Pieces answered but not total pieces, unable to estimate how many items are being taken out of the given pack\nCheck your input",
	"Cannot Add more Stock, 3 max cost prices",
	"Need more info, check your input",
	"Warning\n This item is out of stock in the database\n You can continue, but be aware of this and the possible need to recount inventory",
}

func HandleErrorWindow(err error, w fyne.Window) bool {
	if err != nil {
		dialog.ShowError(err, w)
		return true
	}
	return false
}

func HandleKnownError(id int, condition bool, w fyne.Window) bool {
	if condition {
		dialog.ShowInformation("Error", knownErrors[id], w)
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
