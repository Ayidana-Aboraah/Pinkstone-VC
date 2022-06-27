package Database

import (
	"BronzeHermes/Cam"
	"BronzeHermes/UI"
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

	var expense_frequency uint8
	expense_entry := widget.NewEntry()
	expense_amount := UI.NewNumEntry("The amount gained or lost.")

	inventoryData := binding.BindUntypedList(&[]interface{}{})
	inventoryData.Set(ConvertCart(Databases[0]))

	// Expenses = test.TestExpenses

	expenseData := binding.BindUntypedList(&[]interface{}{})
	expenseData.Set(ConvertExpenses())

	items := []*widget.FormItem{
		widget.NewFormItem("Name ", expense_entry),
		widget.NewFormItem("Amount ", expense_amount),
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

	inventoryList := widget.NewListWithData(inventoryData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("name"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	inventoryList.UpdateItem = func(idx widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(NameKeys[Databases[0][idx].ID])
	}

	inventoryList.OnSelected = func(id widget.ListItemID) {
		item := Databases[0][id]
		values := ConvertSale(item)

		idLabel.SetText(strconv.Itoa(int(item.ID)))
		nameLabel.SetText(NameKeys[item.ID])
		priceLabel.SetText(values[0])
		costLabel.SetText(values[1])
		inventoryLabel.SetText(values[2])
	}

	expenseList := widget.NewListWithData(expenseData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("name"), widget.NewButton("X", func() {}))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {})

	expenseList.UpdateItem = func(idx widget.ListItemID, obj fyne.CanvasObject) {
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(NameKeys[Databases[0][idx].ID])
		obj.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
			dialog.ShowConfirm("Think about this.", "Are you sure you want to delete this expense?", func(b bool) {
				if !b {
					return
				}
				RemoveExpense(idx)
				expenseData.Set(ConvertExpenses())
			}, w)
		}
	}

	expenseList.OnSelected = func(id widget.ListItemID) {
		dialog.ShowForm("Modify", "Done", "Cancel", items, func(b bool) {

		}, w)
		// Create a form that shows the information on the currently selected list item
		// Modify attributes in the form
		// Convert the modified attributes into an expense
		// Fill the data into that index of the old expense
	}

	return container.New(layout.NewGridLayout(2),
		container.NewVBox(
			widget.NewLabelWithStyle("Inventory Info", fyne.TextAlign(1), fyne.TextStyle{Bold: true}),
			idLabel,
			nameLabel,
			priceLabel,
			costLabel,
			inventoryLabel,
			widget.NewButton("New", func() {
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
			widget.NewButton("Modify", func() {
				conID, _ := strconv.Atoi(idLabel.Text)
				idLabel := widget.NewLabel(strconv.Itoa(conID))

				nameEntry := widget.NewEntry()
				nameEntry.SetPlaceHolder("Product Name with _ for spaces.")
				nameEntry.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")

				priceEntry := UI.NewNumEntry("Selling Price")
				costEntry := UI.NewNumEntry("How much you bought it for")
				inventoryEntry := UI.NewNumEntry("Current Inventory")

				dialog.ShowForm("Item", "Save", "Cancel",
					[]*widget.FormItem{
						widget.NewFormItem("ID", idLabel),
						widget.NewFormItem("Name", nameEntry),
						widget.NewFormItem("Price", priceEntry),
						widget.NewFormItem("Cost", costEntry),
						widget.NewFormItem("Inventory", inventoryEntry),
					},
					func(b bool) {
						if !b {
							return
						}

						price, cost, inventory := ConvertString(priceEntry.Text, costEntry.Text, inventoryEntry.Text)
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

						//Updating Entries
						nameLabel.Text = nameEntry.Text
						priceLabel.Text = priceEntry.Text
						costLabel.Text = costEntry.Text
						inventoryLabel.Text = inventoryEntry.Text

						dialog.NewInformation("Success!", "Your data has been saved successfully!", w)
					}, w)
			}),
			widget.NewButton("Expense/Gift", func() {
				dialog.ShowForm("Expense", "Create", "Cancel", items, func(b bool) {
					if !b {
						return
					}

					day, month, y := time.Now().Date()
					year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

					amount, err := strconv.ParseFloat(expense_amount.Text, 32)
					if err != nil {
						log.Println(err)
					}

					Expenses = append(Expenses, Expense{
						Name:      expense_entry.Text,
						Amount:    float32(amount),
						Frequency: expense_frequency,
						Date:      [3]uint8{uint8(day), uint8(month), uint8(year)},
					})
					expenseData.Set(ConvertExpenses())
				}, w)
			}),
		),
		container.NewVSplit(
			container.NewMax(inventoryList),
			container.NewMax(expenseList),
		))
}
