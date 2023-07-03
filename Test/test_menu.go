package Test

// import (
// 	"BronzeHermes/Database"
// 	"fmt"

// 	"fyne.io/fyne/v2"
// 	"fyne.io/fyne/v2/container"
// 	"fyne.io/fyne/v2/data/binding"
// 	"fyne.io/fyne/v2/dialog"
// 	"fyne.io/fyne/v2/widget"
// )

// func TestMenu(shoppingCart *[]Database.Sale, a fyne.App, w fyne.Window) fyne.CanvasObject {
// 	idData := binding.NewIntList()
// 	idData.Set(Database.ConvertItemKeys())

// 	items := widget.NewListWithData(
// 		idData,
// 		func() fyne.CanvasObject {
// 			return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
// 		},
// 		func(item binding.DataItem, obj fyne.CanvasObject) {
// 			val, _ := item.(binding.Int).Get() // NOTE: Err handling
// 			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Database.Item[uint64(val)].Name)
// 		},
// 	)

// 	items.OnSelected = func(id widget.ListItemID) {
// 		*shoppingCart = append(*shoppingCart, Database.Reports[0][id])
// 		items.Unselect(id)
// 	}

// 	return container.NewVBox(
// 		widget.NewButton("Display Database", func() {
// 			dialog.ShowInformation("Databases", fmt.Sprint(Database.Reports), w)
// 		}),
// 		widget.NewButton("Load Test DB", func() {
// 			// Database.ItemKeys = TestItemKeys
// 			Database.Reports = TestDB
// 			fmt.Println(TestDB[0])
// 		}),
// 		widget.NewButton("Load Test Expenses", func() { Database.Expenses = TestExpenses }),
// 		widget.NewButton("Add Item to Shopping Cart", func() {
// 			dialog.ShowCustom("Test Items", "Done", items, w)
// 		}),
// 	)
// }

// // var TestItemKeys = map[uint64]*Database.ItemEV{
// // 	999999999999: {Price: 234.23, Name: "sammy", Quantity: 1},
// // 	674398202423: {Price: 100.50, Name: "Clark", Quantity: 1},
// // 	389432143927: {Price: 3974.89, Name: "Banker", Quantity: 5},
// // 	402933466372: {Price: 1324.89, Name: "Blackest", Quantity: 87},
// // 	198998421024: {Price: 1094.89, Name: "Reeses puffs", Quantity: 4124},
// // 	412341251434: {Price: 3974.89, Name: "Sus", Quantity: 5},
// // }

// var TestDB = [2][]Database.Sale{
// 	{
// 		{Year: 22, Month: 10, Day: 1, ID: 674398202423, Price: 111.23, Quantity: 1},
// 		{Year: 22, Month: 9, Day: 4, ID: 674398202423, Price: 100.50, Quantity: 1},
// 		{Year: 22, Month: 8, Day: 5, ID: 389432143927, Price: 222.89, Quantity: 5},
// 		{Year: 22, Month: 7, Day: 6, ID: 674398202423, Price: 444.22, Quantity: 7},
// 		{Year: 22, Month: 6, Day: 4, ID: 402933466372, Price: 333.21, Quantity: 4},
// 		{Year: 22, Month: 6, Day: 4, ID: 198998421024, Price: 555.22, Quantity: 5},
// 		{Year: 22, Month: 6, Day: 7, ID: 412341251434, Price: 666.22, Quantity: 1},
// 	},
// 	{},
// }

// var TestExpenses = []Database.Expense{
// 	{Date: [3]uint8{7, 6, 22}, Frequency: 1, Amount: 1, Name: "red"},
// 	{Date: [3]uint8{7, 6, 22}, Frequency: 1, Amount: -43},

// 	{Date: [3]uint8{8, 6, 22}, Frequency: 0, Amount: -43},
// 	{Date: [3]uint8{7, 7, 22}, Frequency: 1, Amount: 3},
// 	{Date: [3]uint8{7, 6, 23}, Frequency: 1, Amount: -13},
// }
