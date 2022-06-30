package main

import (
	"BronzeHermes/Cam"
	"BronzeHermes/Database"
	"BronzeHermes/Graph"
	test "BronzeHermes/Test"
	"BronzeHermes/UI"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("Bronze Hermes")
	// go Graph.StartServer()

	Database.DataInit(false)

	CreateWindow(a)
}

func CreateWindow(a fyne.App) {
	w := a.NewWindow("Bronze Hermes")
	w.SetOnClosed(Graph.StopSever)

	if UI.HandleErrorWindow(Database.LoadData(), w) {
		dialog.ShowInformation("Back Up", "Loading BackUp", w)
		UI.HandleErrorWindow(Database.LoadBackUp(), w)
	}

	fmt.Println(Database.Expenses)

	w.SetContent(container.NewVBox(container.NewAppTabs(
		container.NewTabItem("Main", makeMainMenu(a, w)),
		container.NewTabItem("Shop", makeShoppingMenu(w)),
		container.NewTabItem("Inventory", Database.MakeInfoMenu(w)),
		container.NewTabItem("Statistics", makeStatsMenu(w)),
		container.NewTabItem("Debug", test.TestMenu(&shoppingCart, a, w)),
	)))

	w.ShowAndRun()
}

func makeMainMenu(a fyne.App, w fyne.Window) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabelWithStyle("Welcome", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		// NOTE: Not functional
		// widget.NewButton("Save Backup Data", func() {
		// 	go UI.HandleErrorWindow(Database.BackUpAllData(), w)
		// }),
		// widget.NewButton("Save Backup Data", func() {
		// 	dialog.ShowInformation("Loading Back up Data", "Wait until back up is done loading...", w)
		// 	go UI.HandleErrorWindow(Database.LoadBackUp(), w)
		// 	dialog.ShowInformation("Loaded", "Back Up Loaded", w)
		// }),

		//Add inventory features here
	)
}

var shoppingCart []Database.Sale

