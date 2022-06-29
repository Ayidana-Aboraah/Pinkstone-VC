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
	inventoryLabel := widget.NewLabel("Inventory")

	inventoryData := binding.BindUntypedList(&[]interface{}{})
	inventoryData.Set(ConvertCart(Databases[0]))

	expenseData := binding.BindUntypedList(&[]interface{}{})
	expenseData.Set(ConvertExpenses())

	inventoryList := widget.NewListWithData(inventoryData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	inventoryList.UpdateItem = func(idx widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(NameKeys[Databases[0][idx].ID])
	}

	expenseList := widget.NewListWithData(expenseData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	expenseList.UpdateItem = func(idx widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Expenses[idx].Name)
	}

	var target int
	var arr uint8

	inventoryList.OnSelected = func(id widget.ListItemID) {
		item := Databases[0][id]
		values := ConvertSale(item)

		idLabel.SetText(strconv.Itoa(int(item.ID)))
		nameLabel.SetText(NameKeys[item.ID])
		priceLabel.SetText(values[0])
		costLabel.SetText(values[1])
		inventoryLabel.SetText(values[2])
		inventoryList.Unselect(id)

		arr = 0
		target = id
	}

	expenseList.OnSelected = func(id widget.ListItemID) {
		idLabel.SetText("" + Expenses[id].Name)
		nameLabel.SetText(fmt.Sprintf("Amount: %v", Expenses[id].Amount))

		freq_text := "Frequency: "

		switch Expenses[id].Frequency {
		case YEARLY:
			freq_text += "Yearly"
		case MONTHLY:
			freq_text += "Monthly"
		case ONCE:
			freq_text += fmt.Sprintf("20%v/%v/%v", Expenses[id].Date[2], Expenses[id].Date[1], Expenses[id].Date[0])
		}

		priceLabel.SetText(freq_text)
		costLabel.SetText("")
		inventoryLabel.SetText("")
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
			inventoryLabel,
			widget.NewButton("New Item", func() {
				id := Cam.OpenCam(&w)
				if id == 0 {
					return
				}

				result := FindItem(id)
				labels := ConvertSale(result)

				idLabel.SetText(strconv.Itoa(id))
				nameLabel.SetText(NameKeys[result.ID])
				priceLabel.SetText(labels[0])
				costLabel.SetText(labels[1])
				inventoryLabel.SetText(labels[2])

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
					fmt.Println(Expenses) // DEBUG
					UI.HandleErrorWindow(SaveData(), w)
				}, w)
			}),
			widget.NewButton("Remove", func() {
				if arr == 0 { // Inventory
					if len(Databases[ITEMS])-1 == 0 {
						dialog.ShowInformation("Nope", "You cannot remove all your items like this, try adding a new item, and removing the last old one.", w)
						return
					}
					Databases[ITEMS][target] = Databases[ITEMS][len(Databases[ITEMS])-1]
					Databases[ITEMS] = Databases[ITEMS][:len(Databases[ITEMS])-1]
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
			}),

			widget.NewButton("Modify", func() {
				conID, _ := strconv.Atoi(idLabel.Text)

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

						price, cost, inventory := ConvertString(priceText, costText, inventoryText)
						newItem := Sale{ID: uint64(conID), Price: price, Cost: cost, Quantity: inventory}

						Databases[2] = append(Databases[2], newItem)
						NameKeys[uint64(conID)] = nameEntry.Text

						func(found bool) {
							for i, v := range Databases[0] {
								if v.ID == newItem.ID {
									Databases[0][i] = newItem
									found = true
									break
								}
							}

							if !found {
								Databases[0] = append(Databases[0], newItem)
							}
						}(false)

						inventoryData.Set(ConvertCart(Databases[0]))

						UI.HandleErrorWindow(SaveData(), w)

						dialog.NewInformation("Success!", "Your data has been saved successfully!", w)
					}, w)
			}),
		),
		container.NewVSplit(
			container.NewMax(inventoryList),
			container.NewMax(expenseList),
		))
}
