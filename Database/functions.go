package Database

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetLine(selection string, dataType int, database []Sale) ([]string, [][]float32) {
	if selection == "" {
		return nil, nil
	}

	raw := strings.Split(selection, "/")

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil, nil
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, nil
	}

	date := []uint8{
		uint8(year),
		uint8(month),
	}

	//Change the error handling for this to show that you can't convert
	var sales [][]float32
	var names []string

	for id, name := range NameKeys {
		var totals []float32

		for i := uint8(1); i < 32; i++ {
			var total float32

			for v := len(Databases[REPORT]) - 1; i >= 0; i-- {
				if Databases[1][v].ID != id || Databases[1][v].Day != i || Databases[1][v].Month != date[1] || Databases[1][v].Year != date[0] {
					continue
				}
				switch dataType {
				case 0:
					total += Databases[1][v].Price
				case 1:
					total += Databases[1][v].Cost
				case 2:
					total += Databases[1][v].Price - Databases[1][v].Cost
				case 3:
					total += float32(Databases[1][v].Quantity)
				}
			}

			totals = append(totals, total)
		}

		if totals != nil {
			names = append(names, name)
			sales = append(sales, totals)
		}
	}

	return names, sales
}

func GetPie(selection string, dataType int) ([]string, []float32) {
	if selection == "" {
		return nil, nil
	}

	raw := strings.Split(selection, "/")

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil, nil
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, nil
	}

	day, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, nil
	}

	date := []uint8{
		uint8(year),
		uint8(month),
		uint8(day),
	}

	var sales []float32
	var names []string

	for id, name := range NameKeys {
		var total float32

		for i := len(Databases[1]) - 1; i >= 0; i-- {
			if Databases[1][i].ID != id || Databases[1][i].Day != date[2] || Databases[1][i].Month != date[1] || Databases[1][i].Year != date[0] {
				continue
			}
			switch dataType {
			case 0:
				total += Databases[REPORT][i].Price
			case 1:
				total += Databases[REPORT][i].Cost
			case 2:
				total += Databases[REPORT][i].Price - Databases[1][i].Cost
			case 3:
				total += float32(Databases[1][i].Quantity)
			}
		}

		if total != 0 {
			names = append(names, name)
			sales = append(sales, total)
		}
	}

	return names, sales
}

func Report(selection uint8) string {
	now := time.Now()
	year, _ := strconv.Atoi(strconv.Itoa(now.Year())[1:])

	var item_sales [3]float32
	for _, v := range Databases[REPORT] {
		if v.Year != uint8(year) && v.Month != uint8(now.Month()) {
			continue
		}

		if selection == 0 && v.Day != uint8(now.Day()) {
			continue
		}

		item_sales[0] += v.Price
		item_sales[1] += v.Cost
		item_sales[2] += v.Price - v.Cost
	}

	//Add up expenses
	//Add up patreonage
	//Combine item and expenses to get Report Sales

	var expenses float32
	var gifts float32
	for i := len(Expenses) - 1; i >= 0; i-- {
		if Expenses[i].Year != uint8(year) {
			continue
		}

		if selection == Expenses[i].Frequency {
			if Expenses[i].Frequency == ONCE && Expenses[i].Day != uint8(now.Day()) && Expenses[i].Month != uint8(now.Month()) {
				continue
			}

			if Expenses[i].Amount < 0 {
				expenses += Expenses[i].Amount
			} else {
				gifts += Expenses[i].Amount
			}
		}
	}

	return fmt.Sprintf(
		"Item Gain: %f,\n Item Loss: %f,\n Item Profit: %f,\n Expenses: %f,\n Gains: %f, Report Total: %f",
		item_sales[0],                // Sold Amount
		item_sales[1],                // Cost
		item_sales[2],                // Profit
		expenses,                     // Expenses
		gifts,                        // Gifts
		item_sales[2]+expenses+gifts, // Report Total
	)
}

func CustomDateReport(selection string) string {
	raw := strings.Split(selection, "/")

	var variant uint8

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return ""
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		variant = 1
	}

	day, err := strconv.Atoi(raw[1])
	if err != nil {
		variant = 2
	}

	var item_sales [3]float32
	for _, v := range Databases[1] {
		if v.Year != uint8(year) {
			continue
		}

		if variant < 2 && v.Month != uint8(month) {
			continue
		}
		if variant < 1 && v.Day != uint8(day) {
			continue
		}

		item_sales[0] += v.Price
		item_sales[1] += v.Cost
		item_sales[2] += v.Price - v.Cost
	}

	return fmt.Sprintf(
		"Item Gain: %f,\n Item Loss: %f,\n Item Profit: %f,\n Expenses: %f,\n Gains: %f, Report Total: %f",
		item_sales[0],
		item_sales[1],
		item_sales[2],
		0.0, // Expenses
		0.0, // Gifts
		0.0, // Report Total
	)
}

func FindItem(ID int) Sale {
	for i, z := 0, len(Databases[0]); i < len(Databases[0]); i++ {
		if int(Databases[ITEMS][i].ID) == ID {
			return Databases[ITEMS][i]
		}

		z -= 1
		if int(Databases[ITEMS][z].ID) == ID {
			return Databases[ITEMS][z]
		}
	}
	return Sale{}
}
