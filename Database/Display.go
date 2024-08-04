package Database

import (
	"BronzeHermes/Debug"
	"BronzeHermes/UI"
	unknown "BronzeHermes/Unknown"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var InventoryData binding.IntList

var updateInventoryDisplay func(id any)

var RefreshInventory func()

func ItemsMenu(w fyne.Window) fyne.CanvasObject {
	idLabel := widget.NewLabel("ID")
	nameLabel := widget.NewLabel("Name")
	priceLabel := widget.NewLabel("Price")
	costLabel := widget.NewLabel("Cost")
	stockLabel := widget.NewLabel("Stock")
	totalRev := widget.NewLabel("Total Possible Revenue")
	totalCost := widget.NewLabel("Total Cost")

	InventoryData = binding.NewIntList()
	InventoryData.Set(ConvertItemKeys())

	inventoryList := widget.NewListWithData(InventoryData, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, nil, widget.NewLabel("N"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {
		val, _ := item.(binding.Int).Get()
		i, ok := Items[uint16(val)]
		if ok {
			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(i.Name)
		}
	})

	RefreshInventory = func() { InventoryData.Set(ConvertItemKeys()) }

	target := -1

	updateInventoryDisplay = func(idx any) {
		var id int
		switch x := idx.(type) {
		case uint16:
			id = int(x)
		case int:
			id, _ = InventoryData.GetValue(x) // NOTE: Handle Err
			inventoryList.Unselect(x)
		}

		stock := Items[uint16(id)].Quantity[0] + Items[uint16(id)].Quantity[1] + Items[uint16(id)].Quantity[2]
		idLabel.SetText("ID: " + strconv.Itoa(id))
		nameLabel.SetText("Name: " + Items[uint16(id)].Name)
		priceLabel.SetText(fmt.Sprintf("Price: %1.2f", Items[uint16(id)].Price))
		stockLabel.SetText(fmt.Sprintf("Stock: %1.2f\n", stock))
		totalRev.SetText(fmt.Sprintf("Total Possible Revenue: %1.2f", Items[uint16(id)].Price*stock))

		txt := "Cost: "

		if Items[uint16(id)].Cost[0] == 0 {
			txt += "0.00 "
		} else {
			for i := 0; i < 3 && Items[uint16(id)].Cost[i] > 0; i++ {
				txt += fmt.Sprintf("%1.2f, ", Items[uint16(id)].Cost[i])
			}
		}

		costLabel.SetText(txt[:len(txt)-1])

		var costs float32
		for i, v := range Items[uint16(id)].Cost {
			costs += v * Items[uint16(id)].Quantity[i]
		}

		totalCost.SetText(fmt.Sprintf("Total Cost: %1.2f", costs))

		target = id
	}

	inventoryList.OnSelected = func(id widget.ListItemID) {
		updateInventoryDisplay(id)
	}

	return container.New(layout.NewGridLayout(2),
		container.NewVBox(
			widget.NewLabelWithStyle("Inventory Info", fyne.TextAlign(1), fyne.TextStyle{Bold: true}),
			idLabel,
			nameLabel,
			priceLabel,
			costLabel,
			stockLabel,
			totalRev,
			totalCost,

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

					ID, errID := CreateItem(nameEntry.Text, priceEntry.Text, costEntry.Text, invenEntry.Text)
					if Debug.HandleKnownError(errID, errID != Debug.Success, w) {
						return
					}
					target = int(ID)

					Debug.ShowError("Saving Data", SaveData(), w)

					InventoryData.Set(ConvertItemKeys())
					updateInventoryDisplay(ID)
					inventoryList.Refresh()
				}, w)
			}),

			widget.NewButton("Damages", func() {
				if target == -1 {
					return
				}

				dialog.ShowConfirm("Are You sure?", "Are you sure "+Items[uint16(target)].Name+" is damaged?", func(b bool) {
					if !b {
						return
					}

					entry := UI.NewNumEntry("How many were damaged?")
					dialog.ShowForm("Damages", "Confirm", "Cancel", []*widget.FormItem{widget.NewFormItem("Amount", entry)}, func(b bool) {
						if !b {
							return
						}
						errID := AddDamages(uint16(target), entry.Text)
						if Debug.HandleKnownError(errID, errID != Debug.Success, w) {
							return
						}
						Debug.ShowError("Saing Data", SaveData(), w)
						updateInventoryDisplay(uint16(target))
					}, w)

				}, w)
			}),

			widget.NewButton("Add", func() {
				if target == -1 {
					return
				}

				priceEntry := UI.NewNumEntry("Selling Price")
				costEntry := UI.NewNumEntry("Cost Price")
				invenEntry := UI.NewNumEntry("Amont of Stock Being Added")

				dialog.ShowForm("Add", "Done", "Cancel", []*widget.FormItem{
					widget.NewFormItem("ID", widget.NewLabel(strconv.Itoa(target))),
					widget.NewFormItem("Name", widget.NewLabel(Items[uint16(target)].Name)),
					widget.NewFormItem("Price", priceEntry),
					widget.NewFormItem("Cost", costEntry),
					widget.NewFormItem("Added Stock", invenEntry),
				}, func(b bool) {
					if !b {
						return
					}

					errID := AddItem(uint16(target), priceEntry.Text, costEntry.Text, invenEntry.Text)
					if Debug.HandleKnownError(errID, errID != Debug.Success, w) {
						return
					}

					Debug.ShowError("Saving Data", SaveData(), w)
					updateInventoryDisplay(uint16(target))
				}, w)
			}),

			widget.NewButton("Remove", func() {
				if target == -1 {
					return
				}
				dialog.ShowConfirm("Are you sure?", "Are you sure you want to delete "+Items[uint16(target)].Name+"?", func(b bool) {
					if !b {
						return
					}

					unknown.RemoveAndUpdate(&Items[uint16(target)].Name, func() {
						CleanUpDeadItems()
						InventoryData.Set(ConvertItemKeys())
						inventoryList.Refresh()
						target = -1
					})

					Debug.ShowError("Saving Data", SaveData(), w)
				}, w)
			}),
		),

		container.NewVSplit(
			container.NewMax(inventoryList),

			widget.NewButton("Search", func() {
				searchBar := UI.NewSearchBar("Item Name Here...", SearchInventory)

				dialog.ShowCustomConfirm("Scan Item", "Confirm", "Cancel", searchBar, func(confirmed bool) {
					if !confirmed {
						return
					}

					id := searchBar.Result()

					if Debug.HandleKnownError(0, id < 0, w) {
						return
					}
					target = id

					updateInventoryDisplay(uint16(target))
				}, w)
			}),
		))
}
