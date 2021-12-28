package Data

import (
	"fmt"
	"strconv"
)

//For the Line chart
func GetProfitForTimes(variant int, targetSheet, subStr string)([][]float64,[]string){
	items := GetAllIDs(targetSheet, subStr)
	labels := make([]string, 0)
	values := make([][]float64, 0)

	for _, v := range items{
		check := GetProfitForItemTimes(v.ID, targetSheet, subStr)
		labels = append(labels, v.Name)
		//revenue: 0, cost: 1, profit, 2;
		values = append(values, check[variant])
	}

	return values, labels
}

func GetProfitForItemTimes(id int, targetSheet, subStr string) [][]float64 {
	revenue := make([]float64, 0)
	cost := make([]float64, 0)
	profit := make([]float64, 0)

	for i := 0; i < 32; i++ {
		newSelect := subStr + "/" + strconv.Itoa(i)

		totals := GetTotalProfit(id, targetSheet, newSelect)
		revenue = append(revenue, totals[0])
		cost = append(cost, totals[1])
		profit = append(profit, totals[2])
	}
	return [][]float64{
		revenue,
		cost,
		profit,
	}
}

//For the Pie Chart
func GetAllProfits(selectionStr string) ([][]float64, []string) {
	targetSheet := "Report Data"
	items := GetAllIDs(targetSheet, selectionStr)
	fmt.Println(items)

	names := make([]string, 0)
	profits := make([]float64, 0)
	revenue := make([]float64, 0)
	costs := make([]float64, 0)

	for _, v := range items {
		names = append(names, v.Name)
		revenue = append(revenue, v.Price)
		costs = append(revenue, v.Cost)
		profits = append(profits, v.Price - v.Cost)
	}

	return [][]float64{
		revenue,
		costs,
		profits,
	}, names
}

func GetTotalProfit(id int, targetSheet, selectionStr string) []float64 {
	totalRevenue := 0.0
	totalCost := 0.0
	totalProfit := 0.0

	res := FindAll(targetSheet, "G", selectionStr, id)

	for i := 0; i < len(res); i++ {
		for _, v := range res {
			rev := F.GetCellValue(targetSheet, "D"+strconv.Itoa(v))
			cost := F.GetCellValue(targetSheet, "E"+strconv.Itoa(v))

			conRev, _ := strconv.ParseFloat(rev, 64)
			conCos, _ := strconv.ParseFloat(cost, 64)
			prof := conRev - conCos

			totalRevenue += conRev
			totalCost += conCos
			totalProfit += prof
		}
	}

	return []float64{
		totalRevenue,
		totalCost,
		totalProfit,
	}
}