package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"time"
)

// Defining the info used later
type Sale struct {
	id        int
	name      string
	price     float64
}
var f, err = excelize.OpenFile("AppData.xlsx")

var shoppingCart []*Sale
//var totalInventory int

//Add a total for all sales

func main(){
	UpdateLog(Sale{
		id:7,
		name: "Test",
		price: 2,
	})
	ReadVal()
	f.Save()
}

func CreateNewFile(){
	nf := excelize.NewFile()
	idx := nf.NewSheet("Detection Data")
	nf.SetCellValue("Detection Data", "A1", "ID")
	nf.SetCellValue("Long Data", "A2", 100)
	nf.SetActiveSheet(idx)

	if err := f.SaveAs("AppData.xlsx"); err != nil{
		fmt.Println(err)
	}
}

func ReadVal(){
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

func BuyCart(){
	if err != nil{
		fmt.Println(err)
		return
	}

	//Cycle through and check if there is an open space
	//When finding an empty space start unloading the contents of the cart into each cell
	//Until the cart's info is fully displayed
	targetSheet := "Report Data"
	cell := f.GetCellValue(targetSheet, "A1")

	for idx := 1; cell != ""; idx++{
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(idx))
		if cell == "" {
			fmt.Println("Staring at A" + strconv.Itoa(idx))
			for i := 0; i <= len(shoppingCart);{
				f.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), shoppingCart[i].id)
				f.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), shoppingCart[i].name)
				f.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), shoppingCart[i].price)
				//Place the cost of the product(for the seller)
				//Place the profit made from the item (revenue - cost)
				//Show the decrease in inventory
				f.SetCellValue(targetSheet, "F"+strconv.Itoa(idx), time.Now())
				i++
			}
		}
		fmt.Println("Not: " + strconv.Itoa(idx))
	}

	//Clear cart
	ClearCart()
}

func AddToCart(id int){
	//Check the "Detection Data" for the specified ID
	//Grab the row of data
	//Convert the row data to a Sale variable
	//Add the sale variable to the ShoppingCart array
	targetSheet := "Detection Data"
	cell := f.GetCellValue(targetSheet, "A1")
	for idx := 1; strconv.Itoa(id) != cell; idx++{
		cell = f.GetCellValue("Detection Data", "A"+strconv.Itoa(idx))
		if strconv.Itoa(id) == cell {
			fmt.Println("A" + strconv.Itoa(idx))
			if p, err := strconv.ParseFloat(f.GetCellValue(targetSheet, "B" + strconv.Itoa(idx)), 64); err == nil{
				temp := Sale{
					id: id,
					name: f.GetCellValue(targetSheet, "B" + strconv.Itoa(idx)),
					price: p,
				}
				shoppingCart = append(shoppingCart, &temp)
			}
		}
		fmt.Println("Not: " + strconv.Itoa(idx))
	}
}

// RemoveFromCart [Untested]
func RemoveFromCart(id int){
	//Cycle through the cart to find the item id
	//If the id doesn't exist; display error
	//Otherwise; Remove the specified item from the cart
	for i := 0; i != len(shoppingCart);{
		if shoppingCart[i].id == id{
			shoppingCart[i] = shoppingCart[len(shoppingCart)-1] // Copy last element to index i.
			shoppingCart[len(shoppingCart)-1] = &Sale{}   // Erase last element (write zero value).
			shoppingCart = shoppingCart[:len(shoppingCart)-1]   // Truncate slice.
			break
		}
		fmt.Println("Not this one...")
		i++
	}
}


func ClearCart(){
	shoppingCart = shoppingCart[:0]
}

func UpdateLog(item Sale) {
	if err != nil {
		fmt.Println(err)
		return
	}

	cell := f.GetCellValue("Log", "A1")
	targetSheet := "Log"
	/*i := 1
	for {
		cell = f.GetCellValue("Log", "A"+strconv.Itoa(i))
		if cell == "" {
			fmt.Println("A" + strconv.Itoa(i))
			f.SetCellValue(targetSheet, "A"+strconv.Itoa(i), item.id)
			f.SetCellValue(targetSheet, "B"+strconv.Itoa(i), item.name)
			f.SetCellValue(targetSheet, "C"+strconv.Itoa(i), item.price)
			f.SetCellValue(targetSheet, "D"+strconv.Itoa(i), time.Now())
			break
		}
		fmt.Println("No")
		i++
	}

	 */

	for idx := 1; cell != ""; idx++{
		cell = f.GetCellValue("Log", "A"+strconv.Itoa(idx))
		if cell == "" {
			fmt.Println("A" + strconv.Itoa(idx))
			f.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), item.id)
			f.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), item.name)
			f.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), item.price)
			f.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), time.Now())
			break
		}
		fmt.Println("Not: " + strconv.Itoa(idx))
	}
}

func ConvertNowDate() string{
	day := time.Now().Day()
	month := time.Now().Month()
	year := time.Now().Year()
	return string(day) + "/" + string(month) + "/" + string(year)
}