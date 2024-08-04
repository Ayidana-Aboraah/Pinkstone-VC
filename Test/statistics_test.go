package Test

import (
	"BronzeHermes/Database"
	"time"
)

var timestamp = time.Now().Unix()

var testCart2 = []Database.Sale{ // Use 6 since that's the one that has proper data
	{ID: 6, Price: 12, Cost: 5, Quantity: 2, Timestamp: timestamp, Customer: 1},
	{ID: 6, Price: 13, Cost: 6, Quantity: 3, Timestamp: timestamp, Customer: 1},
	{ID: 6, Price: 14, Cost: 7, Quantity: 6, Timestamp: timestamp, Customer: 0},
	{ID: 6, Price: 15, Cost: 8, Quantity: 7, Timestamp: timestamp, Customer: 1},
	{ID: 6, Price: 16, Cost: 9, Quantity: 12, Timestamp: timestamp, Customer: 0},
	{ID: 6, Price: 16, Cost: 9, Quantity: 14, Timestamp: timestamp, Customer: 1},
}

// func TestReportDay(t *testing.T) {
// 	Database.Sales = testCart2
// 	answers := [4]float32{
// 		63,
// 		-28,
// 		0,
// 		35,
// 	}

// 	Database.Sales = testCart2
// 	_, vals := Database.CompileReport(0, []int{})
// 	for i := range vals {
// 		if vals[i] != answers[i] {
// 			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
// 		}
// 	}
// }

// func TestReportMonth(t *testing.T) {
// 	answers := [4]float32{
// 		147,
// 		-70,
// 		0,
// 		77,
// 	}

// 	Database.Sales = testCart2
// 	_, vals := Database.CompileReport(1, []int{})
// 	for i := range vals {
// 		if vals[i] != answers[i] {
// 			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
// 		}
// 	}
// }

// func TestReportYear(t *testing.T) {
// 	answers := [4]float32{
// 		252,
// 		-126,
// 		0,
// 		126,
// 	}

// 	Database.Sales = testCart2
// 	_, vals := Database.CompileReport(2, []int{})
// 	for i := range vals {
// 		if vals[i] != answers[i] {
// 			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
// 		}
// 	}
// }

// func TestReportYearWithDamages(t *testing.T) {
// 	answers := [4]float32{
// 		416,
// 		-234,
// 		-108,
// 		182,
// 	}

// 	Database.Sales = testCart2
// 	_, vals := Database.CompileReport(2, []int{})
// 	for i := range vals {
// 		if vals[i] != answers[i] {
// 			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
// 		}
// 	}
// }

// func TestReportSpecifiedYear(t *testing.T) {
// 	answers := [4]float32{
// 		416,
// 		-234,
// 		-108,
// 		182,
// 	}

// 	Database.Sales = testCart2
// 	y, m, d := time.Now().Date()
// 	_, vals := Database.CompileReport(2, []int{d, int(m), y})
// 	for i := range vals {
// 		if vals[i] != answers[i] {
// 			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
// 		}
// 	}
// }

// func TestReportSpecifiedMonth(t *testing.T) {
// 	answers := [4]float32{
// 		416,
// 		-234,
// 		-108,
// 		182,
// 	}

// 	Database.Sales = testCart2
// 	y, m, d := time.Now().Date()
// 	_, vals := Database.CompileReport(1, []int{d, int(m), y})
// 	for i := range vals {
// 		if vals[i] != answers[i] {
// 			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
// 		}
// 	}
// }

// func TestReportSpecifiedDay(t *testing.T) {
// 	answers := [4]float32{
// 		416,
// 		-234,
// 		-108,
// 		182,
// 	}

// 	Database.Sales = testCart2
// 	y, m, d := time.Now().Date()
// 	_, vals := Database.CompileReport(0, []int{d, int(m), y})
// 	for i := range vals {
// 		if vals[i] != answers[i] {
// 			t.Errorf("Report Result unexpected | have: %f, want: %f", vals[i], answers[i])
// 		}
// 	}
// }
