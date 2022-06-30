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

			for v := len(Databases[1]) - 1; i >= 0; i-- {
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
				total += Databases[1][i].Price
			case 1:
				total += Databases[1][i].Cost
			case 2:
				total += Databases[1][i].Price - Databases[1][i].Cost
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

func Report(selection uint8, date []uint8) string {
	// For Date: 0 = day, 1 = month, 2 = year, 3 = now

	if len(date) == 0 {
		day, month, y := time.Now().Date()
		year, _ := strconv.Atoi(strconv.Itoa(y)[1:])
		date = []uint8{uint8(day), uint8(month), uint8(year)}
	}

	var item_sales [3]float32

	for _, v := range Databases[1] {
		if v.Year != date[2] {
			continue
		}

		if selection != 2 && v.Month != date[1] {
			continue
		}

		if selection == ONCE && v.Day != date[ONCE] {
			continue
		}

		item_sales[0] += v.Price * float32(v.Quantity)
		item_sales[1] += v.Cost * float32(v.Quantity)
		item_sales[2] += (v.Price - v.Cost) * float32(v.Quantity)
	}

	var expenses float32
	var gifts float32

	for i := len(Expenses) - 1; i >= 0; i-- {
		if Expenses[i].Date[YEARLY] != date[YEARLY] {
			continue
		}

		if selection != 2 && Expenses[i].Date[MONTHLY] != date[MONTHLY] {
			continue
		}

		if selection == ONCE && Expenses[i].Date[ONCE] != date[ONCE] {
			continue
		}

		if Expenses[i].Amount < 0 {
			expenses += Expenses[i].Amount
		} else {
			gifts += Expenses[i].Amount
		}
	}

	return fmt.Sprintf(
		"Item Gain: %.2f,\n Item Loss: %.2f,\n Item Profit: %.2f,\n Expenses: %.2f,\n Gains: %.2f,\n Report Total: %.2f",
		item_sales[0],                // Sold Amount
		item_sales[1],                // Cost
		item_sales[2],                // Profit
		expenses,                     // Expenses
		gifts,                        // Gifts
		item_sales[2]+expenses+gifts, // Report Total
	)
}

func FindItem(ID int) Sale {
	for i, z := 0, len(Databases[0]); i < len(Databases[0]); i++ {
		if int(Databases[1][i].ID) == ID {
			return Databases[1][i]
		}

		z -= 1
		if int(Databases[1][z].ID) == ID {
			return Databases[1][z]
		}
	}
	return Sale{}
}
