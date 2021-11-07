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
//Create a total profit for a specified period of time

func main(){
	UpdateLog(Sale{
		id:0,
		name: "Null",
		price: 0,
	},	"Log")

	ReadVal("Log")
	f.Save()
}

func ReadVal(sheet string){
	if err != nil{
		fmt.Println(err)
		return
	}

	//Getting a row
	rows := f.GetRows(sheet)

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
				f.SetCellValue(targetSheet, "G"+strconv.Itoa(idx), time.Now())
				UpdateLog(*shoppingCart[i], "Log")
				//Place the profit made from the item (revenue - cost)
				//Show the decrease in inventory
				newInven := UpdateInvenotry(shoppingCart[i].id)
				f.SetCellValue(targetSheet, "F"+strconv.Itoa(idx), strconv.Itoa(newInven - 1))
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

func UpdateLog(item Sale, targetSheet string) {
	if err != nil {
		fmt.Println(err)
		return
	}

	cell := f.GetCellValue(targetSheet, "A1")
	for idx := 1; cell != ""; idx++{
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(idx))
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

func UpdateInvenotry(id int) int{
	//Check the Detection data for the item
	//pull the inventory data from one of the cells
	//Convert the inventory data to a number
	//take the number and subtract by one
	//Return the result number
	targetSheet := "Detection Data"
	cell := f.GetCellValue(targetSheet, "A1")

	for idx := 1; cell != ""; idx++{
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(idx))
		if cell == strconv.Itoa(id) {
			res := f.GetCellValue(targetSheet, "F" + strconv.Itoa(idx))

			if inven, err := strconv.Atoi(res); err == nil{
				return inven
				}
		}
	}
	fmt.Println("Error: No id found")
	return 0
}

//Changes
func NewPrice(newPrice float64, item Sale){
	item.price = newPrice
	UpdateLog(item, "Price Log")
}