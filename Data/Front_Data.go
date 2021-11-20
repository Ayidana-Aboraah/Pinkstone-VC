package Data

import (
	"github.com/makiuchi-d/gozxing"
	"strconv"
	"time"
)

//Create a total profit for a specified period of time
var tempIndex int
var tempSale *Sale

func NewProfit(variant int, newValue float64, item Sale) {
	switch variant {
	case 0:
		item.price = newValue
		break
	case 1:
		item.cost = newValue
		break
	}
	UpdateLog(item, "Price Log")
	UpdateLog(item, "Items")
}

func NewAppItem(res gozxing.Result) {
	//Make Sure Camera App is open
	//Check id against database
	//if id is in the database, ask to override
	//else open up new item menu
	//store id in temp sale
	//each box fills in a variable in the temp sale
	//paste the temp sale into the rows in items sheet use the index used to find the last row before it was empty
	//Close when and the user it's done

	//Insert image from cam for this.
	convRes := res.GetNumBits()
	//Convert res to int first
	ix := GetIndex("Items", convRes, 0)
	iy := GetIndex("Items", convRes, 1)

	//(if)Checking if the id needs to be overwritten
	//(else)Checking if a new id needs to be made
	if ix > 1 {
		//Ask user to override,
		//if cancel then close camera
		//if override then open change menu
	} else {
		//Open New Item menu
		//OpenItem()
		//Open temp and save the id as res
		tempSale.id = convRes
		//Save the index for the cell for later placement of a new id
		tempIndex = iy
	}
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