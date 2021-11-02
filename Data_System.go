package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"time"
)

// Defining the info used later
type Sale struct {
	id        int
	name      string
	price     float64
}
var shoppingCart []Sale
//var totalInventory int

//Add a total for all sales

func main(){
	red := Sale{
		id:1,
		name: "Hoe",
		price: 50,
	}
	UpdateLog(red)
	ReadVal()
}

func CreateNewFile(){
	f := excelize.NewFile()
	idx := f.NewSheet("Detection Data")
	f.SetCellValue("Detection Data", "A1", "ID")
	f.SetCellValue("Long Data", "A2", 100)
	f.SetActiveSheet(idx)

	if err := f.SaveAs("LTAppData.xlsx"); err != nil{
		fmt.Println(err)
	}
}

func ReadVal(){
	f, err := excelize.OpenFile("LTAppData.xlsx")
	if err != nil{
		fmt.Println(err)
		return
	}

	//Getting a row
	rows := f.GetRows("Log")

	for _, row := range rows{
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}

func BuyItems(){
	f, err := excelize.OpenFile("LTAppData.xlsx")
	if err != nil{
		fmt.Println(err)
		return
	}
	cell := f.GetCellValue("Report Data", "H2")

	for  i := 2; cell != "";{
		cell := f.GetCellValue("Log", "D" + string(i))
		if cell == ""{
			for x := 0; x < len(shoppingCart);{
				f.SetCellValue("Report Data", "A"+ string(i), shoppingCart[x].id)
				f.SetCellValue("Report Data", "B"+ string(i), shoppingCart[x].name)
				f.SetCellValue("Report Data", "C"+ string(i), shoppingCart[x].price)
				f.SetCellValue("Report Data", "D"+ string(i), ConvertNowDate())
/*
				f.SetCellValue("Report Data", "E"+ string(i), shoppingCart[x])
				f.SetCellValue("Report Data", "F"+ string(i), shoppingCart[x])
				f.SetCellValue("Report Data", "G"+ string(i), shoppingCart[x])
				f.SetCellValue("Report Data", "H"+ string(i), shoppingCart[x])
 */
			}
		}
	}

	//Clear cart
	//shoppingCart = nil
}

func UpdateLog(item Sale) {
	f, err := excelize.OpenFile("LTAppData.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}

	i := 2
	cell := f.GetCellValue("Log", "A" + string(i))
	/*
	for {
		i := 2
		cell = f.GetCellValue("Log", "A"+string(i))
		if cell == "" {
			f.SetCellValue("Report Data", "A"+string(i), item.id)
			f.SetCellValue("Report Data", "B"+string(i), item.name)
			f.SetCellValue("Report Data", "C"+string(i), item.price)
			f.SetCellValue("Report Data", "D"+string(i), ConvertNowDate())
			break
		}
	}
	 */
	if cell == ""{
		f.SetCellValue("Report Data", "A"+string(i), item.id)
		f.SetCellValue("Report Data", "B"+string(i), item.name)
		f.SetCellValue("Report Data", "C"+string(i), item.price)
		f.SetCellValue("Report Data", "D"+string(i), ConvertNowDate())
	}else{i++}
}

func ConvertNowDate() string{
	day := time.Now().Day()
	month := time.Now().Month()
	year := time.Now().Year()
	return string(day) + "/" + string(month) + "/" + string(year)
}