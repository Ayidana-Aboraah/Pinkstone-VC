package Data

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"strconv"
	"time"
)

//var F, err = excelize.OpenFile("Assets/AppData.xlsx")
var F, Err = excelize.OpenFile("TestAppData.xlsx")
//save a back-up.
//if an error with reading the file is met; save the back up as the file and try reading again.

func ReadVal(sheet string) {
	//Getting a row
	rows := F.GetRows(sheet)

	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

func SaveFile(){
	err := F.Save()
	if err != nil{
		fmt.Println(err)
	}
}

func SaveBackUp(sourceFile, backUpfile string){
	input, err := ioutil.ReadFile(sourceFile)
	//input, err := ioutil.ReadFile("Assets/" + sourceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("Assets/" + backUpfile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", backUpfile)
		fmt.Println(err)
		return
	}
}

func UpdateData(item Sale, targetSheet string, variant int){
	idx := GetIndex(targetSheet, 0, 0)

	switch variant {
	//Update Items
	case 2:
		checker := GetIndex(targetSheet, item.ID, 1)
		if checker == 0{
			F.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), item.ID)
			F.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), item.Name)
			F.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), item.Price)
			F.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), item.Cost)
			F.SetCellValue(targetSheet, "E"+strconv.Itoa(idx), item.Quantity)
			break
		}
			F.SetCellValue(targetSheet, "A"+strconv.Itoa(checker), item.ID)
			F.SetCellValue(targetSheet, "B"+strconv.Itoa(checker), item.Name)
			F.SetCellValue(targetSheet, "C"+strconv.Itoa(checker), item.Price)
			F.SetCellValue(targetSheet, "D"+strconv.Itoa(checker), item.Cost)
			F.SetCellValue(targetSheet, "E"+strconv.Itoa(checker), item.Quantity)
		break
	//Update Report function
	case 1:
		inventory := SetInventory(item.ID, item.Quantity)
		F.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), item.ID)
		F.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), item.Name)
		F.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), item.Quantity)
		F.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), item.Price)
		F.SetCellValue(targetSheet, "E"+strconv.Itoa(idx), item.Cost)
		F.SetCellValue(targetSheet, "F"+strconv.Itoa(idx), inventory)
		F.SetCellValue(targetSheet, "G"+strconv.Itoa(idx), ConvertDate(time.Now()))
		break
		//Update Log function
	default:
		F.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), item.ID)
		F.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), item.Name)
		F.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), item.Price)
		F.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), item.Cost)
		F.SetCellValue(targetSheet, "E"+strconv.Itoa(idx), ConvertDate(time.Now()))
		break
	}
}

func SetInventory(ID, amount int) int{
	targetSheet := "Items"

	idx := GetIndex(targetSheet, ID, 1)

	res := F.GetCellValue(targetSheet, "F"+strconv.Itoa(idx))

	inven, _ := strconv.Atoi(res)

	newInventory := inven - amount

	F.SetCellValue(targetSheet, "F" + strconv.Itoa(idx), newInventory)
	return newInventory
}

func GetIndex(targetSheet string, ID, searchType int) int{
	i:= 1
	cell := F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
	for  {
		cell = F.GetCellValue(targetSheet, "A"+strconv.Itoa(i))
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

func ConvertDate(date time.Time) string{
	year, month, day := date.Date()
	return strconv.Itoa(year)+ "/"+ strconv.Itoa(int(month)) + "/" + strconv.Itoa(day)
}