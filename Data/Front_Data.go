package Data

import (
	"fmt"
	"strconv"
	"strings"
)

func GetAllIDs(targetSheet, selectionStr string) []Sale {
	var Items []Sale

	cell := F.GetCellValue(targetSheet, "A2")

	complete := false

	for i := 2; cell != ""; i++ {
		complete = false
		cell = F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		checkCell := F.GetCellValue(targetSheet, "G"+strconv.Itoa(i))

		if strings.Contains(checkCell, selectionStr) {
			name := F.GetCellValue(targetSheet, "B"+strconv.Itoa(i))
			quan := F.GetCellValue(targetSheet, "C"+strconv.Itoa(i))
			rev := F.GetCellValue(targetSheet, "D"+strconv.Itoa(i))
			cos := F.GetCellValue(targetSheet, "E"+strconv.Itoa(i))
			revenue, cost, quantity := ConvertStringToSale(rev, cos, quan)
			conID, _ := strconv.Atoi(cell)

			for idx, v := range Items {
				if v.ID == conID {
					Items[idx].Price += revenue
					Items[idx].Cost += cost
					Items[idx].Quantity += quantity
					complete = true
					break
				}
			}
			if !complete {

				temp := Sale{
					ID:       conID,
					Name:     name,
					Price:    revenue,
					Cost:     cost,
					Quantity: quantity,
				}

				fmt.Println(temp)
				Items = append(Items, temp)
			}
		}
	}

	return Items
}

func FindAll(targetSheet, targetAxis, subStr string, ID int) []int {
	var idxes []int
	cell := F.GetCellValue(targetSheet, targetAxis+"2")

	for i := 2; cell != ""; i++ {
		cell = F.GetCellValue(targetSheet, targetAxis+strconv.Itoa(i))

		idCell := F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		conID, _ := strconv.Atoi(idCell)

		if conID == ID || ID == 0 {
			if strings.Contains(cell, subStr) && !strings.Contains(cell, subStr+"0") && !strings.Contains(cell, subStr+"1") {
				idxes = append(idxes, i)
			}
		}
	}

	return idxes
}

func GetAllData(targetSheet string, id int) []Sale {
	var data []Sale

	if id == 0 {
		cell := F.GetCellValue(targetSheet, "A2")
		for i := 2; cell != ""; {
			name := F.GetCellValue(targetSheet, "B"+strconv.Itoa(i))
			price := F.GetCellValue(targetSheet, "C"+strconv.Itoa(i))
			cost := F.GetCellValue(targetSheet, "D"+strconv.Itoa(i))
			quantity := F.GetCellValue(targetSheet, "E"+strconv.Itoa(i))

			conID, _ := strconv.Atoi(cell)
			p, c, q := ConvertStringToSale(price, cost, quantity)

			temp := Sale{
				ID:       conID,
				Name:     name,
				Price:    p,
				Cost:     c,
				Quantity: q,
			}
			data = append(data, temp)
			i++
			cell = F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		}
	} else {
		i := GetIndex(targetSheet, id, 1)

		name := F.GetCellValue(targetSheet, "B"+strconv.Itoa(i))
		price := F.GetCellValue(targetSheet, "C"+strconv.Itoa(i))
		cost := F.GetCellValue(targetSheet, "D"+strconv.Itoa(i))
		quantity := F.GetCellValue(targetSheet, "E"+strconv.Itoa(i))

		p, c, q := ConvertStringToSale(price, cost, quantity)
		temp := Sale{
			ID: id,
			Quantity: q,
			Price: p,
			Cost: c,
			Name: name,
		}
		data = append(data, temp)
	}

	return data
}
