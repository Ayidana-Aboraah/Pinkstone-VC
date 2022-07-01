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

			for v := len(Reports[0]) - 1; i >= 0; i-- {
				if Reports[0][v].ID != id || Reports[0][v].Day != i || Reports[0][v].Month != date[1] || Reports[0][v].Year != date[0] {
					continue
				}
				switch dataType {
				case 0:
					total += Reports[0][v].Price
				case 1:
					total += Reports[0][v].Cost
				case 2:
					total += Reports[0][v].Price - Reports[0][v].Cost
				case 3:
					total += float32(Reports[0][v].Quantity)
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

		for i := len(Reports[0]) - 1; i >= 0; i-- {
			if Reports[0][i].ID != id || Reports[0][i].Day != date[2] || Reports[0][i].Month != date[1] || Reports[0][i].Year != date[0] {
				continue
			}
			switch dataType {
			case 0:
				total += Reports[0][i].Price
			case 1:
				total += Reports[0][i].Cost
			case 2:
				total += Reports[0][i].Price - Reports[0][i].Cost
			case 3:
				total += float32(Reports[0][i].Quantity)
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

	for _, v := range Reports[0] {
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
	for i, z := 0, len(Reports[0]); i < len(Reports[0]); i++ {
		if int(Reports[0][i].ID) == ID {
			return Reports[0][i]
		}

		z -= 1
		if int(Reports[0][z].ID) == ID {
			return Reports[0][z]
		}
	}
	return Sale{}
}
