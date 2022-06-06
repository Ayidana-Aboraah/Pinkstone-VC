package test

import (
	"BronzeHermes/Database"
	"testing"
)

var testItems = []Database.Sale{
	{ID: 674398202423, Price: 234.23, Cost: 1324, Quantity: 1},
	{ID: 674398202423, Price: 100.50, Cost: 1324, Quantity: 1},
	{ID: 389432143927, Price: 3974.89, Cost: 8934.24, Quantity: 5},
	{ID: 674398202423, Price: 90109.22, Cost: 48.24, Quantity: 87},
	{ID: 402933466372, Price: 1324.89, Cost: 21432.24, Quantity: 4124},
	{ID: 198998421024, Price: 1094.89, Cost: 9021038.24, Quantity: 5},
	{ID: 412341251434, Price: 3974.89, Cost: 8934.24, Quantity: 41},
}

func TestToUint40(t *testing.T) {
	value := testItems[0].ID
	buf := make([]byte, 5)
	Database.ToUint40(buf, uint64(value))

	newVal := Database.FromUint40(buf)

	if value != newVal {
		t.Errorf("Values Don't match | Value: %v, New Value: %v", value, newVal)
	}
	t.Log(value)
	t.Log(newVal)
}

func TestLoadBackUp(t *testing.T){
	
}

func TestCart(t *testing.T) {
	var red []Database.Sale // Create Test Cart

	t.Log(len(red))

	// Run functions on the cart
	red = Database.AddToCart(testItems[0], red)

	if len(red) != 1 {
		t.Errorf("Cart not correct size, adding to cart is bugged | cartSize: %v", len(red))
	}

	if red[0] != testItems[0] {
		t.Error("Item 0 in shopping cart does not match up with test item 0")
	}

	red = Database.AddToCart(testItems[1], red)
	if len(red) != 2 {
		t.Errorf("Cart not correct size (addition) | cartSize: %v", len(red))
	}

	if red[1] != testItems[1] {
		t.Error("Item 1 in shopping cart does not match up with test item 1")
	}

	t.Log(len(red))

	red = Database.DecreaseFromCart(testItems[0], red)

	if len(red) != 1 {
		t.Errorf("Cart not correct size (subtraction) | cartSize: %v", len(red))
	}

	if red[0] != testItems[1] {
		t.Error("Item 0 in shopping cart does not match up with test item 1, after shifting 1 to 0")
		t.Logf("Cart 0: %v, Test Items 1: %v", red[0], testItems[1])
	}

	if result := Database.GetCartTotal(red); result != testItems[1].Price {
		t.Errorf("Total does not match up.")
		t.Logf("Got: %v, Expected: %v", result, testItems[1].Price)
	}

	red = Database.AddToCart(testItems[1], red)

	if len(red) != 1 {
		t.Errorf("Cart not correct size (addition) | cartSize: %v", len(red))
	}

	if red[0].Quantity != 2 {
		t.Error("Cart Quantitiy does not match up!")
		t.Log(red[0].Quantity)
	}
}
