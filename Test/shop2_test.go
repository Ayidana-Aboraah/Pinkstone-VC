package Test

import (
	"BronzeHermes/Database"
	"testing"
)

func TestAddingItemNormal(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("", "", "", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		// Check that teh proper transformation
		if s.Quantity != 1 {
			t.Errorf("Illegal Quantity Modification, want: 1, have: %f", s.Quantity)
		}

		if s.Price != 5 {
			t.Errorf("Illegal Price Modificaiton, want: 5, have: %f", s.Price)
		}

		if s.Cost != 5 {
			t.Errorf("Illegal Cost Modificaiton, want: 5, have: %f", s.Cost)
		}
	}
}

func TestAddingItemWithBargin(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("10", "", "", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		// Check that teh proper transformation
		if s.Quantity != 1 {
			t.Errorf("Illegal Quantity Modification, want: 1, have: %f", s.Quantity)
		}

		if s.Price != 10 {
			t.Errorf("Illegal Price Modificaiton, want: 10, have: %f", s.Price)
		}

		if s.Cost != 5 {
			t.Errorf("Illegal Cost Modificaiton, want: 5, have: %f", s.Cost)
		}
	}
}

func TestAddingItemInPieces(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("", "1", "12", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		// Check that teh proper transformation
		if s.Quantity != 1.0/12.0 {
			t.Errorf("Illegal Quantity Modification, want: 1/12, have: %f", s.Quantity)
		}

		if s.Price != 5 {
			t.Errorf("Illegal Price Modificaiton, want: 5, have: %f", s.Price)
		}

		if s.Cost != 5 {
			t.Errorf("Illegal Cost Modificaiton, want: 5, have: %f", s.Cost)
		}
	}
}

func TestAddingItemWithQuantityInPieces(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("25", "1", "12", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		// Check that teh proper transformation
		if s.Quantity != 1.0/12.0 {
			t.Errorf("Illegal Quantity Modification, want: 1/12, have: %f", s.Quantity)
		}

		if s.Price != 25.0*12.0 {
			t.Errorf("Illegal Price Modificaiton, want: 5, have: %f", s.Price)
		}

		if s.Cost != 5 {
			t.Errorf("Illegal Cost Modificaiton, want: 5, have: %f", s.Cost)
		}
	}
}

func TestInvalidBargin(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("-122-123-2", "", "", &s)
	switch err {
	case 0:
	case 1:
		t.Error("No input sent to the piece || total")
	case -1:
		t.Log(s)
		t.Error("This Data is invalid and should not pass")
	}
}

func TestMissingPiece(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("", "", "12", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
	case -1:
		t.Log(s)
		t.Error("This Data is invalid and should not pass")
	}
}

func TestMissingTotal(t *testing.T) {
	s := Database.Sale{Price: 5, Cost: 5, Quantity: 1}
	err := Database.ProcessNewItemData("", "1", "", &s)
	switch err {
	case 0:
		t.Error("Invalid Input passed into ProcessNewItemData")
	case 1:
	case -1:
		t.Log(s)
		t.Error("This Data is invalid and should not pass")
	}
}

func TestCartTotal(t *testing.T) {
	total := Database.GetCartTotal(testSales)
	if total != 30 {
		t.Log(total)
		t.Error("Some Sort of addition error in GetCartTotal")
	}
}
