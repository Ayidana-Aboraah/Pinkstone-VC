package Data

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

// Defining the info used later
type Sale struct {
	id       int
	name     string
	price    float64
	cost 	 float64
	quantity int
}

func NewSale(id int, name string, price, cost float64, quantity int) Sale{
	return Sale{
		id: id,
		name: name,
		price: price,
		cost: cost,
		quantity: quantity,
	}
}

//var f, err = excelize.OpenFile("AppData.xlsx")
var f, err = excelize.OpenFile("TestAppData.xlsx")

var shoppingCart []*Sale

//A test Main
func TestMain() {
	UpdateData(Sale{
		id:    0,
		name:  "Null",
		price: 0,
	}, "Log", 0)

	ReadVal("Log")
	//_ = f.Save()
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
	//Gets the index of a free slot in the sheet
	//idx := GetIndex(targetSheet, 0, 0)

	//Loops through to fill each item on the cart to the stuff
	for i := 0; i <= len(shoppingCart); {
		UpdateData(*shoppingCart[i], targetSheet, 1)
		/*
		newInven := GetInventory(shoppingCart[i].id) - 1
		quantity := float64(shoppingCart[i].quantity)
		f.SetCellValue(targetSheet, "A"+strconv.Itoa(idx), shoppingCart[i].id)
		f.SetCellValue(targetSheet, "B"+strconv.Itoa(idx), shoppingCart[i].name)
		f.SetCellValue(targetSheet, "C"+strconv.Itoa(idx), quantity)
		f.SetCellValue(targetSheet, "D"+strconv.Itoa(idx), shoppingCart[i].price*quantity)
		f.SetCellValue(targetSheet, "E"+strconv.Itoa(idx), shoppingCart[i].cost*quantity)
		f.SetCellValue(targetSheet, "F"+strconv.Itoa(idx), newInven)
		f.SetCellValue(targetSheet, "G"+strconv.Itoa(idx), ConvertDate(time.Now()))
		f.SetCellValue(targetSheet, "H"+strconv.Itoa(idx), ConvertClock(time.Now()))
		UpdateLog(*shoppingCart[i], "Log")
		 */
		i++
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
	i := 0
	for {
		if i < len(shoppingCart) {
			if shoppingCart[i].id == id {
				shoppingCart[i].quantity++
				break
			}
			i++
		} else {
			targetSheet := "Detection Data"
			idx := GetIndex(targetSheet, id, 1)
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

func ConvertStringToSale(price, cost, quantity string) (float64, float64, int){
	newPrice, _ := strconv.ParseFloat(price, 64)
	newCost, _ := strconv.ParseFloat(cost, 64)
	newQuantity, _ := strconv.Atoi(quantity)
	return newPrice, newCost, newQuantity
}