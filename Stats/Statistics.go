package Stats

import (
	"BronzeHermes/Database"
	"fmt"
	"time"
)

const (
	ONCE uint8 = iota
	MONTHLY
	YEARLY
)

func CompileReport(selection uint8, date []int) (string, [4]float64) {
	// For Date: 0 = day, 1 = month, 2 = year, 3 = now

	if len(date) == 0 {
		year, month, day := time.Now().Local().Date()
		date = []int{day, int(month), year}
	}

	var item_sales [3]float64
	var damages float64

	for _, v := range Database.Sales {
		vt := time.Unix(v.Timestamp, 0)

		if vt.Year() != date[YEARLY] {
			continue
		}

		if selection != YEARLY && int(vt.Month()) != date[MONTHLY] {
			continue
		}

		if selection == ONCE && vt.Day() != date[ONCE] {
			continue
		}

		item_sales[0] += v.Price * v.Quantity
		item_sales[1] += v.Cost * v.Quantity
		item_sales[2] += (v.Price - v.Cost) * v.Quantity

		if v.Customer == 0 {
			damages += v.Cost * v.Quantity
		}
	}

	return fmt.Sprintf(
		"Item Revenue: %.2f,\nItem Cost: %.2f\nDamages: %.2f,\nReport Total: %.2f",
		item_sales[0],  // Gain
		-item_sales[1], // Cost
		-damages,       // Damages
		item_sales[2],  // Report Total
	), [4]float64{item_sales[0], -item_sales[1], -damages, item_sales[2]}
}

func RemoveFromSales(index int) {

	if Database.Items[Database.Sales[index].ID].Quantity[0] <= 0 {
		Database.Items[Database.Sales[index].ID].Cost[0] = Database.Sales[index].Cost
		Database.Items[Database.Sales[index].ID].Quantity[0] = Database.Sales[index].Quantity
	} else {
		for i, v := range Database.Items[Database.Sales[index].ID].Cost {
			if v == 0 || v == Database.Sales[index].Cost {
				Database.Items[Database.Sales[index].ID].Cost[i] = Database.Sales[index].Cost
				Database.Items[Database.Sales[index].ID].Quantity[i] += Database.Sales[index].Quantity
				break
			}
		}
	}

	Database.Sales[index] = Database.Sales[len(Database.Sales)-1]
	Database.Sales = Database.Sales[:len(Database.Sales)-1]
}

// func GetLine(selection string, dataType, itemFilter, customerFilter int) ([]string, [][]float32, int) {
// 	if selection == "" {
// 		return nil, nil, Helpers.Need_More_Info
// 	}

// 	raw := strings.Split(selection, "-")

// 	if len(raw) < 2 {
// 		return nil, nil, Helpers.Need_More_Info
// 	}

// 	year, err := strconv.Atoi(raw[0][1:])
// 	if err != nil {
// 		return nil, nil, Helpers.Invalid_Input
// 	}

// 	month, err := strconv.Atoi(raw[1])
// 	if err != nil {
// 		return nil, nil, Helpers.Invalid_Input
// 	}

// 	date := []uint8{
// 		uint8(year),
// 		uint8(month),
// 	}

// 	var sales [][]float32
// 	var names []string

// 	for id, x := range Items {
// 		var totals []float32

// 		for i := 1; i < 32; i++ {
// 			var total float32

// 			for v := len(Sales) - 1; v >= 0; v-- {
// 				vt := time.Unix(Sales[v].Timestamp, 0)
// 				if Sales[v].ID != id || vt.Day() != i || vt.Month() != date[1] || vt.Year() != date[0] {
// 					continue
// 				}

// 				if itemFilter != -1 && Sales[v].ID != uint16(itemFilter) {
// 					continue
// 				}

// 				if customerFilter != -1 && Sales[v].Customer != uint16(customerFilter) {
// 					continue
// 				}

// 				switch dataType {
// 				case 0:
// 					total += Sales[v].Price
// 				case 1:
// 					total += Sales[v].Cost
// 				case 2:
// 					total += Sales[v].Quantity
// 				}
// 			}

// 			totals = append(totals, total)
// 		}

// 		if totals != nil {
// 			names = append(names, x.Name)
// 			sales = append(sales, totals)
// 		}
// 	}
// 	fmt.Println("Line: ", names, sales)

// 	return names, sales, Helpers.Success
// }

// func GetPie(selection string, dataType, itemFilter, customerFilter int) ([]string, []float32, int) {

// 	if selection == "" {
// 		return nil, nil, Helpers.Need_More_Info
// 	}

// 	date, errID := Helpers.ProcessDate2(selection)
// 	if errID != Helpers.Success {
// 		return nil, nil, errID
// 	}

// 	var sales []float32
// 	var names []string

// 	for id, val := range Items {
// 		var total float32

// 		for i := len(Sales) - 1; i >= 0; i-- {
// 			if Sales[i].ID != id || Sales[i].Month != date[1] || Sales[i].Year != date[0] || (Sales[i].Month != date[1] && Sales[i].Month != 0) || (Sales[i].Day != date[2] && Sales[i].Day != 0) {
// 				continue
// 			}

// 			if itemFilter != -1 && Sales[i].ID != uint16(itemFilter) {
// 				continue
// 			}

// 			if customerFilter != -1 && Sales[i].Customer != uint8(customerFilter) {
// 				continue
// 			}

// 			switch dataType {
// 			case 0:
// 				total += Sales[i].Price
// 			case 1:
// 				total += Sales[i].Cost
// 			case 2:
// 				total += Sales[i].Quantity
// 			}
// 		}

// 		if total != 0 {
// 			names = append(names, val.Name)
// 			sales = append(sales, total)
// 		}
// 	}
// 	fmt.Println("Pie: ", names, sales)
// 	return names, sales, Helpers.Success
// }
