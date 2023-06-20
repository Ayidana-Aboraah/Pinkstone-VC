package Test

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
	idData := binding.NewIntList()
	idData.Set(Database.ConvertItemKeys())

	items := widget.NewListWithData(
		idData,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			val, _ := item.(binding.Int).Get() // NOTE: Err handling
			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Database.ItemKeys[uint64(val)].Name)
		},
	)

	items.OnSelected = func(id widget.ListItemID) {
		*shoppingCart = append(*shoppingCart, Database.Reports[0][id])
		items.Unselect(id)
	}

	return container.NewVBox(
		widget.NewButton("Display Database", func() {
			dialog.ShowInformation("Databases", fmt.Sprint(Database.Reports), w)
		}),
		widget.NewButton("Load Test DB", func() {
			Database.ItemKeys = TestItemKeys
			Database.Items = TestItems
			Database.Reports = TestDB
			fmt.Println(TestDB[0])
		}),
		widget.NewButton("Load Test Expenses", func() { Database.Expenses = TestExpenses }),
		widget.NewButton("Add Item to Shopping Cart", func() {
			dialog.ShowCustom("Test Items", "Done", items, w)
		}),
	)
}

var TestItemKeys = map[uint64]*Database.ItemEV{
	999999999999: {Price: 234.23, Name: "sammy", Idxes: []int{0}},
	674398202423: {Price: 100.50, Name: "Clark", Idxes: []int{1, 3}},
	389432143927: {Price: 3974.89, Name: "Banker", Idxes: []int{2}},
	402933466372: {Price: 1324.89, Name: "Blackest", Idxes: []int{4}},
	198998421024: {Price: 1094.89, Name: "Reeses puffs", Idxes: []int{5}},
	412341251434: {Price: 3974.89, Name: "Sus", Idxes: []int{6, 7}},
}

var TestItems = []Database.Item{
	{Cost: 1324, Quantity: 1},
	{Cost: 1324, Quantity: 1},
	{Cost: 8934.24, Quantity: 5},
	{Cost: 48.24, Quantity: 87},
	{Cost: 21432.24, Quantity: 4124},
	{Cost: 9021038.24, Quantity: 5},
	{Cost: 8934.24, Quantity: 41},
	{Cost: 10.0, Quantity: 10},
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
