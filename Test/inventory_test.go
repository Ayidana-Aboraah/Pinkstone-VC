package Test

import (
	"BronzeHermes/Database"
	"strings"
	"testing"
)

var testCustomer = []string{
	"",
	"xxx",
	"128340192830",
	"Ball",
	"Broom",
	"Ballon",
	"Barrow",
	"aa",
}

func TestCustomerSearchWithResults(t *testing.T) {
	Database.Customers = testCustomer
	res, ids := Database.SearchCustomers("Ba")
	for i, v := range res {
		if !strings.Contains(strings.ToLower(v), "ba") {
			t.Log(v)
			t.Error("Somehow the results have nothing to do with the input")
		}
		if v != testCustomer[ids[i]] {
			t.Log(v)
			t.Log(testCustomer[ids[i]])
			t.Error("Either the ids or the string results don't match up")
		}
	}
}

func TestCustomerSearchEmpty(t *testing.T) {
	Database.Customers = testCustomer
	res, ids := Database.SearchCustomers("")
	if len(ids) != len(res) {
		t.Log(ids)
		t.Log(res)
		t.Errorf("string & idx arrays are != | idx: %d, str: %d", len(ids), len(res))
	}

	if len(ids) != len(testCustomer) {
		t.Log(ids)
		t.Log(testCustomer)
		t.Errorf("len of ids & customers are unequal | idx: %d, customers: %d", len(ids), len(Database.Customers))
	}

	for i, v := range testCustomer {
		if v != res[i] {
			t.Errorf("Customer: %d are not equal in both the search result & the Customer list", i)
		}
	}
}

func TestProcessQuantity(t *testing.T) {
	quantity, errID := Database.ProcessQuantity("15")
	if errID != -1 {
		t.Errorf("Some error has occured, have: %d", errID)
	}

	if quantity != 15 {
		t.Errorf("Quantity !=, have: %f, want: 15.0", quantity)
	}
}

func TestProcessQuantityWithError(t *testing.T) {
	quantity, errID := Database.ProcessQuantity("-12-42-321")
	if errID != 0 {
		t.Errorf("err not caught or wrong error passed, have: %d", errID)
	}

	if quantity != 0 {
		t.Errorf("Quantity !=, have: %f, want: 0.0", quantity)
	}
}

func TestProcessFractionQuantity(t *testing.T) {
	quantity, errID := Database.ProcessQuantity(" 1/2")
	if errID != -1 {
		t.Errorf("Some error has occured, have: %d", errID)
	}

	if quantity != 1.0/2.0 {
		t.Errorf("Quantity !=, have: %f, want: %f", quantity, 1.0/2.0)
	}
}

func TestProcessInvalidFraction(t *testing.T) {
	quantity, errID := Database.ProcessQuantity("1/2")
	if errID != 0 {
		t.Errorf("Some error has occured, have: %d", errID)
	}

	if quantity != 0 {
		t.Errorf("Quantity !=, have: %f, want: %f", quantity, 0.0)
	}
}

func TestProcessQuantityWithFraction(t *testing.T) {
	quantity, errID := Database.ProcessQuantity("12 1/2")
	if errID != -1 {
		t.Errorf("Some error has occured, have: %d", errID)
	}

	if quantity != 12.0+1.0/2.0 {
		t.Errorf("Quantity !=, have: %f, want: %f", quantity, 12.0+1.0/2.0)
	}
}

func TestNewItem(t *testing.T) {
	Database.Items = map[uint16]*Database.Entry{}
	id, errID := Database.CreateItem("Piss", "15", "12", "95")
	if errID != -1 {
		t.Errorf("Some error has occured | ErrID: %d", errID)
	}
	val, found := Database.Items[id]
	if !found {
		t.Error("Created Item not found")
	}

	if len(Database.Items) != 1 {
		t.Errorf("Multiple Items Illegally Created, len: %d", len(Database.Items))
	}

	if val.Name != "Piss" {
		t.Errorf("Issue adding name to the created Item | have: %s, want: Piss", val.Name)
	}

	if val.Price != 15.0 {
		t.Errorf("Wrong Price | have: %f, want: %f", val.Price, 15.0)
	}

	if val.Cost[0] != 12.0 {
		t.Errorf("Wrong Cost | have: %f, want: %f", val.Cost[0], 12.0)
	}

	if val.Quantity[0] != 95.0 {
		t.Errorf("Wrong Quantity | have: %f, want: %f", val.Quantity[0], 95.0)
	}

	if val.Cost[1] != 0 || val.Cost[2] != 0 {
		t.Errorf("Other Costs filled Illegally | have: %v, want: [12.0, 0.0, 0.0]", val.Cost)
	}

	if val.Quantity[1] != 0 || val.Quantity[2] != 0 {
		t.Errorf("Other Quantities filled Illegally | have: %v, want: [95.0, 0.0, 0.0]", val.Quantity)
	}
}

