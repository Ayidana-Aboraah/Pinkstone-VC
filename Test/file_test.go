package Test

import (
	"BronzeHermes/Database"
	"testing"
)

var testUserNames = []string{
	"Bob",
	"Penny",
	"Poppy",
	"123481023984012983",
}

func TestSaveUsers(t *testing.T) {
	Database.Users = testUserNames
	Database.SaveNLoadUsers()
	if len(Database.Users) != len(testUserNames) {
		t.Errorf("Unequal lengths after saving")
	}

	for i := range testUserNames {
		if testUserNames[i] != Database.Users[i] {
			t.Errorf("%s != %s | idx: %d of testUsers, error after lodaing data", testUserNames[i], Database.Users[i], i)
		}
	}
}

func TestSaveNoUser(t *testing.T) {
	Database.Users = []string{}
	Database.SaveNLoadUsers()
	if len(Database.Users) != 0 {
		t.Error("Users Databasse lenght is != 0, when handed an empty array")
	}
}

func TestSaveBlankUser(t *testing.T) {
	Database.Users = []string{""}
	Database.SaveNLoadUsers()
	if len(Database.Users) != 1 {
		t.Error("Users Databasse lenght is != 1, when handed an empty string")
	}
}
