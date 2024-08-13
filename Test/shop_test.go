package Test

import (
	"BronzeHermes/Database"
	"testing"
	"time"
)

var testCart []Database.Sale // Use 6 since that's the one that has proper data

func ResetTestCartAndQuantities() {
	resetTestItemsAndSales()
	testCart = []Database.Sale{ // Use 6 since that's the one that has proper data
		{ID: 6, Price: 12, Cost: 2, Quantity: 2},
		{ID: 6, Price: 13, Cost: 2, Quantity: 3},
		{ID: 6, Price: 14, Cost: 2, Quantity: 6},
		{ID: 6, Price: 15, Cost: 2, Quantity: 7},
		{ID: 6, Price: 16, Cost: 2, Quantity: 12},
		{ID: 6, Price: 16, Cost: 2, Quantity: 14},
		{ID: 6, Price: 16, Cost: 2, Quantity: 20},
	}
}

func TestBuyingNormal(t *testing.T) {
	ResetTestCartAndQuantities()
	Item := uint16(6)
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[0]}, 0)

	test_time := time.Now().Unix()

	expectedQuantities := [3]float32{1, 4, 7}
	for i, v := range expectedQuantities {
		if Database.Items[Item].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[Item].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[Item].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{2, 3, 4}
	for i, v := range expectedCosts {
		if Database.Items[Item].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[Item].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[Item].Cost[i], v)
		}
	}

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)
			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testCart[0].ID {
			t.Logf("Test: %d, DB: %d", testCart[0].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testCart[0].Price {
			t.Logf("Test: %f, DB: %f", testCart[0].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testCart[0].Cost {
			t.Logf("Test: %f, DB: %f", testCart[0].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testCart[0].Quantity {
			t.Logf("Test: %f, DB: %f", testCart[0].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}

func TestBuyingAllOfQuantity0(t *testing.T) {
	ResetTestCartAndQuantities()

	Item := uint16(6)
	Database.Items = testItems
	t.Log(Database.Sales)
	Database.BuyCart([]Database.Sale{testCart[1]}, 0)

	test_time := time.Now().Unix()

	expectedQuantities := [3]float32{4, 7, 0}
	for i, v := range expectedQuantities {
		if Database.Items[Item].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[Item].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[Item].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{3, 4, 0}
	for i, v := range expectedCosts {
		if Database.Items[Item].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[Item].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[Item].Cost[i], v)
		}
	}

	if len(Database.Sales) != 1 {
		t.Errorf("Invalid Size of Sales | have: %d, want: 1", len(Database.Sales))
	}

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)

			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testCart[1].ID {
			t.Logf("Test: %d, DB: %d", testCart[1].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testCart[1].Price {
			t.Logf("Test: %f, DB: %f", testCart[1].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testCart[1].Cost {
			t.Logf("Test: %f, DB: %f", testCart[1].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testCart[1].Quantity {
			t.Logf("Test: %f, DB: %f", testCart[1].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}

func TestBuying2Quantities(t *testing.T) {
	ResetTestCartAndQuantities()

	Item := uint16(6)
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[2]}, 0)

	test_time := time.Now().Unix()

	expectedQuantities := [3]float32{1, 7, 0}
	for i, v := range expectedQuantities {
		if Database.Items[Item].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[Item].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[Item].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{3, 4, 0}
	for i, v := range expectedCosts {
		if Database.Items[Item].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[Item].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[Item].Cost[i], v)
		}
	}

	testAnswers := []Database.Sale{
		{ID: 6, Price: 14, Cost: 2, Quantity: 3},
		{ID: 6, Price: 14, Cost: 3, Quantity: 3},
	}

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)

			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testAnswers[i].ID {
			t.Logf("Test: %d, DB: %d", testAnswers[i].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testAnswers[i].Price {
			t.Logf("Test: %f, DB: %f", testAnswers[2].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testAnswers[i].Cost {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testAnswers[i].Quantity {
			t.Logf("Test: %f, DB: %f", testCart[i].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}

func TestBuyingAllQuantities0n1(t *testing.T) {
	ResetTestCartAndQuantities()

	Item := uint16(6)
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[3]}, 0)

	test_time := time.Now().Unix()

	expectedQuantities := [3]float32{7, 0, 0}
	for i, v := range expectedQuantities {
		if Database.Items[Item].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[Item].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[Item].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{4, 0, 0}
	for i, v := range expectedCosts {
		if Database.Items[Item].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[Item].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[Item].Cost[i], v)
		}
	}

	testAnswers := []Database.Sale{
		{ID: 6, Price: 15, Cost: 2, Quantity: 3},
		{ID: 6, Price: 15, Cost: 3, Quantity: 4},
	}

	t.Log(Database.Sales)

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)

			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testAnswers[i].ID {
			t.Logf("Test: %d, DB: %d", testAnswers[i].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testAnswers[i].Price {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testAnswers[i].Cost {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testAnswers[i].Quantity {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}

func TestBuyingInto3rdQuantity(t *testing.T) {
	ResetTestCartAndQuantities()
	Item := uint16(6)

	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[4]}, 0)

	test_time := time.Now().Unix()

	expectedQuantities := [3]float32{2, 0, 0}
	for i, v := range expectedQuantities {
		if Database.Items[Item].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[Item].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[Item].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{4, 0, 0}
	for i, v := range expectedCosts {
		if Database.Items[Item].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[Item].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[Item].Cost[i], v)
		}
	}

	testAnswers := []Database.Sale{
		{ID: 6, Price: 16, Cost: 2, Quantity: 3},
		{ID: 6, Price: 16, Cost: 3, Quantity: 4},
		{ID: 6, Price: 16, Cost: 4, Quantity: 5},
	}

	t.Log(len(Database.Sales))

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)

			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testAnswers[i].ID {
			t.Logf("Test: %d, DB: %d", testAnswers[i].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testAnswers[i].Price {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testAnswers[i].Cost {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testAnswers[i].Quantity {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}

func TestBuyingAllQuantities(t *testing.T) {
	ResetTestCartAndQuantities()

	Item := uint16(6)
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[5]}, 0)

	test_time := time.Now().Unix()

	expectedQuantities := [3]float32{0, 0, 0}
	for i, v := range expectedQuantities {
		if Database.Items[Item].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[Item].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[Item].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{4, 0, 0}
	for i, v := range expectedCosts {
		if Database.Items[Item].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[Item].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[Item].Cost[i], v)
		}
	}

	testAnswers := []Database.Sale{
		{ID: 6, Price: 16, Cost: 2, Quantity: 3},
		{ID: 6, Price: 16, Cost: 3, Quantity: 4},
		{ID: 6, Price: 16, Cost: 4, Quantity: 7},
	}

	t.Log(len(Database.Sales))

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)

			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testAnswers[i].ID {
			t.Logf("Test: %d, DB: %d", testAnswers[i].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testAnswers[i].Price {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testAnswers[i].Cost {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testAnswers[i].Quantity {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}

func TestBuyingOverAllQuantities(t *testing.T) {
	ResetTestCartAndQuantities()

	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[6]}, 0)

	test_time := time.Now().Unix()

	expectedQuantities := [3]float32{-6, 0, 0}
	for i, v := range expectedQuantities {
		if Database.Items[6].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[6].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[6].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{4, 0, 0}
	for i, v := range expectedCosts {
		if Database.Items[6].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[6].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[6].Cost[i], v)
		}
	}

	testAnswers := []Database.Sale{
		{ID: 6, Price: 16, Cost: 2, Quantity: 3},
		{ID: 6, Price: 16, Cost: 3, Quantity: 4},
		{ID: 6, Price: 16, Cost: 4, Quantity: 13},
	}

	t.Log(len(Database.Sales))

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)

			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testAnswers[i].ID {
			t.Logf("Test: %d, DB: %d", testAnswers[i].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testAnswers[i].Price {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testAnswers[i].Cost {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testAnswers[i].Quantity {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}

func TestOverBuying1QuantityWithOnly1Quantity(t *testing.T) {
	resetTestItemsAndSales()
	ResetTestCartAndQuantities()

	id := uint16(3)

	Database.Items = testItems
	s := Database.Sale{
		ID:       id,
		Price:    5,
		Cost:     2,
		Quantity: 5,
	}
	Database.BuyCart([]Database.Sale{s}, 0)

	test_time := time.Now().Unix()

	expectedQuantities := [3]float32{-4, 0, 0}
	for i, v := range expectedQuantities {
		if Database.Items[id].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[id].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[id].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{2, 0, 0}
	for i, v := range expectedCosts {
		if Database.Items[id].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[4].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[id].Cost[i], v)
		}
	}

	testAnswers := []Database.Sale{
		{ID: id, Price: 5, Cost: 2, Quantity: 5},
	}

	t.Log(len(Database.Sales))

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)

			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testAnswers[i].ID {
			t.Logf("Test: %d, DB: %d", testAnswers[i].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testAnswers[i].Price {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testAnswers[i].Cost {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testAnswers[i].Quantity {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}

func TestOverBuyingEmptyQuantity(t *testing.T) {
	ResetTestCartAndQuantities()
	Item := uint16(2)

	Database.Items = testItems
	s := Database.Sale{
		ID:       Item,
		Price:    5,
		Cost:     0,
		Quantity: 5,
	}
	Database.BuyCart([]Database.Sale{s}, 0)

	expectedQuantities := [3]float32{-5, 0, 0}
	for i, v := range expectedQuantities {
		if Database.Items[Item].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[Item].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[Item].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{0, 0, 0}
	for i, v := range expectedCosts {
		if Database.Items[Item].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[Item].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[Item].Cost[i], v)
		}
	}

	testAnswers := []Database.Sale{
		{ID: Item, Price: 5, Cost: 0, Quantity: 5},
	}

	test_time := time.Now().Unix()

	t.Log(Database.Sales[0].Cost)

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)

			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testAnswers[i].ID {
			t.Logf("Test: %d, DB: %d", testAnswers[i].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testAnswers[i].Price {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testAnswers[i].Cost {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testAnswers[i].Quantity {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}

func TestOverBuying2QuantityWith1Empty(t *testing.T) {
	ResetTestCartAndQuantities()
	Item := uint16(4)

	Database.Items = testItems
	s := Database.Sale{
		ID:       Item,
		Price:    5,
		Cost:     0,
		Quantity: 7,
	}
	Database.BuyCart([]Database.Sale{s}, 0)

	test_time := time.Now().Unix()

	expectedQuantities := [3]float32{-2, 0, 0}
	for i, v := range expectedQuantities {
		if Database.Items[Item].Quantity[i] != v {
			t.Logf("Quantity: %v", Database.Items[Item].Quantity)
			t.Errorf("Current and expected quantites don't match | have: %f, want: %f", Database.Items[Item].Quantity[i], v)
		}
	}

	expectedCosts := [3]float32{2, 0, 0}
	for i, v := range expectedCosts {
		if Database.Items[Item].Cost[i] != v {
			t.Logf("Cost: %v", Database.Items[Item].Cost)
			t.Errorf("Current and expected Cost don't match | have: %f, want: %f", Database.Items[Item].Cost[i], v)
		}
	}

	testAnswers := []Database.Sale{
		{ID: Item, Price: 5, Cost: 1, Quantity: 2},
		{ID: Item, Price: 5, Cost: 2, Quantity: 5},
	}

	t.Log(len(Database.Sales))

	for i, v := range Database.Sales {
		if test_time != v.Timestamp {
			t.Logf("Test: %d | DB: %d\n", test_time, v.Timestamp)

			t.Error("Dates are not matching for:", i)
		}

		if v.Customer != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Customer)
			t.Error("Wrong Customer for:", i)
		}

		if v.ID != testAnswers[i].ID {
			t.Logf("Test: %d, DB: %d", testAnswers[i].ID, v.ID)
			t.Error("Unequal IDs for:", i)
		}

		if v.Price != testAnswers[i].Price {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Price, v.Price)
			t.Error("Unequal Price for:", i)
		}

		if v.Cost != testAnswers[i].Cost {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Cost, v.Cost)
			t.Error("Unequal Cost for:", i)
		}

		if v.Quantity != testAnswers[i].Quantity {
			t.Logf("Test: %f, DB: %f", testAnswers[i].Quantity, v.Quantity)
			t.Error("Unequal Quantities for:", i)
		}
	}
}
