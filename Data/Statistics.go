package Data

import (
	"strconv"
)

func GetTotalProfit(id int, targetSheet, selectionStr string) []float64{
	totalRevenue := 0.0
	totalCost := 0.0
	totalProfit := 0.0

	res := FindAll(targetSheet, "G", selectionStr, id)

	for i := 0; i < len(res); i++{
		for _, v := range res{
			rev := f.GetCellValue(targetSheet, "D" + strconv.Itoa(v))
			cost := f.GetCellValue(targetSheet, "E" + strconv.Itoa(v))

			conRev, _ :=  strconv.ParseFloat(rev, 64)
			conCos, _ := strconv.ParseFloat(cost, 64)
			prof := conRev - conCos

			totalRevenue += conRev
			totalCost += conCos
			totalProfit += prof
		}
	}

	return  []float64{
		totalRevenue,
		totalCost,
		totalProfit,
	}
}

func GetProfitForTimes(variant, id int, subStr string) []float64{
	targetSheet := "Report Data"
	profit := []float64{}

	for i := 0; i < 32; i++ {
		newSelect := subStr + "/" + strconv.Itoa(i)

		totals := GetTotalProfit(id, targetSheet, newSelect)
		//0 revenue; 1 cost; 2 profit
		profit = append(profit, totals[variant])
	}
	return profit
}