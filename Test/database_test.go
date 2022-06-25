package test

import (
	"BronzeHermes/Database"
	"fmt"
	"strings"
	"testing"
)

var TestNames = map[uint64]string{
	999999999999: "sammy",
	674398202423: "Cark",
	389432143927: "Banker",
	402933466372: "Blackeg",
	198998421024: "Boomb",
	412341251434: "Sus",
}

var TestDB = [3][]Database.Sale{
	{
		{ID: 999999999999, Price: 234.23, Cost: 1324, Quantity: 1},
		{ID: 674398202423, Price: 100.50, Cost: 1324, Quantity: 1},
		{ID: 389432143927, Price: 3974.89, Cost: 8934.24, Quantity: 5},
		{ID: 674398202423, Price: 90109.22, Cost: 48.24, Quantity: 87},
		{ID: 402933466372, Price: 1324.89, Cost: 21432.24, Quantity: 4124},
		{ID: 198998421024, Price: 1094.89, Cost: 9021038.24, Quantity: 5},
		{ID: 412341251434, Price: 3974.89, Cost: 8934.24, Quantity: 41},
	},
	{
		{Year: 22, Month: 10, Day: 1, ID: 674398202423, Price: 111.23, Cost: 1324, Quantity: 1},
		{Year: 22, Month: 9, Day: 4, ID: 674398202423, Price: 100.50, Cost: 555, Quantity: 1},
		{Year: 22, Month: 8, Day: 5, ID: 389432143927, Price: 222.89, Cost: 332.24, Quantity: 5},
		{Year: 22, Month: 7, Day: 6, ID: 674398202423, Price: 444.22, Cost: 222.24, Quantity: 7},
		{Year: 22, Month: 6, Day: 4, ID: 402933466372, Price: 333.21, Cost: 232.24, Quantity: 4},
		{Year: 22, Month: 6, Day: 4, ID: 198998421024, Price: 555.22, Cost: 938.24, Quantity: 5},
		{Year: 22, Month: 6, Day: 7, ID: 412341251434, Price: 666.22, Cost: 834.24, Quantity: 1},
	},
	{},
}

var testExpenses = []Database.Expense{
	{Date: [3]uint8{7, 6, 22}, Frequency: 1, Amount: 1, Name: "red"},
	{Date: [3]uint8{7, 6, 22}, Frequency: 1, Amount: -43},

	{Date: [3]uint8{8, 6, 22}, Frequency: 0, Amount: -43},
	{Date: [3]uint8{7, 7, 22}, Frequency: 1, Amount: 3},
	{Date: [3]uint8{7, 6, 23}, Frequency: 1, Amount: -13},
}

func TestToUint40(t *testing.T) {
	value := TestDB[0][0].ID
	buf := make([]byte, 5)
	Database.PutUint40(buf, uint64(value))

	newVal := Database.FromUint40(buf)

	if value != newVal {
		t.Errorf("Values Don't match | Value: %v, New Value: %v", value, newVal)
	}
	t.Log(value)
	t.Log(newVal)
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

func TestReport(t *testing.T) {
	Database.Databases[1] = TestDB[1]
	Database.Expenses = testExpenses

	test_report_outputs := func() []string {
		base := "Item Gain: %.2f,\n Item Loss: %.2f,\n Item Profit: %.2f,\n Expenses: %.2f,\n Gains: %.2f,\n Report Total: %.2f"
		return []string{
			fmt.Sprintf(base, 666.22, 834.24, -168.02, -43.0, 1.0, -210.02),    //Day Report
			fmt.Sprintf(base, 4775.16, 6454.4, -1679.24, -86.0, 1.0, -1764.24), //Month Report
			fmt.Sprintf(base, 9210.88, 11550.28, -2339.4, -86.0, 4.0, -2421.4), //Year Report
		}
	}()

	DayReport := Database.Report(0, []uint8{7, 6, 22})
	MonthReport := Database.Report(1, []uint8{7, 6, 22})
	YearReport := Database.Report(2, []uint8{7, 6, 22})

	if strings.Compare(DayReport, test_report_outputs[0]) != 0 {
		t.Error("Day's Report does not match up!")
		t.Log(test_report_outputs[0])
		t.Log(DayReport)
	}
	strings.Compare(DayReport, test_report_outputs[0])

	if strings.Compare(MonthReport, test_report_outputs[1]) != 0 {
		t.Error("Month's report does not match up!")
		t.Log(test_report_outputs[1])
		t.Log(MonthReport)
	}

	if strings.Compare(YearReport, test_report_outputs[2]) != 0 {
		t.Error("Year's report does not match up!")
		t.Log(test_report_outputs[2])
		t.Log(YearReport)
	}

}
