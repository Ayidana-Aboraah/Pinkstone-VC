package test

import (
	"BronzeHermes/Database"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func TestMenu(shoppingCart *[]Database.Sale, a fyne.App, w fyne.Window) fyne.CanvasObject {
	items := widget.NewListWithData(
		binding.BindUntypedList(&[]interface{}{Database.ConvertCart(Database.Reports[0])}),
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {},
	)

	items.UpdateItem = func(idx widget.ListItemID, obj fyne.CanvasObject) {
		// Change this to TestDB if needed
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Database.NameKeys[Database.Reports[0][idx].ID])
	}

	items.OnSelected = func(id widget.ListItemID) {
		*shoppingCart = append(*shoppingCart, Database.Reports[0][id])
		items.Unselect(id)
	}

	return container.NewVBox(
		widget.NewButton("Display Database", func() {
			dialog.ShowInformation("Databases", fmt.Sprint(Database.Reports), w)
			dialog.ShowInformation("Name Keys", fmt.Sprint(Database.NameKeys), w)
		}),
		widget.NewButton("Load Test DB", func() {
			Database.NameKeys = TestNames
			Database.Reports = TestDB
		}),
		widget.NewButton("Load Test Expenses", func() { Database.Expenses = TestExpenses }),
		widget.NewButton("Add Item to Shopping Cart", func() {
			dialog.ShowCustom("Test Items", "Done", items, w)
		}),
	)
}

var TestNames = map[uint64]string{
	999999999999: "sammy",
	674398202423: "Cark",
	389432143927: "Banker",
	402933466372: "Blackeg",
	198998421024: "Boomb",
	412341251434: "Sus",
}

var TestDB = [2][]Database.Sale{
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
