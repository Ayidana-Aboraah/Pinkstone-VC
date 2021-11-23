package Data

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"time"
)

//var f, err = excelize.OpenFile("AppData.xlsx")
var f, _ = excelize.OpenFile("TestAppData.xlsx")


func ReadVal(sheet string) {
	//Getting a row
	rows := f.GetRows(sheet)

	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

func SaveFile(){f.Save()}

func UpdateData(item Sale, targetSheet string, variant int){
	idx := GetIndex(targetSheet, 0, 0)

	switch variant {
	//Update Items [Add]
	case 2:
		checker := GetIndex(targetSheet, item.ID, 1)
		if checker == 0{
			f.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), item.ID)
			f.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), item.Name)
			f.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), item.Price)
			f.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), item.Cost)
			f.SetCellValue(targetSheet, "E"+strconv.Itoa(idx), item.Quantity)
		}else{
			f.SetCellValue(targetSheet, "A"+strconv.Itoa(checker), item.ID)
			f.SetCellValue(targetSheet, "B"+strconv.Itoa(checker), item.Name)
			f.SetCellValue(targetSheet, "C"+strconv.Itoa(checker), item.Price)
			f.SetCellValue(targetSheet, "D"+strconv.Itoa(checker), item.Cost)
			f.SetCellValue(targetSheet, "E"+strconv.Itoa(checker), item.Quantity)
		}
		break
	//Update Report function
	case 1:
		newInven := GetInventory(item.ID) - 1
		f.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), item.ID)
		f.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), item.Name)
		f.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), item.Quantity)
		f.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), item.Price)
		f.SetCellValue(targetSheet, "E"+strconv.Itoa(idx), item.Cost)
		f.SetCellValue(targetSheet, "F"+strconv.Itoa(idx), newInven)
		f.SetCellValue(targetSheet, "G"+strconv.Itoa(idx), ConvertDate(time.Now()))
		f.SetCellValue(targetSheet, "H"+strconv.Itoa(idx), ConvertClock(time.Now()))
		break
		//Update Log function
	default:
		fmt.Println("A" + strconv.Itoa(idx))
		f.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), item.ID)
		f.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), item.Name)
		f.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), item.Price)
		f.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), item.Cost)
		f.SetCellValue(targetSheet, "E"+strconv.Itoa(idx), time.Now())
		break
	}
}

func ModifyItem(item Sale, targetSheet string){
	idx := GetIndex(targetSheet, item.ID, 1)

	f.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), item.Name)
	f.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), item.Price)
	f.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), item.Cost)
	f.SetCellValue(targetSheet, "E"+strconv.Itoa(idx), item.Quantity)
}

func GetInventory(ID int) int {
	targetSheet := "Detection Data"

	idx := GetIndex(targetSheet, ID, 1)

	res := f.GetCellValue(targetSheet, "F"+strconv.Itoa(idx))

	inven, _ := strconv.Atoi(res)
	return inven
}

func GetIndex(targetSheet string, ID, searchType int) int{
	i:= 1
	cell := f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
	for  {
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		conCell, _ := strconv.Atoi(cell)
		switch searchType {
		case 1:
			if conCell == ID{
				return i
			}
			if cell == ""{
				return 0
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

func GetIndexStr(targetSheet, ID string, searchType int) int{
	i:= 1
	found := false
	cell := f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
	for  {
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
		switch searchType {
		case 2:
			if cell == ID{
				found = true
			}
			if found{
				if cell != ID{
					return i
				}
			}
			break
		case 1:
			if cell == ID{
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