func makeShoppingMenu(w fyne.Window) fyne.CanvasObject {

	title := widget.NewLabelWithStyle("Cart Total: 0.0", fyne.TextAlignCenter, fyne.TextStyle{})

	cartList := binding.BindUntypedList(&[]interface{}{})

	shoppingList := widget.NewListWithData(cartList,
		func() fyne.CanvasObject {
			return container.NewBorder(nil, nil, nil, widget.NewButton("X", nil), widget.NewLabel(""))
		}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	shoppingList.OnSelected = func(id widget.ListItemID) {
		shoppingCart[id].Quantity++
		cartList.Reload()
		title.SetText(fmt.Sprintf("Cart Total: %.2f", Database.GetCartTotal(shoppingCart)))
		shoppingList.Unselect(id)
	}

	shoppingList.UpdateItem = func(id widget.ListItemID, obj fyne.CanvasObject) {
		text := obj.(*fyne.Container).Objects[0].(*widget.Label)
		btn := obj.(*fyne.Container).Objects[1].(*widget.Button)
		val := shoppingCart[id]

		text.SetText(Database.NameKeys[val.ID] + " x" + strconv.Itoa(int(val.Quantity)))
		btn.OnTapped = func() {
			cartList.Set(Database.ConvertCart(Database.DecreaseFromCart(val, shoppingCart)))
			title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
			text.SetText(Database.NameKeys[val.ID] + " x" + strconv.Itoa(int(val.Quantity)))
			shoppingList.Refresh()
		}
	}

	screen := container.New(layout.NewGridLayoutWithRows(3),
		title,
		container.NewMax(shoppingList),
		container.NewGridWithColumns(3,
			widget.NewButton("Buy Cart", func() {
				dialog.ShowConfirm("Buying", "Do you want to buy all items in the Cart?", func(b bool) {
					if !b {
						return
					}
					cartList.Set(Database.ConvertCart(Database.BuyCart(shoppingCart)))
					title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
					dialog.ShowInformation("Complete", "You're Purchase has been made.", w)
				}, w)
			}),
			widget.NewButton("Clear Cart", func() {
				cartList.Set([]interface{}{})
				shoppingCart = shoppingCart[:0]
				title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
			}),
			widget.NewButton("New Item", func() {
				id := Cam.OpenCam(&w)
				if id == 0 {
					return
				}

				item := Database.FindItem(id)

				dialog.ShowCustomConfirm("Just Checking...", "Yes", "No", container.NewVBox(widget.NewLabel("Is this the right item: "+Database.NameKeys[item.ID])),
					func(b bool) {
						if !b {
							return
						}
						cartList.Set(Database.ConvertCart(Database.AddToCart(item, shoppingCart)))
						title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
						shoppingList.Refresh()
					}, w)
			}),
			widget.NewButton("Reload", func() { // DEBUG
				cartList.Set(Database.ConvertCart(shoppingCart))
				title.SetText(fmt.Sprintf("Cart Total: %1.1f", Database.GetCartTotal(shoppingCart)))
				shoppingList.Refresh()
			}),
		),
	)
	return screen
}

func makeStatsMenu(w fyne.Window) fyne.CanvasObject {
	u, _ := url.Parse("http://localhost:8081/line")
	r, _ := url.Parse("http://localhost:8081/pie")

	link := widget.NewHyperlink("Go To Graph", u)

	selectionEntry := UI.NewNumEntry("Year/Month")

	var profitDataSelect int
	var buttonType int

	dataSelectOptions := widget.NewSelect([]string{"Revenue", "Cost", "Profit"}, func(dataType string) {
		switch dataType {
		case "Revenue":
			profitDataSelect = 0
		case "Cost":
			profitDataSelect = 1
		case "Profit":
			profitDataSelect = 2
		}
	})

	reportDisplay := widget.NewLabel("")
	financeEntry := UI.NewNumEntry("YYYY/MM/DD [Select Custom again to select the date]")

	return container.NewVScroll(container.NewMax(container.NewVBox(
		widget.NewCard("Financial Reports", "", container.NewVBox(
			financeEntry,
			widget.NewSelect([]string{"Day", "Month", "Year", "Date"}, func(time string) {
				var variant uint8

				financeEntry.Hidden = true

				date := []uint8{}

				switch time {
				case "Day":
					variant = Database.ONCE
				case "Month":
					variant = Database.MONTHLY
				case "Year":
					variant = Database.YEARLY
				case "Date": //The user will have to double tap when using Dates
					if financeEntry.Text == "" {
						return
					}

					//String to Date conversion
					raw := strings.Split(financeEntry.Text, "/")

					year, err := strconv.Atoi(raw[0][1:])
					if err != nil {
						fmt.Println("Something Seems up!") //DEBUG: REMOVE AFTER TESTING
						return
					}

					month, err := strconv.Atoi(raw[1])
					if err != nil {
						variant = 1
					}

					day, err := strconv.Atoi(raw[1])
					if err != nil {
						variant = 2
					}
					financeEntry.Hidden = false

					date = []uint8{uint8(day), uint8(month), uint8(year)}
				}

				reportDisplay.Text = Database.Report(variant, date)

			}),
			reportDisplay,
		)),
		widget.NewCard("Data Graphs", "", container.NewVBox(
			selectionEntry,
			widget.NewSelect([]string{"Items Graph", "Price Changes", "Item Popularity", "Item Sales"}, func(graph string) {
				switch graph {
				case "Items Graph":
					buttonType = 0
					link.URL = u
					dataSelectOptions.Hidden = false
				case "Price Changes":
					buttonType = 1
					link.URL = u
					dataSelectOptions.Hidden = false
				case "Item Popularity":
					buttonType = 2
					link.URL = r
					dataSelectOptions.Hidden = false
				case "Item Sales":
					buttonType = 3
					link.URL = r
					dataSelectOptions.Hidden = true
				case "Sales Over Time":
					buttonType = 4
					link.URL = u
					dataSelectOptions.Hidden = true
				}
			}),
			dataSelectOptions,
			widget.NewButton("Graph", func() {
				switch buttonType {
				case 0:
					Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, profitDataSelect, Database.Databases[1])
				case 1:
					Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, profitDataSelect, Database.Databases[2])
				case 4:
					Graph.Labels, Graph.LineInputs = Database.GetLine(selectionEntry.Text, 3, Database.Databases[1])
				case 2:
					Graph.Labels, Graph.Inputs = Database.GetPie(selectionEntry.Text, profitDataSelect)
				case 3:
					Graph.Labels, Graph.Inputs = Database.GetPie(selectionEntry.Text, 3)
				}
			}),
			link,
		)),
	)))
}
