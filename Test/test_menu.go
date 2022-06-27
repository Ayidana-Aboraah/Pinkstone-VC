package test

import (
	"BronzeHermes/Database"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func TestMenu(a fyne.App, w fyne.Window) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewButton("Display Database", func() {
			dialog.ShowInformation("Databases", fmt.Sprint(Database.Databases), w)
			dialog.ShowInformation("Name Keys", fmt.Sprint(Database.NameKeys), w)
		}))
}

var TestNames = map[uint64]string{
	999999999999: "sammy",
	674398202423: "Cark",
	389432143927: "Banker",
	402933466372: "Blackeg",
	198998421024: "Boomb",
	412341251434: "Sus",
}

var TestDB = [3][]Database.Sale{
	{
		{ID: 999999999999, Price: 234.23, Cost: 1324, Quantity: 1},
		{ID: 674398202423, Price: 100.50, Cost: 1324, Quantity: 1},
		{ID: 389432143927, Price: 3974.89, Cost: 8934.24, Quantity: 5},
		{ID: 674398202423, Price: 90109.22, Cost: 48.24, Quantity: 87},
		{ID: 402933466372, Price: 1324.89, Cost: 21432.24, Quantity: 4124},
		{ID: 198998421024, Price: 1094.89, Cost: 9021038.24, Quantity: 5},
		{ID: 412341251434, Price: 3974.89, Cost: 8934.24, Quantity: 41},
	},
	{
		{Year: 22, Month: 10, Day: 1, ID: 674398202423, Price: 111.23, Cost: 1324, Quantity: 1},
		{Year: 22, Month: 9, Day: 4, ID: 674398202423, Price: 100.50, Cost: 555, Quantity: 1},
		{Year: 22, Month: 8, Day: 5, ID: 389432143927, Price: 222.89, Cost: 332.24, Quantity: 5},
		{Year: 22, Month: 7, Day: 6, ID: 674398202423, Price: 444.22, Cost: 222.24, Quantity: 7},
		{Year: 22, Month: 6, Day: 4, ID: 402933466372, Price: 333.21, Cost: 232.24, Quantity: 4},
		{Year: 22, Month: 6, Day: 4, ID: 198998421024, Price: 555.22, Cost: 938.24, Quantity: 5},
		{Year: 22, Month: 6, Day: 7, ID: 412341251434, Price: 666.22, Cost: 834.24, Quantity: 1},
	},
	{},
}

var TestExpenses = []Database.Expense{
	{Date: [3]uint8{7, 6, 22}, Frequency: 1, Amount: 1, Name: "red"},
	{Date: [3]uint8{7, 6, 22}, Frequency: 1, Amount: -43},

	{Date: [3]uint8{8, 6, 22}, Frequency: 0, Amount: -43},
	{Date: [3]uint8{7, 7, 22}, Frequency: 1, Amount: 3},
	{Date: [3]uint8{7, 6, 23}, Frequency: 1, Amount: -13},
}
