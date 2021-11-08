package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"time"
)

// Defining the info used later
type Sale struct {
	id       int
	name     string
	price    float64
	cost 	 float64
	quantity int
}

var f, err = excelize.OpenFile("AppData.xlsx")

var shoppingCart []*Sale

//Create a total profit for a specified period of time

func main() {
	UpdateLog(Sale{
		id:    0,
		name:  "Null",
		price: 0,
	}, "Log")

	ReadVal("Log")
	f.Save()
}

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

func BuyCart() {
	if err != nil {
		fmt.Println(err)
		return
	}

	//Cycle through and check if there is an open space
	//When finding an empty space start unloading the contents of the cart into each cell
	//Until the cart's info is fully displayed
	targetSheet := "Report Data"
	cell := f.GetCellValue(targetSheet, "A1")

	for idx := 1; cell != ""; idx++ {
		cell = f.GetCellValue(targetSheet, "A"+strconv.Itoa(idx))
		if cell == "" {
			fmt.Println("Staring at A" + strconv.Itoa(idx))
			for i := 0; i <= len(shoppingCart); {
				newInven := UpdateInventory(shoppingCart[i].id) - 1
				f.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), shoppingCart[i].id)
				f.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), shoppingCart[i].name)
				f.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), shoppingCart[i].quantity)
				f.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), shoppingCart[i].price * float64(shoppingCart[i].quantity))
				f.SetCellValue(targetSheet, "F"+strconv.Itoa(idx), time.Now())
				f.SetCellValue(targetSheet, "F"+strconv.Itoa(idx), newInven)
				UpdateLog(*shoppingCart[i], "Log")
				i++
			}
		}
		fmt.Println("Not: " + strconv.Itoa(idx))
	}

	//Clear cart
	ClearCart()
}

//[Untested]
func AddToCart(id int) {
	//Check the "Detection Data" for the specified ID
	//Grab the row of data
	//Convert the row data to a Sale variable
	//Add the sale variable to the ShoppingCart array
	for {
		i := 0
		if i < len(shoppingCart) {
			if shoppingCart[i].id == id {
				shoppingCart[i].quantity++
				break
			}
		}else{
			targetSheet := "Detection Data"
			cell := f.GetCellValue(targetSheet, "A1")
			for idx := 1; strconv.Itoa(id) != cell; idx++ {
				cell = f.GetCellValue("Detection Data", "A"+strconv.Itoa(idx))
				if strconv.Itoa(id) == cell {
					fmt.Println("A" + strconv.Itoa(idx))
					if p, err := strconv.ParseFloat(f.GetCellValue(targetSheet, "B"+strconv.Itoa(idx)), 64); err == nil {
						temp := Sale{
							id:    id,
							name:  f.GetCellValue(targetSheet, "B"+strconv.Itoa(idx)),
							price: p,
						}
						shoppingCart = append(shoppingCart, &temp)
						break
					}
				}
				fmt.Println("Not: " + strconv.Itoa(idx))
			}
		}
	}
}
//[Untested]
func DecreaseFromCart(id int){
	for i := 0; i < len(shoppingCart); {
		if shoppingCart[i].id == id {
			if		shoppingCart[i].quantity- 1 > 0{
				shoppingCart[i].quantity-= 1
			}else{
				RemoveFromCart(i)
			}
		}
	}
}

// RemoveFromCart [Untested]
func RemoveFromCart(i int) {
	//Cycle through the cart to find the item id
	//If the id doesn't exist; display error
	//Otherwise; Remove the specified item from the cart
	shoppingCart[i] = shoppingCart[len(shoppingCart)-1] // Copy last element to index i.
	shoppingCart[len(shoppingCart)-1] = &Sale{}         // Erase last element (write zero value).
	shoppingCart = shoppingCart[:len(shoppingCart)-1]   // Truncate slice.
}

func GetCartTotal() float64{
	total := 0.0
	for i := 0; i < len(shoppingCart); {
		total += shoppingCart[i].price * float64(shoppingCart[i].quantity)
	}
	return total
}

func ClearCart() {
	shoppingCart = shoppingCart[:0]
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

func UpdateInventory(id int) int {
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

//These 2 might need to be changed
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
}