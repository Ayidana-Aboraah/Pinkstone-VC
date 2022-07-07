package Database

import (
	"BronzeHermes/Cam"
	"BronzeHermes/UI"
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeInfoMenu(w fyne.Window) fyne.CanvasObject {
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	costLabel := widget.NewLabel("Cost")

	inventoryData := binding.BindIntList(&[]int{})
	inventoryData.Set(ConvertItemKeys())

	expenseData := binding.BindUntypedList(&[]interface{}{})
	expenseData.Set(ConvertExpenses())

	inventoryList := widget.NewListWithData(inventoryData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {
		val, err := item.(binding.Int).Get() // NOTE: Err handling
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

		ctq := "Cost : Quantity\n"
		for _, v := range ItemKeys[uint64(id)].Idxes {
			ctq += fmt.Sprintf("%f : %d\n", Items[v].Cost, Items[v].Quantity)
		}

		idLabel.SetText(strconv.Itoa(id))
		nameLabel.SetText(ItemKeys[uint64(idx)].Name)
		priceLabel.SetText(fmt.Sprint(ItemKeys[uint64(idx)].Price))
		costLabel.SetText(ctq)
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
			freq_text += "Yearly"
		case MONTHLY:
			freq_text += "Monthly"
		case ONCE:
			freq_text += fmt.Sprintf("Year/Month/Day: 20%v/%v/%v", Expenses[id].Date[2], Expenses[id].Date[1], Expenses[id].Date[0])
		}

		priceLabel.SetText(freq_text)
		costLabel.SetText("")
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
			costLabel,
			widget.NewButton("New Item", func() { // REVIEW: Change name
				id := Cam.OpenCam(&w)
				if id == 0 {
					return
				}

				val, found := ItemKeys[uint64(id)] // NOTE: if !found {Open up modify menu or something of the sort}
				if !found {
					ItemKeys[uint64(id)] = &ItemEV{Idxes: []int{len(Items)}}
					Items = append(Items, Item{})
				}

				ctq := "Cost : Quantity\n"
				for _, v := range ItemKeys[uint64(id)].Idxes {
					ctq += fmt.Sprintf("%2f : %d\n", Items[v].Cost, Items[v].Quantity)
				}

				idLabel.SetText(strconv.Itoa(id))
				nameLabel.SetText(val.Name)
				priceLabel.SetText(fmt.Sprint(val.Price))
				costLabel.SetText(ctq)
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

					Expenses = append(Expenses, Expense{
						Name:      items[0].Widget.(*widget.Entry).Text,
						Amount:    float32(amount),
						Frequency: expense_frequency,
						Date:      [3]uint8{uint8(day), uint8(month), uint8(year)},
					})
					expenseData.Set(ConvertExpenses())
					UI.HandleErrorWindow(SaveData(), w)
				}, w)
			}),
			widget.NewButton("Remove", func() {
				dialog.ShowConfirm("Are you sure?", "Are you sure you want to delete this?", func(b bool) {
					if !b {
						return
					}

					if arr == 0 { // Inventory
						itemData := binding.NewIntList()
						itemData.Set(ConvertItemIdxes(uint64(target)))

						itemList := widget.NewListWithData(itemData,
							func() fyne.CanvasObject { return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N")) },
							func(item binding.DataItem, obj fyne.CanvasObject) {
								val, err := item.(binding.Int).Get()
								UI.HandleError(err)
								obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%f : %d", Items[val].Cost, Items[val].Quantity))
							})

						itemList.OnSelected = func(id widget.ListItemID) {
							RemoveItem(id, uint64(target))
						}

						title := widget.NewLabelWithStyle(ItemKeys[uint64(target)].Name+"Price : Cost\n", fyne.TextAlignCenter, fyne.TextStyle{})

						dialog.ShowCustom("Which one?", "Done", container.NewVBox(title, itemList), w)
						inventoryList.Refresh()
					} else { // Expenses
						if len(Expenses)-1 == 0 {
							dialog.ShowInformation("Nope", "You cannot remove all your items like this, try adding a new item, and removing the last old one.", w)
							return
						}
						RemoveExpense(target)
						expenseList.Refresh()
					}
					UI.HandleErrorWindow(SaveData(), w)
				}, w)
			}),

			widget.NewButton("Modify", func() {
				itemData := binding.NewIntList()
				itemData.Set(ConvertItemIdxes(uint64(target)))

				itemList := widget.NewListWithData(itemData,
					func() fyne.CanvasObject { return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N")) },
					func(item binding.DataItem, obj fyne.CanvasObject) {
						val, err := item.(binding.Int).Get()
						UI.HandleError(err)
						obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%f : %d", Items[val].Cost, Items[val].Quantity))
					})

				itemList.OnSelected = func(id widget.ListItemID) {
					nameEntry := widget.NewEntry()
					nameEntry.SetPlaceHolder("Product Name with _ for spaces.")
					nameEntry.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")

					items := []*widget.FormItem{
						widget.NewFormItem("ID", idLabel),
						widget.NewFormItem("Name", nameEntry),
						widget.NewFormItem("Price", UI.NewNumEntry("Selling Price")),
						widget.NewFormItem("Cost", UI.NewNumEntry("How much you bought it for")),
						widget.NewFormItem("Inventory", UI.NewNumEntry("Current Inventory")),
					}

					dialog.ShowForm("Item", "Save", "Cancel",
						items,
						func(b bool) {
							if !b {
								return
							}

							priceText := items[2].Widget.(*UI.NumEntry).Text
							costText := items[3].Widget.(*UI.NumEntry).Text
							inventoryText := items[4].Widget.(*UI.NumEntry).Text

							val, err := itemData.GetValue(id)
							UI.HandleErrorWindow(err, w)

							price, cost, inventory := ConvertString(priceText, costText, inventoryText)
							y, month, day := time.Now().Date()
							year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

							newItem := Sale{Year: uint8(year), Month: uint8(month), Day: uint8(day), ID: uint64(target), Price: price, Cost: cost, Quantity: inventory}

							Reports[1] = append(Reports[1], newItem)

							ItemKeys[uint64(target)] = &ItemEV{Name: nameEntry.Text, Price: price, Idxes: ItemKeys[uint64(target)].Idxes} // WATCH: Unsure of how this will go
							Items[val] = Item{Quantity: inventory, Cost: cost}

							inventoryData.Set(ConvertItemKeys()) // NOTE: CHANGE this to ItemDB

							UI.HandleErrorWindow(SaveData(), w)

							dialog.NewInformation("Success!", "Your data has been saved successfully!", w)
						}, w)
				}

				title := widget.NewLabelWithStyle(ItemKeys[uint64(target)].Name+"Price : Cost\n", fyne.TextAlignCenter, fyne.TextStyle{})

				costEntry := UI.NewNumEntry("How much did you buy it for?")
				invenEntry := UI.NewNumEntry("How much in stock with this cost?")

				createBtn := widget.NewButton("New Inventory", func() {
					dialog.ShowForm("New Entry", "Create", "Cancel", []*widget.FormItem{
						widget.NewFormItem("ID", widget.NewLabel(strconv.Itoa(target))),
						widget.NewFormItem("Name", widget.NewLabel(ItemKeys[uint64(target)].Name)),
						widget.NewFormItem("Cost", costEntry),
						widget.NewFormItem("Quantity", invenEntry),
					}, func(b bool) {
						if !b {
							return
						}

						_, cost, quan := ConvertString("0", costEntry.Text, invenEntry.Text)
						ItemKeys[uint64(target)].Idxes = append(ItemKeys[uint64(target)].Idxes, len(Items))
						Items = append(Items, Item{Quantity: quan, Cost: cost})

					}, w)
				})

				dialog.ShowCustom("Which one?", "Done", container.NewVBox(title, itemList, createBtn), w)
			}),
		),
		container.NewVSplit(
			container.NewMax(inventoryList),
			container.NewMax(expenseList),
		))
}
