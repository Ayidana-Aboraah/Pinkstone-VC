package Test

import (
	"BronzeHermes/Database"
	"fmt"
	"strings"
	"testing"
)

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
	Database.Reports[0] = TestDB[0]
	Database.Expenses = TestExpenses

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

	t.Log(YearReport)
}
