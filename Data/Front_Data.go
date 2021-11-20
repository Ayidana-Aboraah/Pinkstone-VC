package Data

import (
	"strconv"
	"time"
)

//Create a total profit for a specified period of time

func NewProfit(variant int, newValue float64, item Sale) {
	switch variant {
	case 0:
		item.price = newValue
		break
	case 1:
		item.cost = newValue
		break
	}
	UpdateData(item, "Price Log", 0)
	UpdateData(item, "Items", 0)
}

//Specify the time as a parameter
func GetTotalProfit(selectionType int) (revenue, cost, profit float64){
	targetSheet := "Report Data"
	total_revenue := 0.0
	total_cost := 0.0
	switch selectionType {
	//The Day's profit
	default:
		conDate := ConvertDate(time.Now())
		startIdx := GetIndexStr(targetSheet, conDate, 1)
		endIdx := GetIndexStr(targetSheet, conDate, 2)
		for idx := startIdx; idx < endIdx;{
			rev := f.GetCellValue(targetSheet, "D" + strconv.Itoa(idx))
			cos := f.GetCellValue(targetSheet, "E" + strconv.Itoa(idx))

			conRev, _ :=  strconv.ParseFloat(rev, 64)
			conCos, _ := strconv.ParseFloat(cos, 64)

			total_revenue += conRev
			total_cost += conCos
			idx++
		}
		total_profit := total_revenue - total_cost
		return total_revenue,total_cost, total_profit
	}
}