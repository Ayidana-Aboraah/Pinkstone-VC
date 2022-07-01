package test

import (
	"BronzeHermes/Database"
	"testing"

	"fyne.io/fyne/v2/app"
)

func TestFileSave(t *testing.T) {
	a := app.NewWithID("Testing")

	Database.NameKeys = TestNames
	Database.Reports = TestDB
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

	for k, v := range TestNames {
		if _, found := Database.NameKeys[k]; !found {
			t.Errorf("Key %v !found, value is %v", k, v)
		}
	}

	for i := range TestExpenses {
		if TestExpenses[i] != Database.Expenses[i] {
			t.Errorf("Expenses do not match at %v", i)
		}
	}

	t.Logf("Now just for vals:\n Expenses: %v \n Name Keys: %v\n Databaseses: %v", Database.Expenses, Database.NameKeys, Database.Reports)

	Database.DataInit(true)
	a.Quit()
}

func TestSaveBackUp(t *testing.T) {
	a := app.NewWithID("Testing")

	Database.Reports = TestDB
	Database.NameKeys = TestNames
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

	// for i := range TestDB {
	// 	if len(TestDB[i]) != len(Database.Databases[i]) {
	// 		t.Errorf("Test DB and Datbase %v aren't the same\n", i)
	// 		t.Logf("Test DB: %v, Datbase: %v\n", len(TestDB[i]), len(Database.Databases[i]))
	// 		// t.Logf("Test DB: %v\n DB: %v", TestDB[i], Database.Databases[i])
	// 	}
	// 	for x := range TestDB[i] {
	// 		if Database.Databases[i][x] != TestDB[i][x] {
	// 			t.Errorf("Database %v, entry %v don't match up", i, x)
	// 		}
	// 	}
	// }

	for k, v := range TestNames {
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
