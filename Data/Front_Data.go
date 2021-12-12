package Data

import (
	"fmt"
	"strconv"
	"strings"
)

func GetAllIDs(targetSheet, selectionStr string) []Sale {
	Items := make([]Sale, 0)

	cell := F.GetCellValue(targetSheet, "A2")

	complete := false

	for i := 2; cell != ""; i++{
		complete = false
		cell = F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		checkCell := F.GetCellValue(targetSheet, "G"+strconv.Itoa(i))

		if strings.Contains(checkCell, selectionStr) {
			name := F.GetCellValue(targetSheet, "B"+strconv.Itoa(i))
			rev := F.GetCellValue(targetSheet, "D"+strconv.Itoa(i))
			cos := F.GetCellValue(targetSheet, "E"+strconv.Itoa(i))
			revenue, cost, _ := ConvertStringToSale(rev, cos, "")
			conID, _ := strconv.Atoi(cell)

			for idx, v := range Items{
				if v.ID == conID{
					Items[idx].Price += revenue
					Items[idx].Cost += cost
					complete = true
					break
				}
			}
			if !complete{



				temp := Sale{
					ID: conID,
					Name: name,
					Price: revenue,
					Cost: cost,
					Quantity: 1,
				}

				fmt.Println(temp)
				Items = append(Items, temp)
			}
		}
	}

	Items = Items[:len(Items) -1]

	return Items
}

func FindAll(targetSheet, targetAxis, subStr string, ID int) []int {
	var idxes []int
	cell := F.GetCellValue(targetSheet, targetAxis+"2")

	for i := 2; cell != ""; i++ {
		cell = F.GetCellValue(targetSheet, targetAxis+strconv.Itoa(i))

		idCell := F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		conID, _ := strconv.Atoi(idCell)

		if conID == ID {
			if strings.Contains(cell, subStr) {
				idxes = append(idxes, i)
			}
		}

		if ID == 0 {
			if strings.Contains(cell, subStr) {
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
		for i := 2; cell != ""; i++ {
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
			id,
			name,
			p,
			c,
			q,
		}
		data = append(data, temp)
	}

	return data
}
