package Test

import (
	"BronzeHermes/Database"
	"testing"
)

var testSales = []Database.Sale{
	{ID: 0, Price: 0, Cost: 0, Quantity: 0, Customer: 0},
	{ID: 1, Price: 2, Cost: 1, Quantity: 15, Customer: 1},
	{ID: 12342, Price: 12, Cost: 0, Quantity: 0, Customer: 0},
	{ID: 0, Price: -41, Cost: 0, Quantity: 0, Customer: 0},
	{ID: 0, Price: 0, Cost: 0, Quantity: -4, Customer: 0},
	{ID: 0, Price: 3, Cost: -12, Quantity: 0, Customer: 0},
}

var testItems []Database.Item

func resetTestItemsAndSales() {
	testItems = []Database.Item{
		{Name: " ", Price: 0, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
		{Name: "Viva", Price: -1, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
		{Name: "Val", Price: 1, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
		{Name: "Pop", Price: 1, Cost: [3]float32{2, 0, 0}, Quantity: [3]float32{1, 0, 0}},
		{Name: "Villianous", Price: 0, Cost: [3]float32{1, 2, 0}, Quantity: [3]float32{2, 3, 0}},
		{Name: "Carty", Price: 12, Cost: [3]float32{2, 3, 4}, Quantity: [3]float32{3, 4, 7}},
		{Name: "Pop Daddy", Price: 12, Cost: [3]float32{2, 3, 4}, Quantity: [3]float32{3, 4, 7}},
		{Name: "Bila", Price: 1, Cost: [3]float32{12, 0, 0}, Quantity: [3]float32{-12, 0, 0}},
	}

	Database.Items = testItems

	Database.Sales = []Database.Sale{}
}

func TestSaveCustomers(t *testing.T) {
	Database.Customers = testCustomer
	Database.SaveNLoadCustomers()
	if len(Database.Customers) != len(testCustomer) {
		t.Error("Customers not equal after saving & loading")
	}

	for i := range testCustomer {
		if testCustomer[i] != Database.Customers[i] {
			t.Log("Test: ", testCustomer[i], " DB: ", Database.Customers[i])
			t.Error("Test & Customer don't match for ", i)
		}
	}
}

func TestSaveSales(t *testing.T) {
	Database.Sales = testSales
	Database.SaveNLoadSales()
	if len(testSales) != len(Database.Sales) {
		t.Error("Sales are unequal in length after saving & loading")
	}

	for i := range testSales {
		if testSales[i] != Database.Sales[i] {
			t.Log("Test: ", testSales[i], " DB: ", Database.Sales[i])
			t.Error("Test & Sales !equal for ", i)
		}
	}
}

func TestSaveItems(t *testing.T) {
	Database.Items = testItems
	Database.SaveNLoadKV()
	if len(testItems) != len(Database.Items) {
		t.Error("Unequal Lengths")
	}

	for i := range testItems {
		if testItems[i] != Database.Items[i] {
			t.Log("Test: ", testItems[i], " DB: ", Database.Items[i])
			t.Error("Test & Items !equal for ", i)
		}
	}
}

func TestSaveBackUp(t *testing.T) {
	Database.Items = testItems
	Database.Sales = testSales
	Database.Customers = testCustomer

	Database.SaveNLoadBackUp()

	if len(testItems) != len(Database.Items) {
		t.Error("Unequal Lengths")
	}

	for i := range testItems {
		if testItems[i] != Database.Items[i] {
			t.Log("Test: ", testItems[i], " DB: ", Database.Items[i])
			t.Error("Test & Items !equal for ", i)
		}
	}

	if len(testSales) != len(Database.Sales) {
		t.Error("Sales are unequal in length after saving & loading")
	}

	for i := range testSales {
		if testSales[i] != Database.Sales[i] {
			t.Log("Test: ", testSales[i], " DB: ", Database.Sales[i])
			t.Error("Test & Sales !equal for ", i)
		}
	}

	if len(Database.Customers) != len(testCustomer) {
		t.Error("Customers not equal after saving & loading")
	}

	for i := range testCustomer {
		if testCustomer[i] != Database.Customers[i] {
			t.Log("Test: ", testCustomer[i], " DB: ", Database.Customers[i])
			t.Error("Test & Customer don't match for ", i)
		}
	}
}

// Fuzz_Testing Binary Patching Save Files
func FuzzBinaryPatchSavingItems(f *testing.F) {
	resetTestItemsAndSales()

	for i, testItem := range testItems {
		f.Add(i, testItem.Price, testItem.Cost[0], testItem.Cost[1], testItem.Cost[2], testItem.Quantity[0], testItem.Quantity[1], testItem.Quantity[2], testItem.Name)
	}

	Database.Items = make([]Database.Item, 255)

	f.Fuzz(func(t *testing.T, i uint8, price, cost, cost2, cost3, quantity, quantity2, quantity3 float32, name string) {
		Database.Items[i] = Database.Item{Price: price, Cost: [3]float32{cost, cost2, cost3}, Quantity: [3]float32{quantity, quantity2, quantity3}, Name: name}
		Database.SaveData()
		// overwrite line with data
		// load the entire document
		// check if the id is valid and if it matches the data
		// check if the values are valid
	})
}
