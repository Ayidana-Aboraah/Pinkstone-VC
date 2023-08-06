package UI

import (
	"strings"

	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type SearchBar struct {
	search func(input string) ([]string, []uint16)
	names  []string
	idxs   []uint16
	widget.SelectEntry
}

func (n *SearchBar) Keyboard() mobile.KeyboardType {
	return mobile.DefaultKeyboard
}

func NewSearchBar(placeHolder string, search func(string) ([]string, []uint16)) *SearchBar {
	e := &SearchBar{}
	e.ExtendBaseWidget(e)
	e.PlaceHolder = placeHolder
	e.search = search
	e.names, e.idxs = search("")
	e.SetOptions(e.names)
	return e
}

func (e *SearchBar) TypedRune(r rune) {
	e.Entry.TypedRune(r)
	e.names, e.idxs = e.search(e.Text)
	e.SetOptions(e.names)
}

func (e *SearchBar) Result() int {
	for i, idx := range e.idxs {
		if strings.EqualFold(e.names[i], e.Text) {
			return int(idx)
		}
	}
	return -1
}

// Warning, is only usable on arrays especially if the array doesn't have remved items
func (e *SearchBar) ResultOrCreate(create func()) int {
	x := e.Result()
	if x == -1 {
		x = len(e.names)
		create()
	}
	return x
}
