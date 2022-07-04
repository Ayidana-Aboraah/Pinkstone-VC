package test

import (
	"BronzeHermes/Database"
	"testing"

	"fyne.io/fyne/v2/app"
)

func TestFileSave(t *testing.T) {
	a := app.NewWithID("Testing")

	Database.Items = TestItems
	Database.Reports = TestDB
	Database.ItemKeys = TestItemKeys
	Database.Expenses = TestExpenses

	Database.DataInit(false)
	err := Database.SaveData()
	if err != nil {
		t.Error(err)
	}

	Database.Expenses = []Database.Expense{}

	err = Database.LoadData()
	if err != nil {
		t.Error(err)
	}

	if len(TestItems) != len(Database.Items) {
		t.Errorf("Lengths of test and normal items don't match up: Test: %v, Items: %v", len(TestItems), len(Database.Items))
	}

	for i := range TestItems {
		if TestItems[i] != Database.Items[i] {
			t.Errorf("DB Item %d does not match up with Test", i)
			t.Errorf("DB Item: %v\n Test Item: %v", Database.Items[i], TestItems[i])
		}
	}

	for i := range TestDB {
		if len(TestDB[i]) != len(Database.Reports[i]) {
			t.Logf("Test DB and Datbase %v aren't the same", i)
		}
		for x := range TestDB[i] {
			if Database.Reports[i][x] != TestDB[i][x] {
				t.Errorf("Database %v, entry %v don't match up", i, x)
			}
		}
	}

	for k, v := range TestItemKeys {
		if val, found := Database.ItemKeys[k]; !found {
			t.Errorf("Key %v !found, value is %v", k, v)
		} else {
			if val != v {
				t.Errorf("Key %v, values do not match up", k)
				t.Errorf("Test: %v, DB Keys: %v", v, val)
			}
		}

	}

	for i := range TestExpenses {
		if TestExpenses[i] != Database.Expenses[i] {
			t.Errorf("Expenses do not match at %v", i)
		}
	}

	// t.Logf("Now just for vals:\n Expenses: %v \n Name Keys: %v\n Databaseses: %v", Database.Expenses, Database.NameKeys, Database.Reports)

	// Database.DataInit(true)
	a.Quit()
}

func TestSaveBackUp(t *testing.T) {
	a := app.NewWithID("Testing")

	Database.Reports = TestDB
	// Database.NameKeys = TestNames
	Database.ItemKeys = TestItemKeys
	Database.Expenses = TestExpenses

	Database.DataInit(false)
	err := Database.BackUpAllData()
	if err != nil {
		t.Error(err)
	}

	Database.Expenses = []Database.Expense{}

	err = Database.LoadBackUp()
	if err != nil {
		t.Error(err)
	}

	for i := range TestDB {
		if len(TestDB[i]) != len(Database.Reports[i]) {
			t.Errorf("Test DB and Datbase %v aren't the same\n", i)
			t.Logf("Test DB: %v, Datbase: %v\n", len(TestDB[i]), len(Database.Reports[i]))
			// t.Logf("Test DB: %v\n DB: %v", TestDB[i], Database.Databases[i])
		}
		for x := range TestDB[i] {
			if Database.Reports[i][x] != TestDB[i][x] {
				t.Errorf("Database %v, entry %v don't match up", i, x)
			}
		}
	}

	for k, v := range TestItemKeys {
		if _, found := Database.NameKeys[k]; !found {
			t.Errorf("Key %v !found, value is %v", k, v)
		}
	}

	if len(TestExpenses) != len(Database.Expenses) {
		t.Errorf("The lengths of the expenses do not match up\n")
		t.Logf("Test: %v, Expenses: %v", len(TestExpenses), len(Database.Expenses))
	}

	for i := range TestExpenses {
		// t.Logf("Idx: %v\n Te s: %v\n Exp: %v\n", i, testExpenses[i], Database.Expenses[i])
		if TestExpenses[i] != Database.Expenses[i] {
			t.Errorf("Expenses do not match at %v", i)
		}
	}

	// t.Logf("Now just for vals:\n Expenses: %v \n Name Keys: %v\n Databaseses: %v", Database.Expenses, Database.NameKeys, Database.Databases)

	Database.DataInit(true)
	a.Quit()
}
