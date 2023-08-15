package Debug

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var knownErrors = []string{
	"Invalid Input\nCheck your input",
	"Cannot Add more Stock, 3 max cost prices",
	"Need more info, check your input",
	"Warning\n This item is out of stock in the database\n You can continue, but be aware of this and the possible need to recount inventory",
}

const Success = -1

const (
	Invalid_Input = iota
	Maxed_Out_Stocks
	Need_More_Info
	Empty_Quantity_Warning
)

func HandleKnownError(id int, condition bool, w fyne.Window) bool {
	if condition {
		dialog.ShowInformation("Error", knownErrors[id], w)
		return true
	}
	return false
}

func ShowError(state string, err error, w fyne.Window) bool {
	if err != nil {
		dialog.ShowInformation("Oops", "Error while"+state+"\nError: "+err.Error(), w)
		return true
	}
	return false
}
