package Test

import (
	"BronzeHermes/Database"
	"strconv"
	"testing"
	"time"
)

var testCart = []Database.Sale{ // Use 6 since that's the one that has proper data
	{ID: 6, Price: 12, Cost: 5, Quantity: 2},
	{ID: 6, Price: 13, Cost: 6, Quantity: 3},
	{ID: 6, Price: 14, Cost: 7, Quantity: 6},
	{ID: 6, Price: 15, Cost: 8, Quantity: 7},
	{ID: 6, Price: 16, Cost: 9, Quantity: 12},
	{ID: 6, Price: 16, Cost: 9, Quantity: 14},
}

func TestBuyingNormal(t *testing.T) {
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[0]}, 0)

	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])
	// Check the stored quantities
	if Database.Items[6].Quantity[0] != 1 {
		t.Logf("Quantity[0]: %1.2f", Database.Items[6].Quantity[0])
		t.Error("Issue subtracting the stock of the testQuantity[6] when buying cart")
	}

	for i, v := range Database.Sales {
		if v.Year != uint8(year) || v.Month != uint8(month) || v.Day != uint8(day) {
			t.Logf("Test: %d, %d, %d | DB: %d, %d, %d", day, month, year, v.Day, v.Month, v.Year)
			t.Error("Dates are not matching for:", i)
		}

		if v.Usr != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Usr)
			t.Error("Wrong Usr Set for:", i)
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
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[1]}, 0)

	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

	if Database.Items[6].Quantity[0] != 4 {
		t.Logf("Quantities: %1.2f, %1.2f", Database.Items[6].Quantity[0], Database.Items[6].Quantity[1])
		t.Error("Issue subtracting the stock of the testQuantity[6] when buying cart")
	}

	for i, v := range Database.Sales {
		if v.Year != uint8(year) || v.Month != uint8(month) || v.Day != uint8(day) {
			t.Logf("Test: %d, %d, %d | DB: %d, %d, %d", day, month, year, v.Day, v.Month, v.Year)
			t.Error("Dates are not matching for:", i)
		}

		if v.Usr != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Usr)
			t.Error("Wrong Usr Set for:", i)
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
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[2]}, 0)

	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

	if Database.Items[6].Quantity[0] != 1 {
		t.Logf("Sale Quantity: %1.2f", testCart[2].Quantity)
		t.Logf("Quantities: %1.2f, %1.2f", Database.Items[6].Quantity[0], Database.Items[6].Quantity[1])
		t.Error("Issue subtracting the stock of the testQuantity[6] when buying cart")
	}

	testAnswers := []Database.Sale{
		{ID: 6, Price: 14, Cost: 3, Quantity: 3},
		{ID: 6, Price: 14, Cost: 7, Quantity: 3},
	}

	for i, v := range Database.Sales {
		if v.Year != uint8(year) || v.Month != uint8(month) || v.Day != uint8(day) {
			t.Logf("Test: %d, %d, %d | DB: %d, %d, %d", day, month, year, v.Day, v.Month, v.Year)
			t.Error("Dates are not matching for:", i)
		}

		if v.Usr != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Usr)
			t.Error("Wrong Usr Set for:", i)
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
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[3]}, 0)

	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

	if Database.Items[6].Quantity[0] != 7 {
		t.Logf("Sale Quantity: %1.2f", testCart[3].Quantity)
		t.Logf("Quantities: %1.2f, %1.2f, %1.2f", Database.Items[6].Quantity[0], Database.Items[6].Quantity[1], Database.Items[6].Quantity[2])
		t.Error("Issue subtracting the stock of the testQuantity[6] when buying cart")
	}

	testAnswers := []Database.Sale{
		{ID: 6, Price: 15, Cost: 3, Quantity: 4},
		{ID: 6, Price: 15, Cost: 8, Quantity: 3},
	}

	t.Log(Database.Sales)

	for i, v := range Database.Sales {
		if v.Year != uint8(year) || v.Month != uint8(month) || v.Day != uint8(day) {
			t.Logf("Test: %d, %d, %d | DB: %d, %d, %d", day, month, year, v.Day, v.Month, v.Year)
			t.Error("Dates are not matching for:", i)
		}

		if v.Usr != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Usr)
			t.Error("Wrong Usr Set for:", i)
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
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[4]}, 0)

	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

	if Database.Items[6].Quantity[0] != 2 {
		t.Logf("Sale Quantity: %1.2f", testCart[4].Quantity)
		t.Logf("Quantities: %1.2f, %1.2f, %1.2f", Database.Items[6].Quantity[0], Database.Items[6].Quantity[1], Database.Items[6].Quantity[2])
		t.Error("Issue subtracting the stock of the testQuantity[6] when buying cart")
	}

	testAnswers := []Database.Sale{
		{ID: 6, Price: 16, Cost: 4, Quantity: 5},
		{ID: 6, Price: 16, Cost: 3, Quantity: 4},
		{ID: 6, Price: 16, Cost: 9, Quantity: 3},
	}

	t.Log(len(Database.Sales))

	for i, v := range Database.Sales {
		if v.Year != uint8(year) || v.Month != uint8(month) || v.Day != uint8(day) {
			t.Logf("Test: %d, %d, %d | DB: %d, %d, %d", day, month, year, v.Day, v.Month, v.Year)
			t.Error("Dates are not matching for:", i)
		}

		if v.Usr != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Usr)
			t.Error("Wrong Usr Set for:", i)
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
	Database.Items = testItems
	Database.BuyCart([]Database.Sale{testCart[5]}, 0)

	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

	if Database.Items[6].Quantity[0] != 0 {
		t.Logf("Sale Quantity: %1.2f", testCart[4].Quantity)
		t.Logf("Quantities: %1.2f, %1.2f, %1.2f", Database.Items[6].Quantity[0], Database.Items[6].Quantity[1], Database.Items[6].Quantity[2])
		t.Error("Issue subtracting the stock of the testQuantity[6] when buying cart")
	}

	testAnswers := []Database.Sale{
		{ID: 6, Price: 16, Cost: 4, Quantity: 7},
		{ID: 6, Price: 16, Cost: 3, Quantity: 4},
		{ID: 6, Price: 16, Cost: 9, Quantity: 3},
	}

	t.Log(len(Database.Sales))

	for i, v := range Database.Sales {
		if v.Year != uint8(year) || v.Month != uint8(month) || v.Day != uint8(day) {
			t.Logf("Test: %d, %d, %d | DB: %d, %d, %d", day, month, year, v.Day, v.Month, v.Year)
			t.Error("Dates are not matching for:", i)
		}

		if v.Usr != 0 {
			t.Logf("Test: %d, DB: %d", 0, v.Usr)
			t.Error("Wrong Usr Set for:", i)
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
