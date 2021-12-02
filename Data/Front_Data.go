package Data

import (
	"fmt"
	"strconv"
	"strings"
)

func GetTotalProfit(selectionStr string, variants int) ([]string, []float64){
	targetSheet := "Report Data"
	//totalRevenue := 0.0
	//totalCost := 0.0

	names := []string{}
	res := []float64{}

	IDs := GetAllIDs(targetSheet)
	fmt.Println(IDs)
	for _,v  := range IDs{
		results := FindAll(targetSheet, "G", selectionStr, v)
		fmt.Println(results)

		for	_, r := range results{
			name := f.GetCellValue(targetSheet, "B" + strconv.Itoa(r))
			rev := f.GetCellValue(targetSheet, "D" + strconv.Itoa(r))
			cos := f.GetCellValue(targetSheet, "E" + strconv.Itoa(r))
			conRev, _ :=  strconv.ParseFloat(rev, 64)
			conCos, _ := strconv.ParseFloat(cos, 64)
			prof := conRev - conCos

			/*
			temp := Sale{
				v,
				name,
				conRev,
				conCos,
				0,
			}
			 */

			names = append(names, name)

			switch variants {
			case 2:
				res = append(res, prof)
				break
			case 1:
				res = append(res, conCos)
				break
			default:
				res = append(res, conRev)
				break
			}
		}
	}


	return names, res
}

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

func GetData(targetSheet string,id int) []string{
	i := GetIndex(targetSheet, id, 1)

	name := f.GetCellValue(targetSheet, "B"+ strconv.Itoa(i))
	price := f.GetCellValue(targetSheet, "C"+ strconv.Itoa(i))
	cost := f.GetCellValue(targetSheet, "D"+ strconv.Itoa(i))
	quantity := f.GetCellValue(targetSheet, "E"+ strconv.Itoa(i))
	return []string{
		name,
		price,
		cost,
		quantity,
	}
}

func GetAllData(targetSheet string) []Sale{
	var data []Sale

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
	return data
}