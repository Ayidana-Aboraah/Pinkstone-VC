package Test

import (
	"BronzeHermes/Database"
	"testing"
)

var testCart2 = []Database.Sale{ // Use 6 since that's the one that has proper data
	{ID: 6, Price: 12, Cost: 5, Quantity: 2, Year: 23, Month: 02, Day: 01},
	{ID: 6, Price: 13, Cost: 6, Quantity: 3, Year: 23, Month: 02, Day: 01},
	{ID: 6, Price: 14, Cost: 7, Quantity: 6, Year: 23, Month: 02, Day: 02},
	{ID: 6, Price: 15, Cost: 8, Quantity: 7, Year: 23, Month: 03, Day: 02},
	{ID: 6, Price: 16, Cost: 9, Quantity: 12, Year: 24, Month: 03, Day: 01, Usr: 255},
	{ID: 6, Price: 16, Cost: 9, Quantity: 14, Year: 24, Month: 01, Day: 01},
}

func TestReportDay(t *testing.T) {
	answers := [4]float32{
		63,
		-28,
		0,
		35,
	}

	Database.Sales = testCart2
	_, vals := Database.CompileReport(0, []uint8{01, 02, 23})
	for i := range vals {
		if vals[i] != answers[i] {
			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
		}
	}
}

func TestReportMonth(t *testing.T) {
	answers := [4]float32{
		147,
		-70,
		0,
		77,
	}

	Database.Sales = testCart2
	_, vals := Database.CompileReport(1, []uint8{01, 02, 23})
	for i := range vals {
		if vals[i] != answers[i] {
			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
		}
	}
}

func TestReportYear(t *testing.T) {
	answers := [4]float32{
		252,
		-126,
		0,
		126,
	}

	Database.Sales = testCart2
	_, vals := Database.CompileReport(2, []uint8{01, 02, 23})
	for i := range vals {
		if vals[i] != answers[i] {
			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
		}
	}
}

func TestReportYearWithDamages(t *testing.T) {
	answers := [4]float32{
		416,
		-234,
		-108,
		182,
	}

	Database.Sales = testCart2
	_, vals := Database.CompileReport(2, []uint8{01, 02, 24})
	for i := range vals {
		if vals[i] != answers[i] {
			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
		}
	}
}
