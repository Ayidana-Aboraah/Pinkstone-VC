package Data

import (
	"strconv"
	"strings"
)

func GetAllIDs(targetSheet string) []int{
	IDs := []int{}

	cell := f.GetCellValue(targetSheet, "A2")

	for i := 1; cell != ""; {
		conID, _ := strconv.Atoi(cell)
		IDs = append(IDs, conID)
		i++
		cell = f.GetCellValue(targetSheet, "A" + strconv.Itoa(i))
	}

	return IDs
}

func FindAll(targetSheet, targetAxis, subStr string, ID int) []int {
	var idxes []int
	cell := f.GetCellValue(targetSheet, targetAxis + "1")

	for i := 1; cell != "";  i++{
		cell =  f.GetCellValue(targetSheet, targetAxis + strconv.Itoa(i))

		idCell := f.GetCellValue(targetSheet, "A" + strconv.Itoa(i))
		conID, _ := strconv.Atoi(idCell)

		if conID == ID{
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

func GetAllData(targetSheet string, id int) []Sale{
	var data []Sale

	if id == 0{
		cell := f.GetCellValue(targetSheet, "A2")
		for i := 1; cell != ""; i++{
			name := f.GetCellValue(targetSheet, "B"+ strconv.Itoa(i))
			price := f.GetCellValue(targetSheet, "C"+ strconv.Itoa(i))
			cost := f.GetCellValue(targetSheet, "D"+ strconv.Itoa(i))
			quantity := f.GetCellValue(targetSheet, "E"+ strconv.Itoa(i))

			conID, _ := strconv.Atoi(cell)
			p, c, q := ConvertStringToSale(price, cost, quantity)

			temp := Sale{
				ID:    conID,
				Name:  name,
				Price: p,
				Cost:  c,
				Quantity: q,
			}
			data = append(data, temp)
			cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		}
	}else{
		i := GetIndex(targetSheet, id, 1)

		name := f.GetCellValue(targetSheet, "B"+ strconv.Itoa(i))
		price := f.GetCellValue(targetSheet, "C"+ strconv.Itoa(i))
		cost := f.GetCellValue(targetSheet, "D"+ strconv.Itoa(i))
		quantity := f.GetCellValue(targetSheet, "E"+ strconv.Itoa(i))

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