func TestNewItemWithFraction(t *testing.T) {
	Database.Items = map[uint16]*Database.Entry{}
	id, errID := Database.CreateItem("Piss", "15", "12", "95 3/4")
	if errID != -1 {
		t.Errorf("Some error has occured | ErrID: %d", errID)
	}
	val, found := Database.Items[id]
	if !found {
		t.Error("Created Item not found")
	}

	if len(Database.Items) != 1 {
		t.Errorf("Multiple Items Illegally Created, len: %d", len(Database.Items))
	}

	if val.Name != "Piss" {
		t.Errorf("Issue adding name to the created Item | have: %s, want: Piss", val.Name)
	}

	if val.Price != 15.0 {
		t.Errorf("Wrong Price | have: %f, want: %f", val.Price, 15.0)
	}

	if val.Cost[0] != 12.0 {
		t.Errorf("Wrong Cost | have: %f, want: %f", val.Cost[0], 12.0)
	}

	if val.Quantity[0] != 95.0+3.0/4.0 {
		t.Errorf("Wrong Quantity | have: %f, want: %f", val.Quantity[0], 95.0+3.0/4.0)
	}

	if val.Cost[1] != 0 || val.Cost[2] != 0 {
		t.Errorf("Other Costs filled Illegally | have: %v, want: [12.0, 0.0, 0.0]", val.Cost)
	}

	if val.Quantity[1] != 0 || val.Quantity[2] != 0 {
		t.Errorf("Other Quantities filled Illegally | have: %v, want: [95.75, 0.0, 0.0]", val.Quantity)
	}
}

func TestFailedItemCreation(t *testing.T) {
	Database.Items = map[uint16]*Database.Entry{}
	id, errID := Database.CreateItem("Piss", "15", "12", "953/4")
	if errID != 0 {
		t.Errorf("An error has slipped through | have: %d, want: 0", errID)
	}

	val, found := Database.Items[id]
	if found {
		t.Error("Illegal Item Creation")
	}
	if val != nil {
		t.Errorf("Illegal Value Creation | have: %v", val)
	}
}

func TestAddItem(t *testing.T) {
	resetTestItems()
	Database.Items = testItems
	errID := Database.AddItem(0, "1", "12", "32")
	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	val, found := Database.Items[0]
	if val == nil || !found {
		t.Errorf("val: An error occured retrieving the value")
		t.FailNow()
	}

	if val.Name != "" {
		t.Errorf("Error occured with the name | have: %s, want: '' ", val.Name)
	}

	if val.Price != 1 {
		t.Errorf("Error occured with price | have: %f, want: 1.0", val.Price)
	}

	if val.Cost[0] != 12.0 {
		t.Errorf("Error occured with cost | have: %f, want: 12.0", val.Cost[0])
	}

	if val.Quantity[0] != 32.0 {
		t.Errorf("Error occured with cost | have: %f, want: 32.0", val.Quantity[0])
	}
}

func TestAddWithFraction(t *testing.T) {
	resetTestItems()
	Database.Items = testItems
	errID := Database.AddItem(0, "1", "12", "32 1/2")
	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	val, found := Database.Items[0]
	if val == nil || !found {
		t.Errorf("val: An error occured retrieving the value")
		t.FailNow()
	}

	if val.Name != "" {
		t.Errorf("Error occured with the name | have: %s, want: '' ", val.Name)
	}

	if val.Price != 1 {
		t.Errorf("Error occured with price | have: %f, want: 1.0", val.Price)
	}

	if val.Cost[0] != 12.0 {
		t.Errorf("Error occured with cost | have: %f, want: 12.0", val.Cost[0])
	}

	if val.Quantity[0] != 32.5 {
		t.Errorf("Error occured with cost | have: %f, want: 32.5", val.Quantity[0])
	}
}

