package Test

// import (
// 	"BronzeHermes/Database"
// 	"BronzeHermes/Debug"
// 	"testing"
// 	"time"
// )

// func TestCreatingSale(t *testing.T) {
// 	resetTestItemsAndSales()
// 	// have the data
// 	Database.Items = testItems
// 	Database.Customers = testCustomer

// 	errID := Database.CreateSale(1, "12", "20", "30", 3) // push to Database.Sales

// 	if errID != Debug.Success {
// 		t.Errorf("An error with ID: %d has occured", errID)
// 	}

// 	// check the len of Database.Sales

// 	if len(Database.Sales) != 1 {
// 		t.Errorf("Sale was not added or duplicated | have: %d, want: 1", len(Database.Sales))
// 	}

// 	// Check if it has the correct data
// 	for _, v := range Database.Sales {
// 		if Database.Items[v.ID].Name != "Viva" {
// 			t.Errorf("Incorrect Name | have: %s, want: 'Viva'", Database.Items[v.ID].Name)
// 		}

// 		if v.ID != 1 {
// 			t.Errorf("Incorrect ID | have: %d, want: 1", v.ID)
// 		}

// 		ti := time.Now()
// 		if v.Timestamp != ti.Unix() {
// 			t.Errorf("Timestamps Not Matching Up | have: %s \t want: %s ", time.Unix(v.Timestamp, 0).String(), ti.String())
// 		}

// 		if v.Price != 12 {
// 			t.Errorf("Incorrect Price | have; %f, want: 12.0", v.Price)
// 		}

// 		if v.Cost != 20 {
// 			t.Errorf("Incorrect Cost | have; %f, want: 20.0", v.Cost)
// 		}

// 		if v.Quantity != 30 {
// 			t.Errorf("Incorrect Quantity | have; %f, want: 30.0", v.Quantity)
// 		}

// 		if v.Customer != 3 {
// 			t.Errorf("Incorrect Customer | have; %d, want: 3", v.Customer)
// 		}

// 	}
// }

// func TestCreatingSaleFraction(t *testing.T) {
// 	resetTestItemsAndSales()
// 	Database.Items = testItems
// 	Database.Customers = testCustomer

// 	errID := Database.CreateSale(1, "12", "20", "30 7/8", 3) // push to Database.Sales

// 	if errID != Debug.Success {
// 		t.Errorf("An error with ID: %d has occured", errID)
// 	}

// 	// check the len of Database.Sales
// 	if len(Database.Sales) != 1 {
// 		t.Errorf("Sale was not added or duplicated | have: %d, want: 1", len(Database.Sales))
// 	}

// 	// Check if it has the correct data
// 	for _, v := range Database.Sales {
// 		if Database.Items[v.ID].Name != "Viva" {
// 			t.Errorf("Incorrect Name | have: %s, want: 'Viva'", Database.Items[v.ID].Name)
// 		}

// 		if v.ID != 1 {
// 			t.Errorf("Incorrect ID | have: %d, want: 1", v.ID)
// 		}

// 		if v.Year != 23 {
// 			t.Errorf("Incorrect Year | have: %d, want: 23", v.Year)
// 		}

// 		if v.Month != 06 {
// 			t.Errorf("Incorrect Month | have: %d, want: 06", v.Month)
// 		}

// 		if v.Day != 02 {
// 			t.Errorf("Incorrect Day | have: %d, want: 02", v.Day)
// 		}

// 		if v.Price != 12 {
// 			t.Errorf("Incorrect Price | have; %f, want: 12.0", v.Price)
// 		}

// 		if v.Cost != 20 {
// 			t.Errorf("Incorrect Cost | have; %f, want: 20.0", v.Cost)
// 		}

// 		if v.Quantity != 30.0+7.0/8.0 {
// 			t.Errorf("Incorrect Quantity | have; %f, want: %f", v.Quantity, 30.0+7.0/8.0)
// 		}

// 		if v.Customer != 3 {
// 			t.Errorf("Incorrect Customer | have; %d, want: 3", v.Customer)
// 		}

// 		if v.Usr != 0 {
// 			t.Errorf("Incorrect User | have; %d, want: 0", v.Usr)
// 		}

// 	}
// }

// func TestFailedCreatingSale(t *testing.T) {
// 	resetTestItemsAndSales()
// 	// have the data
// 	Database.Items = testItems
// 	Database.Customers = testCustomer

// 	errID := Database.CreateSale(1, "2023-06-02", "12", "20--23--40", "30", 3) // push to Database.Sales

// 	if errID != 0 {
// 		t.Errorf("An error with ID: %d has occured or slipped through| want: 0", errID)
// 	}

// 	// check the len of Database.Sales

// 	if len(Database.Sales) > 0 {
// 		t.Errorf("Sale Incorrectly added| have: %d, want: 0", len(Database.Sales))
// 	}

// }

// func TestFailedCreatingSaleDate(t *testing.T) {
// 	resetTestItemsAndSales()
// 	// have the data
// 	Database.Items = testItems
// 	Database.Customers = testCustomer

// 	errID := Database.CreateSale(1, "202//3-0.6-0,2.", "12", "20", "30", 3) // push to Database.Sales

// 	if errID != 0 {
// 		t.Errorf("An error with ID: %d has occured or slipped through| want: 0", errID)
// 	}

// 	// check the len of Database.Sales

// 	if len(Database.Sales) > 0 {
// 		t.Errorf("Sale Incorrectly added| have: %d, want: 0", len(Database.Sales))
// 	}

// }

// func TestFailedCreatingSaleFraction(t *testing.T) {
// 	resetTestItemsAndSales()
// 	// have the data
// 	Database.Items = testItems
// 	Database.Customers = testCustomer

// 	errID := Database.CreateSale(1, "12", "20", "30 //73--20340-/-8-342--", 3) // push to Database.Sales

// 	if errID != 0 {
// 		t.Errorf("An error with ID: %d has occured or slipped through| want: 0", errID)
// 	}

// 	// check the len of Database.Sales

// 	if len(Database.Sales) > 0 {
// 		t.Errorf("Sale Incorrectly added| have: %d, want: 0", len(Database.Sales))
// 	}

// }
