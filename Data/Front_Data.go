package Data

import (
	"strconv"
	"strings"
)

func GetAllIDs(targetSheet, selectionStr string) ([]int, []string) {
	IDs := make([]int,0)
	Names := make([]string,0)

	cell := F.GetCellValue(targetSheet, "A2")

	for i := 2; cell != ""; i++ {
		cell = F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		checkCell := F.GetCellValue(targetSheet, "G"+strconv.Itoa(i))

		if strings.Contains(checkCell, selectionStr) {
			name := F.GetCellValue(targetSheet, "B"+strconv.Itoa(i))
			conID, _ := strconv.Atoi(cell)

			for _, v := range IDs{
				if v == conID{
					i++
				}
			}

			Names = append(Names, name)
			IDs = append(IDs, conID)
		}
	}

	return IDs, Names
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