func TestFailedAdd(t *testing.T) {
	resetTestItems()
	Database.Items = testItems
	errID := Database.AddItem(0, "3", "12", "--32 -1-/-0.9-")
	if errID != 0 {
		t.Errorf("An Error has slipped through | have: %d, want: 0", errID)
	}

	val, found := Database.Items[0]
	if val == nil || !found {
		t.Errorf("val: An error occured retrieving the value")
		t.FailNow()
	}

	if val.Name != "" {
		t.Errorf("Error occured with the name | have: %s, want: '' ", val.Name)
	}

	if val.Price != 0.0 {
		t.Errorf("Error occured with price | have: %f, want: 0.0", val.Price)
	}

	if val.Cost[0] != 0.0 {
		t.Errorf("Error occured with cost | have: %f, want: 0.0", val.Cost[0])
	}

	if val.Quantity[0] != 0.0 {
		t.Errorf("Error occured with cost | have: %f, want: 0.0", val.Quantity[0])
	}
}

func TestAdd2Costs(t *testing.T) {
	resetTestItems()
	Database.Items = testItems
	errID := Database.AddItem(4, "1", "12", "32")
	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	val, found := Database.Items[4]
	if val == nil || !found {
		t.Errorf("val: An error occured retrieving the value")
		t.FailNow()
	}

	if val.Name != "Pop" {
		t.Errorf("Error occured with the name | have: %s, want: '' ", val.Name)
	}

	if val.Price != 1 {
		t.Errorf("Error occured with price | have: %f, want: 1.0", val.Price)
	}

	if val.Cost[0] != 2.0 {
		t.Errorf("Error occured with cost | have: %f, want: 2.0", val.Cost[0])
	}

	if val.Cost[1] != 12.0 {
		t.Errorf("Error occured with cost | have: %f, want: 12.0", val.Cost[1])
	}

	if val.Quantity[0] != 1.0 {
		t.Errorf("Error occured with cost | have: %f, want: 1.0", val.Quantity[0])
	}

	if val.Quantity[1] != 32.0 {
		t.Errorf("Error occured with cost | have: %f, want: 32.0", val.Quantity[1])
	}
}

func TestAdd2CostsWithFraction(t *testing.T) {
	resetTestItems()
	Database.Items = testItems
	errID := Database.AddItem(4, "1", "12", "32 1/2")
	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	val, found := Database.Items[4]
	if val == nil || !found {
		t.Errorf("val: An error occured retrieving the value")
		t.FailNow()
	}

	if val.Name != "Pop" {
		t.Errorf("Error occured with the name | have: %s, want: '' ", val.Name)
	}

	if val.Price != 1 {
		t.Errorf("Error occured with price | have: %f, want: 1.0", val.Price)
	}

	if val.Cost[0] != 2.0 {
		t.Errorf("Error occured with cost | have: %f, want: 2.0", val.Cost[0])
	}

	if val.Cost[1] != 12.0 {
		t.Errorf("Error occured with cost | have: %f, want: 12.0", val.Cost[1])
	}

	if val.Quantity[0] != 1.0 {
		t.Errorf("Error occured with cost | have: %f, want: 1.0", val.Quantity[0])
	}

	if val.Quantity[1] != 32.5 {
		t.Errorf("Error occured with cost | have: %f, want: 32.5", val.Quantity[1])
	}
}

