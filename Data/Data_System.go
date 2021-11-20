package Data

import (
	"fmt"
	"strconv"
	"time"
)

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

func GetIndex(targetSheet string, id, searchType int) int{
	i:= 1
	cell := f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
	for  {
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		conCell, _ := strconv.Atoi(cell)
		switch searchType {
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

func GetIndexStr(targetSheet, id string, searchType int) int{
	i:= 1
	found := false
	cell := f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
	for  {
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		switch searchType {
		case 2:
			if cell == id{
				found = true
			}
			if found{
				if cell != id{
					return i
				}
			}
			break
		case 1:
			if cell == id{
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
	day, month, year := date.Date()
	return strconv.Itoa(day) + "/" + strconv.Itoa(int(month))  + "/" + strconv.Itoa(year)
}

func ConvertClock(clock time.Time) string {
	hr, min, sec := clock.Clock()
	return strconv.Itoa(hr) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)
}