package Database

import (
	"BronzeHermes/Debug"
	unknown "BronzeHermes/Unknown"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func CompileReport(selection uint8, date []uint8) (string, [4]float32) {
	// For Date: 0 = day, 1 = month, 2 = year, 3 = now

	if len(date) == 0 {
		y, month, day := time.Now().Date()
		year, _ := strconv.Atoi(strconv.Itoa(y)[1:])
		date = []uint8{uint8(day), uint8(month), uint8(year)}
	}

	var item_sales [3]float32
	var damages float32

	for _, v := range Sales {
		if v.Year != date[YEARLY] {
			continue
		}

		if selection != YEARLY && v.Month != date[MONTHLY] {
			continue
		}

		if selection == ONCE && v.Day != date[ONCE] {
			continue
		}

		item_sales[0] += v.Price * v.Quantity
		item_sales[1] += v.Cost * v.Quantity
		item_sales[2] += (v.Price - v.Cost) * v.Quantity

		if v.Usr == 255 {
			damages += v.Cost * v.Quantity
		}
	}

	return fmt.Sprintf(
		"Item Revenue: %.2f,\nItem Cost: %.2f\nDamages: %.2f,\nReport Total: %.2f",
		item_sales[0],  // Gain
		-item_sales[1], // Cost
		-damages,       // Damages
		item_sales[2],  // Report Total
	), [4]float32{item_sales[0], -item_sales[1], -damages, item_sales[2]}
}

func RemoveFromSales(index int) {

	if Items[Sales[index].ID].Quantity[0] <= 0 {
		Items[Sales[index].ID].Cost[0] = Sales[index].Cost
		Items[Sales[index].ID].Quantity[0] = Sales[index].Quantity
	} else {
		for i, v := range Items[Sales[index].ID].Cost {
			if v == 0 || v == Sales[index].Cost {
				Items[Sales[index].ID].Cost[i] = Sales[index].Cost
				Items[Sales[index].ID].Quantity[i] += Sales[index].Quantity
				break
			}
		}
	}

	Sales[index] = Sales[len(Sales)-1]
	Sales = Sales[:len(Sales)-1]
}

func GetLine(selection string, dataType, itemFilter, customerFilter int) ([]string, [][]float32, int) {
	if selection == "" {
		return nil, nil, Debug.Need_More_Info
	}

	raw := strings.Split(selection, "-")

	if len(raw) < 2 {
		return nil, nil, Debug.Need_More_Info
	}

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil, nil, Debug.Invalid_Input
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, nil, Debug.Invalid_Input
	}

	date := []uint8{
		uint8(year),
		uint8(month),
	}

	var sales [][]float32
	var names []string

	for id, x := range Items {
		var totals []float32

		for i := uint8(1); i < 32; i++ {
			var total float32

			for v := len(Sales) - 1; v >= 0; v-- {
				if Sales[v].ID != id || Sales[v].Day != i || Sales[v].Month != date[1] || Sales[v].Year != date[0] {
					continue
				}

				if itemFilter != -1 && Sales[v].ID != uint16(itemFilter) {
					continue
				}

				if customerFilter != -1 && Sales[v].Customer != uint8(customerFilter) {
					continue
				}

				switch dataType {
				case 0:
					total += Sales[v].Price
				case 1:
					total += Sales[v].Cost
				case 2:
					total += Sales[v].Quantity
				}
			}

			totals = append(totals, total)
		}

		if totals != nil {
			names = append(names, x.Name)
			sales = append(sales, totals)
		}
	}
	fmt.Println("Line: ", names, sales)

	return names, sales, Debug.Success
}

func GetPie(selection string, dataType, itemFilter, customerFilter int) ([]string, []float32, int) {

	if selection == "" {
		return nil, nil, Debug.Need_More_Info
	}

	date, errID := unknown.ProcessDate2(selection)
	if errID != Debug.Success {
		return nil, nil, errID
	}

	var sales []float32
	var names []string

	for id, val := range Items {
		var total float32

		for i := len(Sales) - 1; i >= 0; i-- {
			if Sales[i].ID != id || Sales[i].Month != date[1] || Sales[i].Year != date[0] || (Sales[i].Month != date[1] && Sales[i].Month != 0) || (Sales[i].Day != date[2] && Sales[i].Day != 0) {
				continue
			}

			if itemFilter != -1 && Sales[i].ID != uint16(itemFilter) {
				continue
			}

			if customerFilter != -1 && Sales[i].Customer != uint8(customerFilter) {
				continue
			}

			switch dataType {
			case 0:
				total += Sales[i].Price
			case 1:
				total += Sales[i].Cost
			case 2:
				total += Sales[i].Quantity
			}
		}

		if total != 0 {
			names = append(names, val.Name)
			sales = append(sales, total)
		}
	}
	fmt.Println("Pie: ", names, sales)
	return names, sales, Debug.Success
}
