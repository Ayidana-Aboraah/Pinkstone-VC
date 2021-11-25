package Data

import (
	"fmt"
	"strconv"
)

// Defining the Sale info used later
type Sale struct {
	ID       int
	Name     string
	Price    float64
	Cost 	 float64
	Quantity int
}

//A constructor for a NewSale
func NewSale(ID int, Name string, Price, Cost float64, Quantity int) Sale{
	return Sale{
		ID: ID,
		Name: Name,
		Price: Price,
		Cost: Cost,
		Quantity: Quantity,
	}
}

//A test Main
func TestMain() {
	UpdateData(Sale{
		ID:    0,
		Name:  "Null",
		Price: 0,
	}, "Log", 0)

	ReadVal("Log")
}

func BuyCart(ShoppingCart []*Sale) {
	targetSheet := "Report Data"

	/*
	for i := 0; i < len(ShoppingCart); {
		UpdateData(*ShoppingCart[i], targetSheet, 1)
		i++
	}
	 */

	for _, v := range ShoppingCart {
		UpdateData(*v, targetSheet, 1)
	}
	//Clear cart
	ClearCart(ShoppingCart)
}

//Must pass as the new value of Shopping Cart similar to appending to an array
func AddToCart(ID int, ShoppingCart []*Sale) []*Sale{
	targetSheet := "Items"
	//i := 0
	for {
		/*
		if i < len(ShoppingCart) {
			if ShoppingCart[i].ID == ID {
				ShoppingCart[i].Quantity++
				break
			}
			i++
		}
		 */

		for _, v := range ShoppingCart{
			if v.ID == ID {
				v.Quantity++
				break
			}
		}
		idx := GetIndex(targetSheet, ID, 1)
		p, _ := strconv.ParseFloat(f.GetCellValue(targetSheet, "B"+strconv.Itoa(idx)), 64)

		temp := &Sale{
			ID:    ID,
			Name:  f.GetCellValue(targetSheet, "B"+strconv.Itoa(idx)),
			Price: p,
		}
		fmt.Println(temp)
		ShoppingCart = append(ShoppingCart, temp)
		return ShoppingCart
	}
	return ShoppingCart
}

//[Untested]
func DecreaseFromCart(ID int, ShoppingCart []*Sale) {
	/*
	for i := 0; i < len(ShoppingCart); {
		if ShoppingCart[i].ID == ID {
			if ShoppingCart[i].Quantity-1 > 0 {
				ShoppingCart[i].Quantity--
			} else {
				RemoveFromCart(i, ShoppingCart)
			}
		}
		i++
	}
	 */

	for i, v := range ShoppingCart {
		if v.ID == ID {
			if v.Quantity-1 > 0 {
				v.Quantity--
			} else {
				RemoveFromCart(i, ShoppingCart)
			}
		}
	}
}

// RemoveFromCart [Untested]
func RemoveFromCart(i int, ShoppingCart []*Sale) {
	ShoppingCart[i] = ShoppingCart[len(ShoppingCart)-1] // Copy last element to index i.
	ShoppingCart[len(ShoppingCart)-1] = &Sale{}         // Erase last element (write zero value).
	ShoppingCart = ShoppingCart[:len(ShoppingCart)-1]   // Truncate slice.
}

func GetCartTotal(ShoppingCart []*Sale) (float64 , string) {
	total := 0.0
	for _, v := range ShoppingCart{
		total += v.Price * float64(v.Quantity)
	}
	strTotal := fmt.Sprint(total)
	return total, strTotal
}

//Removes all items from the shopping cart
func ClearCart(ShoppingCart []*Sale) []*Sale {
	return ShoppingCart[:0]
}

func ConvertStringToSale(Price, Cost, Quantity string) (float64, float64, int){
	newPrice, _ := strconv.ParseFloat(Price, 64)
	newCost, _ := strconv.ParseFloat(Cost, 64)
	newQuantity, _ := strconv.Atoi(Quantity)
	return newPrice, newCost, newQuantity
}