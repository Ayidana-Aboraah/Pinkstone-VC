package Database

import (
	"BronzeHermes/UI"
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeInfoMenu(w fyne.Window) fyne.CanvasObject {
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	stockLabel := widget.NewLabel("Stock")

	inventoryData := binding.NewIntList()
	inventoryData.Set(ConvertItemKeys())

	expenseData := binding.NewUntypedList()
	expenseData.Set(ConvertExpenses())

	inventoryList := widget.NewListWithData(inventoryData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {
		val, err := item.(binding.Int).Get()
		UI.HandleError(err)
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(ItemKeys[uint64(val)].Name)
	})

	expenseList := widget.NewListWithData(expenseData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	expenseList.UpdateItem = func(idx widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Expenses[idx].Name)
	}

	var target int
	var arr uint8

	inventoryList.OnSelected = func(idx widget.ListItemID) {
		id, _ := inventoryData.GetValue(idx) // NOTE: Handle Err

		idLabel.SetText("ID: " + strconv.Itoa(id))
		nameLabel.SetText("Name: " + ItemKeys[uint64(id)].Name)
		priceLabel.SetText(fmt.Sprintf("Price: %1.1f", ItemKeys[uint64(id)].Price))
		stockLabel.SetText(fmt.Sprintf("Stock: %d\n", ItemKeys[uint64(id)].Quantity))
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

					v, found := ItemKeys[uint64(id)]
					if !found {
						dialog.ShowInformation("Lol", "Item Not Found, may not be registered", w)
						return
					}

					idLabel.SetText(selectionEntry.Text)
					nameLabel.SetText(ItemKeys[uint64(id)].Name)
					priceLabel.SetText(fmt.Sprintf("%1.1f", v.Price))
					stockLabel.SetText(fmt.Sprintf("%d", v.Quantity))
				}, w)
			}),

			widget.NewButton("New Item", func() {
				idEntry := UI.NewNumEntry("Click and Scan")
				nameEntry := widget.NewEntry()
				priceEntry := UI.NewNumEntry("Price to Sell it at")
				invenEntry := UI.NewNumEntry("Current amount in Stock")

				dialog.ShowForm("New Entry", "Create", "Cancel", []*widget.FormItem{
					widget.NewFormItem("ID", idEntry),
					widget.NewFormItem("Name", nameEntry),
					widget.NewFormItem("Price", priceEntry),
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

					price, quan := ConvertString(priceEntry.Text, invenEntry.Text)

					ItemKeys[uint64(id)] = &ItemEV{Price: price, Name: nameEntry.Text, Quantity: quan}
					delete(ItemKeys, 0)

					UI.HandleErrorWindow(SaveData(), w)

					inventoryData.Set(ConvertItemKeys()) //TODO: figure out how to ensure the list is updated visually
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
					if err != nil {
						log.Println(err)
					}

					if Expenses[0].Name == "" {
						Expenses[0] = Expense{
							Name:      items[0].Widget.(*widget.Entry).Text,
							Amount:    float32(amount),
							Frequency: expense_frequency,
							Date:      [3]uint8{uint8(day), uint8(month), uint8(year)},
						}
					} else {
						Expenses = append(Expenses, Expense{
							Name:      items[0].Widget.(*widget.Entry).Text,
							Amount:    float32(amount),
							Frequency: expense_frequency,
							Date:      [3]uint8{uint8(day), uint8(month), uint8(year)},
						})
					}

					expenseData.Set(ConvertExpenses())
					UI.HandleErrorWindow(SaveData(), w)
				}, w)
			}),

			widget.NewButton("Modify", func() {
				if target == 0 {
					return
				}
				priceEntry := UI.NewNumEntry("Price to Sell it at")
				invenEntry := UI.NewNumEntry("Current amount in stock")

				dialog.ShowForm("New Entry", "Create", "Cancel", []*widget.FormItem{
					widget.NewFormItem("ID", widget.NewLabel(strconv.Itoa(target))),
					widget.NewFormItem("Name", widget.NewLabel(ItemKeys[uint64(target)].Name)),
					widget.NewFormItem("Price", priceEntry),
					widget.NewFormItem("Current Stock", invenEntry),
				}, func(b bool) {
					if !b {
						return
					}

					price, quan := ConvertString(priceEntry.Text, invenEntry.Text)

					ItemKeys[uint64(target)] = &ItemEV{Price: price, Name: ItemKeys[uint64(target)].Name, Quantity: quan}

					UI.HandleErrorWindow(SaveData(), w)

					inventoryData.Set(ConvertItemKeys()) //TODO: figure out how to ensure the list is updated visually
					inventoryList.Select(inventoryList.Length() - 1)
					inventoryList.Refresh()
				}, w)
			}),

			widget.NewButton("Remove", func() { //TODO: make adjustments to functionality
				dialog.ShowConfirm("Are you sure?", "Are you sure you want to delete this?", func(b bool) {
					if !b {
						return
					}

					if arr == 0 { // Inventory
						if target == 0 {
							return
						}

						delete(ItemKeys, uint64(target))
						if len(ItemKeys) == 0 {
							ItemKeys[0] = &ItemEV{}
						}

						inventoryList.Refresh()
					} else { // Expenses

						RemoveExpense(target)
						if len(Expenses) == 0 {
							Expenses = append(Expenses, Expense{})
						}

						expenseList.Refresh()
					}
					UI.HandleErrorWindow(SaveData(), w)
				}, w)
			}),
		),

		container.NewVSplit(
			container.NewMax(inventoryList),
			container.NewMax(expenseList),
		))
}
