package Database

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetLine(selection string, dataType int) ([]string, [][]float32) {
	if selection == "" {
		return nil, nil
	}

	raw := strings.Split(selection, "-")

	if len(raw) < 2 {
		return nil, nil
	}

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

	for id, x := range Item {
		var totals []float32

		for i := uint8(1); i < 32; i++ {
			var total float32

			for v := len(Report) - 1; v >= 0; v-- {
				if Report[v].ID != id || Report[v].Day != i || Report[v].Month != date[1] || Report[v].Year != date[0] {
					continue
				}
				switch dataType {
				case 0:
					total += Report[v].Price
				case 1:
					total += float32(Report[v].Quantity)
				}
			}

			totals = append(totals, total)
		}

		if totals != nil {
			names = append(names, x.Name)
			sales = append(sales, totals)
		}
	}

	return names, sales
}

func GetPie(selection string, dataType int) ([]string, []float32) {
	if selection == "" {
		return nil, nil
	}

	raw := strings.Split(selection, "-")

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

	for id, val := range Item {
		var total float32

		for i := len(Report) - 1; i >= 0; i-- {
			if Report[i].ID != id || Report[i].Day != date[2] || Report[i].Month != date[1] || Report[i].Year != date[0] {
				continue
			}
			switch dataType {
			case 0:
				total += Report[i].Price
			case 1:
				total += float32(Report[i].Quantity)
			}
		}

		if total != 0 {
			names = append(names, val.Name)
			sales = append(sales, total)
		}
	}

	return names, sales
}

func CompileReport(selection uint8, date []uint8) string {
	// For Date: 0 = day, 1 = month, 2 = year, 3 = now

	if len(date) == 0 {
		y, month, day := time.Now().Date()
		year, _ := strconv.Atoi(strconv.Itoa(y)[1:])
		date = []uint8{uint8(day), uint8(month), uint8(year)}
	}

	var item_sales [3]float32
	var damages float32

	for _, v := range Report {
		if v.Year != date[YEARLY] {
			continue
		}

		if selection != YEARLY && v.Month != date[MONTHLY] {
			continue
		}

		if selection == ONCE && v.Day != date[ONCE] {
			continue
		}

		fmt.Println((v.Price-v.Cost)*v.Quantity, v.Price, v.Cost, Item[v.ID].Name)

		item_sales[0] += v.Price * v.Quantity
		item_sales[1] += v.Cost * v.Quantity
		item_sales[2] += (v.Price - v.Cost) * v.Quantity

		if v.Usr == 255 {
			damages += v.Cost * v.Quantity
		}
	}

	return fmt.Sprintf(
		"Item Revenue: %.2f,\nItem Cost: %.2f\nDamages: -%.2f,\nReport Total: %.2f",
		item_sales[0],  // Gain
		-item_sales[1], // Cost
		damages,        // Damages
		item_sales[2],  // Report Total
	)
}
