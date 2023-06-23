package Database

import (
	"BronzeHermes/UI"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var InventoryData binding.IntList
var ExpenseData binding.UntypedList

func MakeInfoMenu(w fyne.Window) fyne.CanvasObject {
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	stockLabel := widget.NewLabel("Stock")

	InventoryData = binding.NewIntList()
	InventoryData.Set(ConvertItemKeys())

	ExpenseData = binding.NewUntypedList()
	ExpenseData.Set(ConvertExpenses())

	inventoryList := widget.NewListWithData(InventoryData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {
		val, err := item.(binding.Int).Get()
		UI.HandleError(err)
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Item[uint64(val)].Name)
	})

	expenseList := widget.NewListWithData(ExpenseData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	expenseList.UpdateItem = func(idx widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Expenses[idx].Name)
	}

	target := -1
	var arr uint8

	inventoryList.OnSelected = func(idx widget.ListItemID) {
		id, _ := InventoryData.GetValue(idx) // NOTE: Handle Err

		idLabel.SetText("ID: " + strconv.Itoa(id))
		nameLabel.SetText("Name: " + Item[uint64(id)].Name)
		priceLabel.SetText(fmt.Sprintf("Price: %1.1f", Item[uint64(id)].Price))
		stockLabel.SetText(fmt.Sprintf("Stock: %d\n", Item[uint64(id)].Quantity[0]+Item[uint64(id)].Quantity[1]+Item[uint64(id)].Quantity[2]))
		inventoryList.Unselect(idx)

		arr = 0
		target = id
	}

	expenseList.OnSelected = func(id widget.ListItemID) {
		idLabel.SetText("ID: " + Expenses[id].Name)
		nameLabel.SetText(fmt.Sprintf("Amount: %v", Expenses[id].Amount))

		freq_text := "Frequency: "

		switch Expenses[id].Frequency {
		case YEARLY:
			freq_text = "Frequency Yearly"
		case MONTHLY:
			freq_text = "Frequency Monthly"
		case ONCE:
			freq_text = fmt.Sprintf("Once: 20%v/%v/%v", Expenses[id].Date[2], Expenses[id].Date[1], Expenses[id].Date[0])
		}

		priceLabel.SetText(freq_text)
		stockLabel.SetText("")
		expenseList.Unselect(id)

		arr = 1
		target = id
	}

	return container.New(layout.NewGridLayout(2),
		container.NewVBox(
			widget.NewLabelWithStyle("Inventory Info", fyne.TextAlign(1), fyne.TextStyle{Bold: true}),
			idLabel,
			nameLabel,
			priceLabel,
			stockLabel,

			widget.NewButton("Search", func() {
				selectionEntry := UI.NewNumEntry("Click and Scan")

				dialog.ShowCustomConfirm("Scan Item", "Confirm", "Cancel", selectionEntry, func(confirmed bool) {
					if !confirmed {
						return
					}

					id, err := strconv.Atoi(selectionEntry.Text)
					if err != nil {
						dialog.ShowInformation("Hahaha", "Invalid Barcode", w)
						return
					}

					v, found := Item[uint64(id)]
					if !found {
						dialog.ShowInformation("Lol", "Item Not Found, may not be registered", w)
						return
					}

					idLabel.SetText(selectionEntry.Text)
					nameLabel.SetText(Item[uint64(id)].Name)
					priceLabel.SetText(fmt.Sprintf("%1.1f", v.Price))
					stockLabel.SetText(fmt.Sprintf("%d", v.Quantity[0]+v.Quantity[1]+v.Quantity[2]))
				}, w)
			}),

			widget.NewButton("New Item", func() {
				idEntry := UI.NewNumEntry("Click and Scan")
				nameEntry := widget.NewEntry()
				priceEntry := UI.NewNumEntry("Selling Price")
				costEntry := UI.NewNumEntry("Cost Price")
				invenEntry := UI.NewNumEntry("Current amount in Stock")

				dialog.ShowForm("New Entry", "Create", "Cancel", []*widget.FormItem{
					widget.NewFormItem("ID", idEntry),
					widget.NewFormItem("Name", nameEntry),
					widget.NewFormItem("Price", priceEntry),
					widget.NewFormItem("Cost", costEntry),
					widget.NewFormItem("Current Stock", invenEntry),
				}, func(b bool) {
					if !b {
						return
					}

					id, err := strconv.ParseUint(idEntry.Text, 10, 64)
					if err != nil {
						dialog.ShowInformation("Nope...", "Invalid Barcode", w)
						return
					}

					price, cost, quan := ConvertString(priceEntry.Text, costEntry.Text, invenEntry.Text)

					Item[uint64(id)] = &Entry{Price: price, Name: nameEntry.Text, Quantity: [3]uint16{quan}, Cost: [3]float32{cost}}

					UI.HandleErrorWindow(SaveData(), w)

					InventoryData.Set(ConvertItemKeys())
					inventoryList.Select(inventoryList.Length() - 1)
					inventoryList.Refresh()
				}, w)
			}),

			widget.NewButton("New Expense/Gift", func() {
				var expense_frequency uint8

				items := []*widget.FormItem{
					widget.NewFormItem("Name ", widget.NewEntry()),
					widget.NewFormItem("Amount ", UI.NewNumEntry("The amount lost or gained")),
					widget.NewFormItem("Frequency ", widget.NewSelect([]string{"Once", "Monthly", "Yearly"}, func(s string) {
						switch s {
						case "Once":
							expense_frequency = ONCE
						case "Monthly":
							expense_frequency = MONTHLY
						case "Yearly":
							expense_frequency = YEARLY
						}
					})),
				}

				dialog.ShowForm("Expense", "Create", "Cancel", items, func(b bool) {
					if !b {
						return
					}

					y, month, day := time.Now().Date()
					year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

					amount, err := strconv.ParseFloat(items[1].Widget.(*UI.NumEntry).Text, 32)
					UI.HandleError(err)

					Expenses = append(Expenses, Expense{
						Name:      items[0].Widget.(*widget.Entry).Text,
						Amount:    float32(amount),
						Frequency: expense_frequency,
						Date:      [3]uint8{uint8(day), uint8(month), uint8(year)},
					})

					ExpenseData.Set(ConvertExpenses())
					UI.HandleErrorWindow(SaveData(), w)
				}, w)
			}),

			widget.NewButton("Modify", func() {
				if target == -1 {
					return
				}

				if arr == 0 {
					priceEntry := UI.NewNumEntry("Selling Price")
					costEntry := UI.NewNumEntry("Cost Price")
					invenEntry := UI.NewNumEntry("Added stock")

					dialog.ShowForm("New Entry", "Create", "Cancel", []*widget.FormItem{
						widget.NewFormItem("ID", widget.NewLabel(strconv.Itoa(target))),
						widget.NewFormItem("Name", widget.NewLabel(Item[uint64(target)].Name)),
						widget.NewFormItem("Price", priceEntry),
						widget.NewFormItem("cost", costEntry),

						widget.NewFormItem("Current Stock", invenEntry),
					}, func(b bool) {
						if !b {
							return
						}

						price, cost, quan := ConvertString(priceEntry.Text, costEntry.Text, invenEntry.Text)

						Item[uint64(target)] = &Entry{Price: price, Name: Item[uint64(target)].Name, Quantity: [3]uint16{quan}}

						for i := 0; i < 3; i++ {
							if Item[uint64(target)].Quantity[i] == 0 {
								Item[uint64(target)].Quantity[i] = quan
								Item[uint64(target)].Cost[i] = cost
								break
							}
						}

						UI.HandleErrorWindow(SaveData(), w)

						InventoryData.Set(ConvertItemKeys())
						inventoryList.Select(inventoryList.Length() - 1)
						inventoryList.Refresh()
					}, w)
				} else {
					var frequency uint8
					amountEntry := UI.NewNumEntry("Amount")
					dialog.ShowForm("New Entry", "Create", "Cancel", []*widget.FormItem{
						widget.NewFormItem("Expense", widget.NewLabel(Expenses[target].Name)),
						widget.NewFormItem("Amount", amountEntry),
						widget.NewFormItem("Frequency ", widget.NewSelect([]string{"Once", "Monthly", "Yearly"}, func(s string) {
							switch s {
							case "Once":
								frequency = ONCE
							case "Monthly":
								frequency = MONTHLY
							case "Yearly":
								frequency = YEARLY
							}
						})),
					}, func(b bool) {
						if !b {
							return
						}
						amount, err := strconv.ParseFloat(amountEntry.Text, 32)
						UI.HandleError(err)
						Expenses[target].Frequency = frequency
						Expenses[target].Amount = float32(amount)

					}, w)
				}
			}),

			widget.NewButton("Remove", func() {
				dialog.ShowConfirm("Are you sure?", "Are you sure you want to delete this?", func(b bool) {
					if !b {
						return
					}

					if target == -1 {
						return
					}

					if arr == 0 { // Inventory

						delete(Item, uint64(target))
						InventoryData.Set(ConvertItemKeys())
						inventoryList.Refresh()
					} else { // Expenses

						RemoveExpense(target)
						ExpenseData.Set(ConvertExpenses())
						expenseList.Refresh()
					}
					target = -1
					UI.HandleErrorWindow(SaveData(), w)
				}, w)
			}),
		),

		container.NewVSplit(
			container.NewMax(inventoryList),
			container.NewMax(expenseList),
		))
}
