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


//A test Main
func TestMain() {
	UpdateData(Sale{
		id:    0,
		name:  "Null",
		price: 0,
	}, "Log", 0)

	ReadVal("Log")
}

func BuyCart(ShoppingCart []*Sale) {
	targetSheet := "Report Data"
	for i := 0; i < len(ShoppingCart); {
		UpdateData(*ShoppingCart[i], targetSheet, 1)
		i++
	}
	//Clear cart
	ClearCart(ShoppingCart)
}

//[Untested]
func AddToCart(id int, ShoppingCart []*Sale) {
	targetSheet := "Items"
	i := 0
	for {
		if i < len(ShoppingCart) {
			if ShoppingCart[i].id == id {
				ShoppingCart[i].quantity++
				break
			}
			i++
		}
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

//[Untested]
func DecreaseFromCart(id int, ShoppingCart []*Sale){
	for i := 0; i < len(ShoppingCart); {
		if ShoppingCart[i].id == id {
			if		ShoppingCart[i].quantity- 1 > 0{
				ShoppingCart[i].quantity--
			}else{
				RemoveFromCart(i, ShoppingCart)
			}
		}
		i++
	}
}

// RemoveFromCart [Untested]
func RemoveFromCart(i int, ShoppingCart []*Sale) {
	ShoppingCart[i] = ShoppingCart[len(ShoppingCart)-1] // Copy last element to index i.
	ShoppingCart[len(ShoppingCart)-1] = &Sale{}         // Erase last element (write zero value).
	ShoppingCart = ShoppingCart[:len(ShoppingCart)-1]   // Truncate slice.
}

func GetCartTotal(ShoppingCart []*Sale) float64 {
	total := 0.0
	for i := 0; i < len(ShoppingCart); {
		total += ShoppingCart[i].price * float64(ShoppingCart[i].quantity)
	}
	return total
}

func ClearCart(ShoppingCart []*Sale) {
	ShoppingCart = ShoppingCart[:0]
}

func ConvertStringToSale(price, cost, quantity string) (float64, float64, int){
	newPrice, _ := strconv.ParseFloat(price, 64)
	newCost, _ := strconv.ParseFloat(cost, 64)
	newQuantity, _ := strconv.Atoi(quantity)
	return newPrice, newCost, newQuantity
}