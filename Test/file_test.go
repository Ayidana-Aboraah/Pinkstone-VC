package Test

import (
	"BronzeHermes/Database"
	"testing"
)

var testUserNames = []string{
	"Bob",
	"Penny",
	"Poppy",
	"123481023984012983",
}

var testSales = []Database.Sale{
	{ID: 0, Price: 0, Cost: 0, Quantity: 0, Customer: 0, Usr: 0},
	{ID: 1, Price: 2, Cost: 1, Quantity: 15, Customer: 1, Usr: 2},
	{ID: 12342, Price: 12, Cost: 0, Quantity: 0, Customer: 0, Usr: 0},
	{ID: 0, Price: -41, Cost: 0, Quantity: 0, Customer: 0, Usr: 0},
	{ID: 0, Price: 0, Cost: 0, Quantity: -4, Customer: 0, Usr: 0},
	{ID: 0, Price: 3, Cost: -12, Quantity: 0, Customer: 0, Usr: 0},
}

var testItems = map[uint16]*Database.Entry{
	0: {Name: "", Price: 0, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
	1: {Name: "Viva", Price: -1, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
	2: {Name: "Val", Price: 1, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
	4: {Name: "Pop", Price: 1, Cost: [3]float32{2, 0, 0}, Quantity: [3]float32{1, 0, 0}},
	5: {Name: "Villianous", Price: 0, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
	6: {Name: "Carty", Price: 12, Cost: [3]float32{2, 3, 4}, Quantity: [3]float32{3, 4, 7}},
}

func resetTestItemsAndSales() {
	testItems = map[uint16]*Database.Entry{
		0: {Name: "", Price: 0, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
		1: {Name: "Viva", Price: -1, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
		2: {Name: "Val", Price: 1, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
		4: {Name: "Pop", Price: 1, Cost: [3]float32{2, 0, 0}, Quantity: [3]float32{1, 0, 0}},
		5: {Name: "Villianous", Price: 0, Cost: [3]float32{0, 0, 0}, Quantity: [3]float32{0, 0, 0}},
		6: {Name: "Carty", Price: 12, Cost: [3]float32{2, 3, 4}, Quantity: [3]float32{3, 4, 7}},
	}

	Database.Sales = []Database.Sale{}
}

func TestSaveUsers(t *testing.T) {
	Database.Users = testUserNames
	Database.SaveNLoadUsers()
	if len(Database.Users) != len(testUserNames) {
		t.Errorf("Unequal lengths after saving")
	}

	for i := range testUserNames {
		if testUserNames[i] != Database.Users[i] {
			t.Errorf("%s != %s | idx: %d of testUsers, error after lodaing data", testUserNames[i], Database.Users[i], i)
		}
	}
}

func TestSaveNoUser(t *testing.T) {
	Database.Users = []string{}
	Database.SaveNLoadUsers()
	if len(Database.Users) != 0 {
		t.Error("Users Databasse lenght is != 0, when handed an empty array")
	}
}

func TestSaveBlankUser(t *testing.T) {
	Database.Users = []string{""}
	Database.SaveNLoadUsers()
	if len(Database.Users) != 1 {
		t.Error("Users Databasse lenght is != 1, when handed an empty string")
	}
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
	Database.Users = testUserNames
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

	if len(Database.Users) != len(testUserNames) {
		t.Errorf("Unequal lengths after saving")
	}

	for i := range testUserNames {
		if testUserNames[i] != Database.Users[i] {
			t.Errorf("%s != %s | idx: %d of testUsers, error after lodaing data", testUserNames[i], Database.Users[i], i)
		}
	}
}
