package test

import (
	"BronzeHermes/Database"
	"testing"
)

var TestDB = [3][]Database.Sale{
	{
		{ID: 674398202423, Price: 234.23, Cost: 1324, Quantity: 1},
		{ID: 674398202423, Price: 100.50, Cost: 1324, Quantity: 1},
		{ID: 389432143927, Price: 3974.89, Cost: 8934.24, Quantity: 5},
		{ID: 674398202423, Price: 90109.22, Cost: 48.24, Quantity: 87},
		{ID: 402933466372, Price: 1324.89, Cost: 21432.24, Quantity: 4124},
		{ID: 198998421024, Price: 1094.89, Cost: 9021038.24, Quantity: 5},
		{ID: 412341251434, Price: 3974.89, Cost: 8934.24, Quantity: 41},
	},
	{
		{Year: 22, Month: 10, Day: 1, ID: 674398202423, Price: 234.23, Cost: 1324, Quantity: 1},
		{Year: 22, Month: 9, Day: 4, ID: 674398202423, Price: 100.50, Cost: 1324, Quantity: 1},
		{Year: 22, Month: 8, Day: 5, ID: 389432143927, Price: 3974.89, Cost: 8934.24, Quantity: 5},
		{Year: 22, Month: 7, Day: 6, ID: 674398202423, Price: 90109.22, Cost: 48.24, Quantity: 87},
		{Year: 22, Month: 6, Day: 4, ID: 402933466372, Price: 1324.89, Cost: 21432.24, Quantity: 4124},
		{Year: 22, Month: 6, Day: 4, ID: 198998421024, Price: 1094.89, Cost: 9021038.24, Quantity: 5},
		{Year: 22, Month: 6, Day: 7, ID: 412341251434, Price: 3974.89, Cost: 8934.24, Quantity: 41},
	},
	{},
}

func TestToUint40(t *testing.T) {
	value := TestDB[0][0].ID
	buf := make([]byte, 5)
	Database.ToUint40(buf, uint64(value))

	newVal := Database.FromUint40(buf)

	if value != newVal {
		t.Errorf("Values Don't match | Value: %v, New Value: %v", value, newVal)
	}
	t.Log(value)
	t.Log(newVal)
}

func TestLoadBackUp(t *testing.T) {

}

func TestReport(t *testing.T) {
	Database.Databases[1] = TestDB[1]
	Database.Expenses = []Database.Expense{
		{Year: 22, Month: 6, Day: 7, Frequency: 1, Amount: 1},
		{Year: 22, Month: 6, Day: 7, Frequency: 1, Amount: -43},
	}

	DayReport := Database.Report(0, []uint8{7, 6, 22})
	MonthReport := Database.Report(1, []uint8{7, 6, 22})
	YearReport := Database.Report(2, []uint8{7, 6, 22})

	// if DayReport != "Item Gain: 3974.889893,\n Item Loss: 8934.240234,\n Item Profit: -4959.350586,\n Expenses: 0.000000,\n	Gains: 0.000000,\n Report Total: -4959.350586" {
	// 	t.Error("Day's Report does not match up!")
	// }

	// if MonthReport != "Item Gain: 6394.669922,\n Item Loss: 9051404.000000,\n Item Profit: -9045009.000000,\n Expenses: 0.000000,\n Gains: 1.000000,\n Report Total: -9045008.000000" {
	// 	t.Error("Month's report does not match up!")
	// }

	// if YearReport != "Item Gain: 100813.507812,\n Item Loss: 9063035.000000,\n Item Profit: -8962221.000000,\n Expenses: 0.000000,\n Gains: 1.000000,\n Report Total: -8962220.000000" {
	// 	t.Error("Year's report does not match up!")
	// }

	//Find a way to actually compare the string result

	t.Log(DayReport)
	t.Log(MonthReport)
	t.Log(YearReport)
}

func TestCart(t *testing.T) {
	var red []Database.Sale // Create Test Cart

	t.Log(len(red))

	// Run functions on the cart
	red = Database.AddToCart(TestDB[0][0], red)

	if len(red) != 1 {
		t.Errorf("Cart not correct size, adding to cart is bugged | cartSize: %v", len(red))
	}

	if red[0] != TestDB[0][0] {
		t.Error("Item 0 in shopping cart does not match up with test item 0")
	}

	red = Database.AddToCart(TestDB[0][1], red)
	if len(red) != 2 {
		t.Errorf("Cart not correct size (addition) | cartSize: %v", len(red))
	}

	if red[1] != TestDB[0][1] {
		t.Error("Item 1 in shopping cart does not match up with test item 1")
	}

	t.Log(len(red))

	red = Database.DecreaseFromCart(TestDB[0][0], red)

	if len(red) != 1 {
		t.Errorf("Cart not correct size (subtraction) | cartSize: %v", len(red))
	}

	if red[0] != TestDB[0][1] {
		t.Error("Item 0 in shopping cart does not match up with test item 1, after shifting 1 to 0")
		t.Logf("Cart 0: %v, Test Items 1: %v", red[0], TestDB[1])
	}

	if result := Database.GetCartTotal(red); result != TestDB[0][1].Price {
		t.Errorf("Total does not match up.")
		t.Logf("Got: %v, Expected: %v", result, TestDB[0][1].Price)
	}

	red = Database.AddToCart(TestDB[0][1], red)

	if len(red) != 1 {
		t.Errorf("Cart not correct size (addition) | cartSize: %v", len(red))
	}

	if red[0].Quantity != 2 {
		t.Error("Cart Quantitiy does not match up!")
		t.Log(red[0].Quantity)
	}
}
