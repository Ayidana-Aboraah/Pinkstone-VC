package Data

import (
	"fmt"
	"strconv"
	"strings"
)

func GetAllIDs(targetSheet, selectionStr string) ([]int, []string) {
	var IDs []int
	var Names []string

	checkCell := F.GetCellValue(targetSheet, "G2")

	for i := 2; checkCell != ""; i++ {
		complete := false

		cell := F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		checkCell = F.GetCellValue(targetSheet, "G"+strconv.Itoa(i))

		if !strings.Contains(checkCell, selectionStr) {
			continue
		}

		conID, _ := strconv.Atoi(cell)

		for _, v := range IDs {
			if v != conID {continue}
			complete = true
			break
		}

		if complete {
			continue
		}

		name := F.GetCellValue(targetSheet, "B"+strconv.Itoa(i))
		fmt.Println(conID)
		IDs = append(IDs, conID)
		Names = append(Names, name)
	}

	return IDs, Names
}

func GetAllData(targetSheet, selectionStr string) []Sale {
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

		if conID != ID || ID != 0 {	continue }

		if strings.Contains(cell, subStr) && !strings.Contains(cell, subStr+"0") && !strings.Contains(cell, subStr+"1") {
			idxes = append(idxes, i)
		}
	}

	return idxes
}

//Gets Data from Items
func GetData(targetSheet string, id int) []Sale {
	var data []Sale

	//id 0 means all
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
		i := GetIndex(targetSheet, id, true)

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
