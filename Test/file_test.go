package Test

// import (
// 	"BronzeHermes/Database"
// 	"testing"

// 	"fyne.io/fyne/v2/app"
// )

// func TestFileSave(t *testing.T) {
// 	a := app.NewWithID("Testing")

// 	Database.Reports = TestDB
// 	Database.ItemKeys = TestItemKeys
// 	Database.Expenses = TestExpenses

// 	Database.DataInit(false)
// 	err := Database.SaveData()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	Database.Expenses = []Database.Expense{}

// 	err = Database.LoadData()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	SaveAndLoadTest(t)

// 	Database.DataInit(true)
// 	a.Quit()
// }

// func TestSaveBackUp(t *testing.T) {
// 	a := app.NewWithID("Testing")

// 	Database.Reports = TestDB
// 	Database.ItemKeys = TestItemKeys
// 	Database.Expenses = TestExpenses

// 	Database.DataInit(false)
// 	err := Database.SaveBackUp()
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	Database.Expenses = []Database.Expense{}

// 	err = Database.LoadBackUp()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	SaveAndLoadTest(t)

// 	// Database.DataInit(true)
// 	a.Quit()
// }

// func SaveAndLoadTest(t *testing.T) {
// 	for i := range TestDB {
// 		if len(TestDB[i]) != len(Database.Reports[i]) {
// 			t.Logf("Test DB and Datbase %v aren't the same", i)
// 		}
// 		for x := range TestDB[i] {
// 			if Database.Reports[i][x] != TestDB[i][x] {
// 				t.Errorf("Database %v, entry %v don't match up", i, x)
// 			}
// 		}
// 	}

// 	for k, v := range TestItemKeys {
// 		if val, found := Database.ItemKeys[k]; !found {
// 			t.Errorf("Key %v !found, value is %v", k, v)
// 		} else {
// 			if val != v {
// 				t.Errorf("Key %v, values do not match up", k)
// 				t.Errorf("Test: %v, DB Keys: %v", v, val)
// 			}
// 		}

// 	}

// 	for i := range TestExpenses {
// 		if TestExpenses[i] != Database.Expenses[i] {
// 			t.Errorf("Expenses do not match at %v", i)
// 		}
// 	}

// 	t.Logf("%v\n", Database.Expenses)
// 	t.Logf("%v\n", Database.Reports)
// 	for k, v := range Database.ItemKeys {
// 		t.Logf("%v : %v\n", k, v)
// 	}
// }
