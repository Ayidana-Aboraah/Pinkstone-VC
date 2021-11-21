package Data

import (
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
var f, _ = excelize.OpenFile("TestAppData.xlsx")

var ShoppingCart []*Sale

//A test Main
func TestMain() {
	UpdateData(Sale{
		id:    0,
		name:  "Null",
		price: 0,
	}, "Log", 0)

	ReadVal("Log")
}

func BuyCart() {
	//Cycle through and check if there is an open space
	//When finding an empty space start unloading the contents of the cart into each cell
	//Until the cart's info is fully displayed

	targetSheet := "Report Data"

	//Loops through to fill each item on the cart to the stuff
	for i := 0; i <= len(ShoppingCart); {
		UpdateData(*ShoppingCart[i], targetSheet, 1)
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
		if i < len(ShoppingCart) {
			if ShoppingCart[i].id == id {
				ShoppingCart[i].quantity++
				break
			}
			i++
		} else {
			targetSheet := "Detection Data"
			idx := GetIndex(targetSheet, id, 1)
			p, _ := strconv.ParseFloat(f.GetCellValue(targetSheet, "B"+strconv.Itoa(idx)), 64)

			temp := Sale{
				id:    id,
				name:  f.GetCellValue(targetSheet, "B"+strconv.Itoa(idx)),
				price: p,
			}

			ShoppingCart = append(ShoppingCart, &temp)
			break
		}
	}
}

//[Untested]
func DecreaseFromCart(id int){
	for i := 0; i < len(ShoppingCart); {
		if ShoppingCart[i].id == id {
			if		ShoppingCart[i].quantity- 1 > 0{
				ShoppingCart[i].quantity--
			}else{
				RemoveFromCart(i)
			}
		}
		i++
	}
}

// RemoveFromCart [Untested]
func RemoveFromCart(i int) {
	//Cycle through the cart to find the item id
	//If the id doesn't exist; display error
	//Otherwise; Remove the specified item from the cart

	ShoppingCart[i] = ShoppingCart[len(ShoppingCart)-1] // Copy last element to index i.
	ShoppingCart[len(ShoppingCart)-1] = &Sale{}         // Erase last element (write zero value).
	ShoppingCart = ShoppingCart[:len(ShoppingCart)-1]   // Truncate slice.
}

func GetCartTotal() float64{
	total := 0.0
	for i := 0; i < len(ShoppingCart); {
		total += ShoppingCart[i].price * float64(ShoppingCart[i].quantity)
	}
	return total
}

func ClearCart() {
	ShoppingCart = ShoppingCart[:0]
}

func ConvertStringToSale(price, cost, quantity string) (float64, float64, int){
	newPrice, _ := strconv.ParseFloat(price, 64)
	newCost, _ := strconv.ParseFloat(cost, 64)
	newQuantity, _ := strconv.Atoi(quantity)
	return newPrice, newCost, newQuantity
}