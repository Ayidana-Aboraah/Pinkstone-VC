package UI

import (
	"fmt"
	"math/big"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

var knownErrors = []string{
	"Invalid Input\nCheck your input",
	"Pieces answered but not total pieces, unable to estimate how many items are being taken out of the given pack\nCheck your input",
	"Cannot Add more Stock, 3 max cost prices",
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

func Truncate(f float64, units float64) float64 {
	bV := big.NewFloat(0).SetPrec(1000).SetFloat64(f)
	bU := big.NewFloat(0).SetPrec(1000).SetFloat64(units)
	bV.Quo(bV, bU)

	i := big.NewInt(0)
	bV.Int(i)
	bV.SetInt(i)

	f, _ = bV.Mul(bV, bU).Float64()
	return f
}
