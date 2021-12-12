package Data

import (
	"fmt"
	"strconv"
)

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

func GetProfitForTimes(variant int, targetSheet, subStr string)([][]float64,[]string){
	ids, labels := GetAllIDs(targetSheet, subStr)
	values := make([][]float64, 0)

	for _, v := range ids{
		check := GetProfitForItemTimes(v, targetSheet, subStr)
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
		//0 revenue; 1 cost; 2 profit
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

func GetAllProfits(selectionStr string) ([][]float64, []string) {
	targetSheet := "Report Data"
	IDs, Names := GetAllIDs(targetSheet, selectionStr)
	profits := make([]float64, 0)
	revenue := make([]float64, 0)
	costs := make([]float64, 0)

	fmt.Println(IDs)

	for _, v := range IDs {
		totals := GetTotalProfit(v, targetSheet, selectionStr)
		revenue = append(revenue, totals[0])
		costs = append(revenue, totals[1])
		profits = append(profits, totals[2])
	}

	return [][]float64{
		revenue,
		costs,
		profits,
	}, Names
}

func ProcessAllProfit(values []float64) float64{
	total := 0.0

	for _, v := range values{
		total += v
	}

	return total
}
