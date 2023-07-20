package Test

import (
	"BronzeHermes/Database"
	"strconv"
	"testing"
	"time"
)

func TestRemoveSale(t *testing.T) {
	resetTestItemsAndSales()
	Database.Items = testItems
	Database.Sales = []Database.Sale{
		{ID: 6, Cost: 2, Quantity: 15},
	}

	Database.RemoveFromSales(0)

	// check the length of sales
	if len(Database.Sales) != 0 {
		t.Errorf("Item not removed | len: %d", len(Database.Sales))
	}
	// Check if item[6]'s quantity has increase for cost 2
	if Database.Items[6].Quantity[0] != 18 {
		t.Errorf("Error occured with Item's quantity[0] | have: %f", Database.Items[6].Quantity[0])
	}

	if Database.Items[6].Quantity[1] != 4 || Database.Items[6].Quantity[2] != 7 {
		t.Errorf("Error occured with other quantiites | want [1]: 4 & want [2]: 7, have: %v", Database.Items[6].Quantity)
	}
}

func TestRemoveSale2(t *testing.T) {
	resetTestItemsAndSales()
	Database.Items = testItems
	Database.Sales = []Database.Sale{
		{ID: 6, Cost: 3, Quantity: 15},
	}

	Database.RemoveFromSales(0)

	if len(Database.Sales) != 0 {
		t.Errorf("Item not removed | len: %d", len(Database.Sales))
	}

	if Database.Items[6].Quantity[1] != 19 {
		t.Errorf("Error occured with Item's quantity[0] | have: %f", Database.Items[6].Quantity[0])
	}

	if Database.Items[6].Quantity[0] != 3 || Database.Items[6].Quantity[2] != 7 {
		t.Errorf("Error occured with other quantiites | want [1]: 4 & want [2]: 7, have: %v", Database.Items[6].Quantity)
	}
}

func TestAddingDamages(t *testing.T) {
	resetTestItemsAndSales()

	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

	answer := Database.Sale{
		ID:       6,
		Price:    0,
		Cost:     Database.Items[6].Cost[0],
		Quantity: 2,
		Year:     uint8(year),
		Month:    uint8(month),
		Day:      uint8(day),
		Usr:      255,
	}

	errID := Database.AddDamages(6, "2")
	if errID != -1 {
		t.Errorf("Some error has occured | have: %d, want: -1", errID)
	}

	if len(Database.Sales) != 1 {
		t.Log(Database.Sales)
		t.Errorf("Issue adding the damages to sales | have: %d, want: 1", len(Database.Sales))
	}

	if Database.Sales[0] != answer {
		t.Errorf("Sales do not match up with the expected | have: %v, want: %v", Database.Sales[0], answer)
	}
}

func TestFailedAddingDamages(t *testing.T) {
	resetTestItemsAndSales()

	errID := Database.AddDamages(6, "-12.4.23.")
	if errID != 0 {
		t.Errorf("An Error has slipped through | have: %d, want: 0", errID)
	}

	if len(Database.Sales) != 0 {
		t.Log(Database.Sales)
		t.Errorf("Illegal Adding of Damages | have: %d, want: 0", len(Database.Sales))
	}
}

func TestSearchingInventoryEmpty(t *testing.T) {
	resetTestItemsAndSales()
	names, IDs := Database.SearchInventory("")
	if len(names) != len(testItems) {
		t.Errorf("Search Results !valid | have: %d, want: %d", len(names), len(testItems))
	}

	if len(names) != len(IDs) {
		t.Errorf("Len of IDs & Names != | names: %d, IDs: %d", len(names), len(IDs))
	}
	for k, v := range testItems {
		found := false
		for i, name := range names {
			if v.Name == name && k == IDs[i] {
				found = true
			} else if k == IDs[i] {
				t.Errorf("ID and Name don't share the same ID | have: %d & %s, want: %d & %s", IDs[i], name, k, v.Name)
			} else if v.Name == name {
				t.Errorf("ID and Name don't share the same ID | have: %d & %s, want: %d & %s", IDs[i], name, k, v.Name)
			}
		}
		if !found {
			t.Errorf("KV ! found | have: %s: %d", v.Name, k)
		}
	}
}

