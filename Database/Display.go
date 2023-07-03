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

type ListID interface {
	int | uint16
}

var UpdateInventoryDisplay func(id any)

func MakeInfoMenu(w fyne.Window) fyne.CanvasObject {
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	costLabel := widget.NewLabel("Cost")
	stockLabel := widget.NewLabel("Stock")

	InventoryData = binding.NewIntList()
	InventoryData.Set(ConvertItemKeys())

	inventoryList := widget.NewListWithData(InventoryData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {
		val, err := item.(binding.Int).Get()
		UI.HandleError(err)
		obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(Item[uint16(val)].Name)
	})

	target := -1

	UpdateInventoryDisplay = func(idx any) {
		var id int
		switch x := idx.(type) {
		case uint16:
			id = int(x)
		case int:
			id, _ = InventoryData.GetValue(x) // NOTE: Handle Err
			inventoryList.Unselect(x)
		}

		idLabel.SetText("ID: " + strconv.Itoa(id))
		nameLabel.SetText("Name: " + Item[uint16(id)].Name)
		priceLabel.SetText(fmt.Sprintf("Price: %1.2f", Item[uint16(id)].Price))
		stockLabel.SetText(fmt.Sprintf("Stock: %1.2f\n", Item[uint16(id)].Quantity[0]+Item[uint16(id)].Quantity[1]+Item[uint16(id)].Quantity[2]))

		txt := "Cost: "

		for i := 0; i < 3 && Item[uint16(id)].Cost[i] > 0; i++ {
			txt += fmt.Sprintf("%1.2f, ", Item[uint16(id)].Cost[i])
		}

		if Item[uint16(id)].Cost[0] == 0 {
			txt += "0.00 "
		}

		costLabel.SetText(txt[:len(txt)-1])
		target = id
	}

	inventoryList.OnSelected = func(id widget.ListItemID) {
		UpdateInventoryDisplay(id)
	}

	return container.New(layout.NewGridLayout(2),
		container.NewVBox(
			widget.NewLabelWithStyle("Inventory Info", fyne.TextAlign(1), fyne.TextStyle{Bold: true}),
			idLabel,
			nameLabel,
			priceLabel,
			costLabel,
			stockLabel,

			widget.NewButton("New Item", func() {
				nameEntry := widget.NewEntry()
				priceEntry := UI.NewNumEntry("Selling Price")
				costEntry := UI.NewNumEntry("Cost Price")
				invenEntry := UI.NewNumEntry("Current amount in Stock, always round down if it exceeds 2 decminal places")

				dialog.ShowForm("New Entry", "Create", "Cancel", []*widget.FormItem{
					widget.NewFormItem("Name", nameEntry),
					widget.NewFormItem("Price", priceEntry),
					widget.NewFormItem("Cost", costEntry),
					widget.NewFormItem("Current Stock", invenEntry),
				}, func(b bool) {
					if !b {
						return
					}

					price, cost, _ := ConvertString(priceEntry.Text, costEntry.Text, "")
					quantity := ProcessQuantity(invenEntry.Text, w)

					// Check for an open slot
					ID := uint16(len(Item))
					v, found := Item[ID]

					for found && v != nil {
						v, found = Item[ID]
						fmt.Println(ID, found, v != nil)
						if !found || v == nil {
							break
						}
						ID += 1
					}

					Item[ID] = &Entry{Price: price, Name: nameEntry.Text, Quantity: [3]float32{quantity, 0, 0}, Cost: [3]float32{cost, 0, 0}}
					fmt.Println("!Found, Adding: ", Item[ID])

					target = int(ID)

					UI.HandleErrorWindow(SaveData(), w)

					InventoryData.Set(ConvertItemKeys())
					UpdateInventoryDisplay(ID)
					inventoryList.Refresh()
				}, w)
			}),

			widget.NewButton("Damages", func() {
				if target == -1 {
					return
				}

				dialog.ShowConfirm("Are You sure?", "Are you sure "+Item[uint16(target)].Name+" is damaged?", func(b bool) {
					if !b {
						return
					}

					entry := UI.NewNumEntry("How many were damaged?")
					dialog.ShowForm("Damages", "Confirm", "Cancel", []*widget.FormItem{widget.NewFormItem("Amount", entry)}, func(b bool) {
						if !b {
							return
						}
						quantity, err := strconv.ParseFloat(entry.Text, 32)
						UI.HandleError(err)

						y, month, day := time.Now().Date()
						year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

						s := Sale{
							ID:       uint16(target),
							Price:    0,
							Cost:     Item[uint16(target)].Cost[0],
							Quantity: float32(quantity),
							Usr:      255,
							Day:      uint8(day),
							Month:    uint8(month),
							Year:     uint8(year),
						}

						BuyCart([]Sale{s}, 0)
						UpdateInventoryDisplay(uint16(target))
					}, w)

				}, w)
			}),

			widget.NewButton("Add", func() {
				if target == -1 {
					return
				}

				priceEntry := UI.NewNumEntry("Selling Price")
				costEntry := UI.NewNumEntry("Cost Price")
				invenEntry := UI.NewNumEntry("Added stock, if only adding a fraction add a space before the fraction")

				dialog.ShowForm("Add", "Done", "Cancel", []*widget.FormItem{
					widget.NewFormItem("ID", widget.NewLabel(strconv.Itoa(target))),
					widget.NewFormItem("Name", widget.NewLabel(Item[uint16(target)].Name)),
					widget.NewFormItem("Price", priceEntry),
					widget.NewFormItem("Cost", costEntry),
					widget.NewFormItem("Added Stock", invenEntry),
				}, func(b bool) {
					if !b {
						return
					}
					price, cost, _ := ConvertString(priceEntry.Text, costEntry.Text, "")
					quan := ProcessQuantity(invenEntry.Text, w)
					Item[uint16(target)].Price = price

					for i := 0; i < 3; i++ {
						if Item[uint16(target)].Cost[i] == cost {
							Item[uint16(target)].Quantity[i] += quan
							break
						}
						if Item[uint16(target)].Quantity[i] == 0 {
							Item[uint16(target)].Quantity[i] = quan
							Item[uint16(target)].Cost[i] = cost
							break
						}
					}
					UI.HandleErrorWindow(SaveData(), w)
					UpdateInventoryDisplay(uint16(target))
				}, w)
			}),

			widget.NewButton("Remove", func() {
				dialog.ShowConfirm("Are you sure?", "Are you sure you want to delete this?", func(b bool) {
					if !b || target == -1 {
						return
					}

					Item[uint16(target)].Name = string([]byte{216}) + Item[uint16(target)].Name

					InventoryData.Set(ConvertItemKeys())
					inventoryList.Refresh()
					target = -1

					UI.HandleErrorWindow(SaveData(), w)
				}, w)
			}),
		),

		container.NewVSplit(
			container.NewMax(inventoryList),

			widget.NewButton("Search", func() {
				searchBar := UI.NewSearchBar("Type Item Name", SearchInventory)

				dialog.ShowCustomConfirm("Scan Item", "Confirm", "Cancel", searchBar, func(confirmed bool) {
					if !confirmed {
						return
					}

					id := searchBar.Result()

					if UI.HandleKnownError(0, id < 0, w) {
						return
					}
					target = id

					UpdateInventoryDisplay(uint16(target))
				}, w)
			}),
		))
}
