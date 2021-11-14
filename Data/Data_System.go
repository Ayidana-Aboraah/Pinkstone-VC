package main

import (
	"fmt"
	"strconv"
	"time"
)

//Create a total profit for a specified period of time

func ReadVal(sheet string) {
	if err != nil {
		fmt.Println(err)
		return
	}

	//Getting a row
	rows := f.GetRows(sheet)

	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

func UpdateLog(item Sale, targetSheet string) {
	if err != nil {
		fmt.Println(err)
		return
	}

	cell := f.GetCellValue(targetSheet, "A1")
	for idx := 1; cell != ""; idx++ {
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(idx))
		if cell == "" {
			fmt.Println("A" + strconv.Itoa(idx))
			f.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), item.id)
			f.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), item.name)
			f.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), item.price)
			f.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), item.cost)
			f.SetCellValue(targetSheet, "E"+strconv.Itoa(idx), time.Now())
			break
		}
		fmt.Println("Not: " + strconv.Itoa(idx))
	}
}

func GetInventory(id int) int {
	//Check the Detection data for the item
	//pull the inventory data from one of the cells
	//Convert the inventory data to a number
	//take the number and subtract by one
	//Return the result number
	targetSheet := "Detection Data"
	cell := f.GetCellValue(targetSheet, "A1")

	for idx := 1; cell != ""; idx++ {
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(idx))
		if cell == strconv.Itoa(id) {
			res := f.GetCellValue(targetSheet, "F"+strconv.Itoa(idx))

			if inven, err := strconv.Atoi(res); err == nil {
				return inven
			}
		}
	}
	fmt.Println("Error: No id found")
	return 0
}

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

func NewAppItem(){
	//Make Sure Camera App is open
	//Check id against database
	//if id is in the database, ask to override
	//else open up new item menu
	//store id in temp sale
	//each box fills in a variable in the temp sale
	//paste the temp sale into the rows in items sheet (use the index used to find the last row before it was empty
	//Close when and the user it's done

	cell := f.GetCellValue("Items", "A1")
		//Convert res to int first
		i := GetIndex("Items", res, 2)
		if strconv.Atoi(cell) == res{
			//Ask user to override,
			//if cancel then close camera
			//if override then open change menu
		}
		if cell == ""{
			//Open New Item menu
			//Open temp and save the id as res
			//Save the index for the cell for later placement of a new id
			break
		}
	}
}

//Specify the time as a parameter
func GetTotalProfit(selectionType int) (revenue, cost, profit float64){
	targetSheet := "Report Data"
	return 0, 0, 0
}

func GetIndex(targetSheet string, id, searchType int) int{
	i:= 1
	cell := f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
	for  {
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		conCell, _ := strconv.Atoi(cell)
		switch searchType {
		case 2:
			if conCell == id{
				return i
			}
			if cell == ""{
				return i
			}
			break
		case 1:
			if conCell == id{
				return i
			}
			break
		default:
			if cell == ""{
				return i
			}
			break
		}

		i++
	}
}

func ConvertDate(date time.Time) string{
	day, month, year := time.Now().Date()
	return strconv.Itoa(day) + "/" + string(month) + "/" + strconv.Itoa(year)
}

func ConvertClock() string {
	hr, min, sec := time.Now().Clock()
	return strconv.Itoa(hr) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)
}