func TestSearchingInventory(t *testing.T) {
	resetTestItemsAndSales()
	names, IDs := Database.SearchInventory("V")
	if len(names) != 3 {
		t.Errorf("Search Results !valid | have: %d, want: 3", len(names))
	}

	if len(names) != len(IDs) {
		t.Errorf("Len of IDs & Names != | names: %d, IDs: %d", len(names), len(IDs))
	}

	expecrtedIDs := [3]uint16{1, 2, 5}
	expectedNames := [3]string{
		"Viva",
		"Val",
		"Villianous",
	}

	for i, ID := range expecrtedIDs {
		found := false

		for nI, name := range names {
			if expectedNames[i] == name && ID == IDs[nI] {
				found = true
			} else if expectedNames[i] == name {
				t.Errorf("ID and Name don't share the same ID | have: %d & %s, want: %d & %s", IDs[i], name, ID, expectedNames[i])
			} else if ID == IDs[nI] {
				t.Errorf("ID and Name don't share the same ID | have: %d & %s, want: %d & %s", IDs[i], name, ID, expectedNames[i])
			}

		}

		if !found {
			t.Errorf("KV ! found | have: %s: %d", expectedNames[i], ID)
		}
	}
}

func TestFailedSearchingInventory(t *testing.T) {
	resetTestItemsAndSales()
	names, IDs := Database.SearchInventory("PEEEPEE")
	if len(names) != 0 {
		t.Errorf("Search Results !valid | have: %d, want: 0", len(names))
	}

	if len(names) != len(IDs) {
		t.Errorf("Len of IDs & Names != | names: %d, IDs: %d", len(names), len(IDs))
	}

}

func TestSearchingInventoryWithTrailingSpace(t *testing.T) {
	resetTestItemsAndSales()
	names, IDs := Database.SearchInventory("V  ")
	if len(names) != 3 {
		t.Errorf("Search Results !valid | have: %d, want: 3", len(names))
	}

	if len(names) != len(IDs) {
		t.Errorf("Len of IDs & Names != | names: %d, IDs: %d", len(names), len(IDs))
	}

	expecrtedIDs := [3]uint16{1, 2, 5}
	expectedNames := [3]string{
		"Viva",
		"Val",
		"Villianous",
	}

	for i, ID := range expecrtedIDs {
		found := false

		for nI, name := range names {
			if expectedNames[i] == name && ID == IDs[nI] {
				found = true
			} else if expectedNames[i] == name {
				t.Errorf("ID and Name don't share the same ID | have: %d & %s, want: %d & %s", IDs[i], name, ID, expectedNames[i])
			} else if ID == IDs[nI] {
				t.Errorf("ID and Name don't share the same ID | have: %d & %s, want: %d & %s", IDs[i], name, ID, expectedNames[i])
			}

		}

		if !found {
			t.Errorf("KV ! found | have: %s: %d", expectedNames[i], ID)
		}
	}
}

func TestSearchingInventoryWithLeadingSpace(t *testing.T) {
	resetTestItemsAndSales()
	names, IDs := Database.SearchInventory("   V")
	if len(names) != 3 {
		t.Errorf("Search Results !valid | have: %d, want: 3", len(names))
	}

	if len(names) != len(IDs) {
		t.Errorf("Len of IDs & Names != | names: %d, IDs: %d", len(names), len(IDs))
	}

	expecrtedIDs := [3]uint16{1, 2, 5}
	expectedNames := [3]string{
		"Viva",
		"Val",
		"Villianous",
	}

	for i, ID := range expecrtedIDs {
		found := false

		for nI, name := range names {
			if expectedNames[i] == name && ID == IDs[nI] {
				found = true
			} else if expectedNames[i] == name {
				t.Errorf("ID and Name don't share the same ID | have: %d & %s, want: %d & %s", IDs[i], name, ID, expectedNames[i])
			} else if ID == IDs[nI] {
				t.Errorf("ID and Name don't share the same ID | have: %d & %s, want: %d & %s", IDs[i], name, ID, expectedNames[i])
			}

		}

		if !found {
			t.Errorf("KV ! found | have: %s: %d", expectedNames[i], ID)
		}
	}
}

func TestSearchingInventoryWithSpace(t *testing.T) {
	resetTestItemsAndSales()
	names, IDs := Database.SearchInventory("Pop D")
	if len(names) != 1 {
		t.Errorf("Search Results !valid | have: %d, want: 1", len(names))
	}

	if len(names) != len(IDs) {
		t.Errorf("Len of IDs & Names != | names: %d, IDs: %d", len(names), len(IDs))
	}

	if names[0] != "Pop Daddy" {
		t.Errorf("Name is not equal to it's expected value | have: %s, want: Pop Daddy", names[0])
	}

	if IDs[0] != 7 {
		t.Errorf("ID is not equal to it's expected value | have: %d, want: 7", IDs[0])
	}
}
