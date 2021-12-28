package UI

import (
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type NumEntry struct {
	widget.Entry
}

func (n *NumEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func NewNumEntry(placeHolder string) *NumEntry {
	e := &NumEntry{}
	e.ExtendBaseWidget(e)
	e.Validator = validation.NewRegexp(`\d`, "Must contain a number")
	e.PlaceHolder = placeHolder
	return e
}

func (e *NumEntry) TypedRune(r rune) {
	switch r {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '/':
		e.Entry.TypedRune(r)
	}
}
