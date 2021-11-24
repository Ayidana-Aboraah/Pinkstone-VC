package Data

import (
	"strconv"
	"strings"
)

//Probably going to be removed and just put in functions in a later thing
func UpdateProfit(item Sale) {
	UpdateData(item, "Price Log", 0)
	ModifyItem(item, "Items")
}

func GetTotalProfit(selectionStr string) (revenue, cost, profit float64){
	targetSheet := "Report Data"
	totalRevenue := 0.0
	totalCost := 0.0

	results := FindAll(targetSheet, "G", selectionStr)
	/*
	for i := 0; i < len(results); {
		rev := f.GetCellValue(targetSheet, "D" + strconv.Itoa(results[i]))
		cos := f.GetCellValue(targetSheet, "E" + strconv.Itoa(results[i]))
		conRev, _ :=  strconv.ParseFloat(rev, 64)
		conCos, _ := strconv.ParseFloat(cos, 64)

		totalRevenue += conRev
		totalCost += conCos
	}
	 */

	for	_, r := range results{
		rev := f.GetCellValue(targetSheet, "D" + strconv.Itoa(r))
		cos := f.GetCellValue(targetSheet, "E" + strconv.Itoa(r))
		conRev, _ :=  strconv.ParseFloat(rev, 64)
		conCos, _ := strconv.ParseFloat(cos, 64)

		totalRevenue += conRev
		totalCost += conCos
	}

	totalProfit := totalRevenue - totalCost
	return totalRevenue, totalCost, totalProfit
}

func FindAll(targetSheet, targetAxis, subStr string) []int {
	var idxes []int
	cell := f.GetCellValue(targetSheet, targetAxis + "1")
	for i := 1; cell != "";  {
		if strings.Contains(cell, subStr) {
			idxes = append(idxes, i)
		}
		i++
		cell = f.GetCellValue(targetSheet, targetAxis+strconv.Itoa(i))
	}
	return idxes
}
