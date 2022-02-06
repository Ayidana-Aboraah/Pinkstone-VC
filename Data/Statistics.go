package Data

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// GetProfitForTimes For the Line chart
func GetProfitForTimes(variant int, targetSheet, subStr string) ([][]float64, []string) {
	IDs, Names := GetAllIDs(targetSheet, subStr)

	var values [][]float64

	for i := range IDs {
		check := GetProfitForItemTimes(IDs[i], targetSheet, subStr)
		values = append(values, check[variant])
	}

	return values, Names
}

func GetProfitForItemTimes(id int, targetSheet, subStr string) [][]float64 {
	var revenue []float64
	var cost []float64
	var profit []float64

	for i := 0; i < 32; i++ {
		newSelect := subStr + "/" + strconv.Itoa(i)
		if !strings.Contains(subStr, "/") {
			newSelect = strconv.Itoa(time.Now().Year()) + "/" + strconv.Itoa(int(time.Now().Month())) + "/" + strconv.Itoa(i)
		}

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

// GetAllProfits For the Pie Chart
func GetAllProfits(selectionStr string) ([][]float64, []string) {
	items := GetAllData("Report Data", selectionStr)
	fmt.Println(items)

	var names []string
	var revenue []float64
	var costs []float64
	var profits []float64

	for _, v := range items {
		names = append(names, v.Name)
		revenue = append(revenue, v.Price)
		costs = append(revenue, v.Cost)
		profits = append(profits, v.Price-v.Cost)
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

	var (
		targetAxis string
		log        = false
	)

	if strings.Contains(targetSheet, "Log") || strings.Contains(targetSheet, "log") {
		log = true
	}
	if log {
		targetAxis = "E"
	} else {
		targetSheet = "G"
	}

	results := FindAll(targetSheet, targetAxis, selectionStr, id)

	for _, v := range results {
		rev := ""
		cost := ""
		quan := ""

		if log {
			rev = F.GetCellValue(targetSheet, "C"+strconv.Itoa(v))
			cost = F.GetCellValue(targetSheet, "D"+strconv.Itoa(v))
			quan = "1"
		} else {
			rev = F.GetCellValue(targetSheet, "D"+strconv.Itoa(v))
			cost = F.GetCellValue(targetSheet, "E"+strconv.Itoa(v))
			quan = F.GetCellValue(targetSheet, "C"+strconv.Itoa(v))
		}

		conRev, _ := strconv.ParseFloat(rev, 64)
		conCos, _ := strconv.ParseFloat(cost, 64)
		quantity, _ := strconv.Atoi(quan)
		prof := conRev - conCos

		totalRevenue += conRev * float64(quantity)
		totalCost += conCos * float64(quantity)
		totalProfit += prof * float64(quantity)
	}

	return []float64{
		totalRevenue,
		totalCost,
		totalProfit,
	}
}

func GetSalesForTime(selectionStr string) ([]float64, []string){
	targetSheet := "Report Data"
	var sales []float64

	IDs, names := GetAllIDs(targetSheet, selectionStr)

	for index := range IDs{
		checkCell := F.GetCellValue(targetSheet, "G1")
		mSales := 0

		for i := 1; checkCell != ""; i++{
			checkCell = F.GetCellValue(targetSheet, "G" + strconv.Itoa(i))
			println(checkCell)

			if !strings.Contains(checkCell, selectionStr){continue}

			idCell := F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
			conID, _ := strconv.Atoi(idCell)

			if conID != IDs[index] {continue}

			SalesCell := F.GetCellValue(targetSheet, "C"+strconv.Itoa(i))
			println(SalesCell)
			tempSale, _ := strconv.Atoi(SalesCell)


			mSales += tempSale
			println(i)
			println(tempSale)
			println(mSales)
		}

		sales = append(sales, float64(mSales))
	}

	//GetIDS(targetSheet, selectionStr)
	//Go through each data and see if it contains selectionStr
	// If yes -> Take the targetAxis;
	// If no -> Skip to next
	// if Null, check next and then return
	return sales, names
}