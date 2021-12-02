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

func NewSale(ID int, Name string, Price, Cost float64, Quantity int) Sale{
	return Sale{
		ID: ID,
		Name: Name,
		Price: Price,
		Cost: Cost,
		Quantity: Quantity,
	}
}

func BuyCart(ShoppingCart []Sale) []Sale{
	targetSheet := "Report Data"

	for _, v := range ShoppingCart {
		UpdateData(v, targetSheet, 1)
	}

	ReadVal(targetSheet)
	SaveFile()
	return ClearCart(ShoppingCart)
}

//Must pass as the new value of Shopping Cart similar to appending to an array
func AddToCart(ID int, ShoppingCart []Sale) []Sale{
	targetSheet := "Items"
	for {
		for i, v := range ShoppingCart{
			if v.ID == ID {
				//v.Quantity += 1
				ShoppingCart[i].Quantity += 1
				return ShoppingCart
			}
		}
		idx := GetIndex(targetSheet, ID, 1)
		p, _ := strconv.ParseFloat(f.GetCellValue(targetSheet, "B"+strconv.Itoa(idx)), 64)
		c, _ := strconv.ParseFloat(f.GetCellValue(targetSheet, "C"+strconv.Itoa(idx)), 64)

		temp := Sale{
			ID:    ID,
			Name:  f.GetCellValue(targetSheet, "B"+strconv.Itoa(idx)),
			Price: p,
			Cost: c,
			Quantity: 1,
		}
		fmt.Println(temp)
		ShoppingCart = append(ShoppingCart, temp)
		return ShoppingCart
	}
}

func DecreaseFromCart(ID int, ShoppingCart []Sale) []Sale{
	for i, v := range ShoppingCart {
		if v.ID == ID {
			if v.Quantity-1 > 0 {
				ShoppingCart[i].Quantity -= 1
			} else {
				ShoppingCart = RemoveFromCart(i, ShoppingCart)
			}
		}
	}

	return ShoppingCart
}

// RemoveFromCart [Untested]
func RemoveFromCart(i int, ShoppingCart []Sale) []Sale{
	ShoppingCart[i] = ShoppingCart[len(ShoppingCart)-1] // Copy last element to index i.
	ShoppingCart[len(ShoppingCart)-1] = Sale{}         // Erase last element (write zero value).
	ShoppingCart = ShoppingCart[:len(ShoppingCart)-1]   // Truncate slice.
	return ShoppingCart
}

func GetCartTotal(ShoppingCart []Sale) float64 {
	total := 0.0
	for _, v := range ShoppingCart{
		total += v.Price * float64(v.Quantity)
	}
	return total
}

//Removes all items from the shopping cart
func ClearCart(ShoppingCart []Sale) []Sale {
	ShoppingCart = ShoppingCart[:0]
	return ShoppingCart
}

func ConvertStringToSale(Price, Cost, Quantity string) (float64, float64, int){
	newPrice, _ := strconv.ParseFloat(Price, 64)
	newCost, _ := strconv.ParseFloat(Cost, 64)
	newQuantity, _ := strconv.Atoi(Quantity)
	return newPrice, newCost, newQuantity
}

func ConvertSaleToString(price, cost float64, inventory int) []string{
	p := fmt.Sprint(price)
	c := fmt.Sprint(cost)
	inven := strconv.Itoa(inventory)

	return []string{
		p,
		c,
		inven,
	}
}