func TestAdd3Costs(t *testing.T) {
	resetTestItems()
	Database.Items = testItems
	errID := Database.AddItem(4, "1", "12", "32")

	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	errID = Database.AddItem(4, "1", "13", "31")

	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	val, found := Database.Items[4]
	if val == nil || !found {
		t.Errorf("val: An error occured retrieving the value")
		t.FailNow()
	}

	if val.Name != "Pop" {
		t.Errorf("Error occured with the name | have: %s, want: '' ", val.Name)
	}

	if val.Price != 1 {
		t.Errorf("Error occured with price | have: %f, want: 1.0", val.Price)
	}

	if val.Cost[0] != 2.0 {
		t.Errorf("Error occured with cost | have: %f, want: 2.0", val.Cost[0])
	}

	if val.Cost[1] != 12.0 {
		t.Errorf("Error occured with cost | have: %f, want: 12.0", val.Cost[1])
	}

	if val.Cost[2] != 13.0 {
		t.Errorf("Error occured with cost | have: %f, want: 13.0", val.Cost[2])
	}

	if val.Quantity[0] != 1.0 {
		t.Errorf("Error occured with cost | have: %f, want: 1.0", val.Quantity[0])
	}

	if val.Quantity[1] != 32 {
		t.Errorf("Error occured with cost | have: %f, want: 32", val.Quantity[1])
	}

	if val.Quantity[2] != 31 {
		t.Errorf("Error occured with cost | have: %f, want: 31", val.Quantity[2])
	}
}

func TestAdd3CostsWithFraction(t *testing.T) {
	resetTestItems()
	Database.Items = testItems
	errID := Database.AddItem(4, "1", "12", "32 1/2")

	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	errID = Database.AddItem(4, "1", "13", "32 3/4")

	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	val, found := Database.Items[4]
	if val == nil || !found {
		t.Errorf("val: An error occured retrieving the value")
		t.FailNow()
	}

	if val.Name != "Pop" {
		t.Errorf("Error occured with the name | have: %s, want: '' ", val.Name)
	}

	if val.Price != 1 {
		t.Errorf("Error occured with price | have: %f, want: 1.0", val.Price)
	}

	if val.Cost[0] != 2.0 {
		t.Errorf("Error occured with cost | have: %f, want: 2.0", val.Cost[0])
	}

	if val.Cost[1] != 12.0 {
		t.Errorf("Error occured with cost | have: %f, want: 12.0", val.Cost[1])
	}

	if val.Cost[2] != 13.0 {
		t.Errorf("Error occured with cost | have: %f, want: 13.0", val.Cost[2])
	}

	if val.Quantity[0] != 1.0 {
		t.Errorf("Error occured with cost | have: %f, want: 1.0", val.Quantity[0])
	}

	if val.Quantity[1] != 32.5 {
		t.Errorf("Error occured with cost | have: %f, want: 32.5", val.Quantity[1])
	}

	if val.Quantity[2] != 32.75 {
		t.Errorf("Error occured with cost | have: %f, want: 32.75", val.Quantity[2])
	}
}

func TestFailedAdd4CostsWithFraction(t *testing.T) {
	resetTestItems()
	Database.Items = testItems
	errID := Database.AddItem(4, "1", "12", "32 1/2")

	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	errID = Database.AddItem(4, "1", "13", "32 3/4")

	if errID != -1 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	errID = Database.AddItem(4, "1", "15", "3")

	if errID != 2 {
		t.Errorf("An Error has occured | have: %d, want: -1", errID)
	}

	val, found := Database.Items[4]
	if val == nil || !found {
		t.Errorf("val: An error occured retrieving the value")
		t.FailNow()
	}

	if val.Name != "Pop" {
		t.Errorf("Error occured with the name | have: %s, want: '' ", val.Name)
	}

	if val.Price != 1 {
		t.Errorf("Error occured with price | have: %f, want: 1.0", val.Price)
	}

	if val.Cost[0] != 2.0 {
		t.Errorf("Error occured with cost | have: %f, want: 2.0", val.Cost[0])
	}

	if val.Cost[1] != 12.0 {
		t.Errorf("Error occured with cost | have: %f, want: 12.0", val.Cost[1])
	}

	if val.Cost[2] != 13.0 {
		t.Errorf("Error occured with cost | have: %f, want: 13.0", val.Cost[2])
	}

	if val.Quantity[0] != 1.0 {
		t.Errorf("Error occured with cost | have: %f, want: 1.0", val.Quantity[0])
	}

	if val.Quantity[1] != 32.5 {
		t.Errorf("Error occured with cost | have: %f, want: 32.5", val.Quantity[1])
	}

	if val.Quantity[2] != 32.75 {
		t.Errorf("Error occured with cost | have: %f, want: 32.75", val.Quantity[2])
	}
}
