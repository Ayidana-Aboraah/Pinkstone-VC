package Database

import (
	"fmt"
	"strconv"
)

func BuyCart(ShoppingCart []Sale) []Sale {
	Reports[0] = append(Reports[0], ShoppingCart...)
	SaveData()
	return ShoppingCart[:0]
}

func AddToCart(item Sale, ShoppingCart []Sale) []Sale {
	for i, v := range ShoppingCart {
		if v.ID != item.ID || v.Price != item.Price {
			continue
		}

		ShoppingCart[i].Quantity++
		return ShoppingCart
	}
	ShoppingCart = append(ShoppingCart, item)
	return ShoppingCart
}

func DecreaseFromCart(item Sale, ShoppingCart []Sale) []Sale {
	for i, v := range ShoppingCart {
		if v.ID != item.ID || v.Price != item.Price {
			continue
		}

		if v.Quantity-1 > 0 {
			ShoppingCart[i].Quantity -= 1
		} else {
			ShoppingCart[i] = ShoppingCart[len(ShoppingCart)-1] // Copy last element to index i.
			ShoppingCart = ShoppingCart[:len(ShoppingCart)-1]   // Truncate slice.
		}
	}

	return ShoppingCart
}

func GetCartTotal(ShoppingCart []Sale) (total float32) {
	for _, v := range ShoppingCart {
		total += v.Price * float32(v.Quantity)
	}
	return
}

func ConvertCart(shoppingCart []Sale) (intercart []interface{}) {
	for i := range shoppingCart {
		intercart = append(intercart, shoppingCart[i])
	}
	return
}

func ConvertItems() (inter []int) {
	for k := range ItemKeys {
		inter = append(inter, int(k))
	}
	return
}

func ConvertExpenses() (inter []interface{}) {
	for i := range Expenses {
		inter = append(inter, Expenses[i])
	}
	return
}

func RemoveExpense(index int) {
	Expenses[index] = Expenses[len(Expenses)-1]
	Expenses = Expenses[:len(Expenses)-1]
}

func ConvertString(Price, Cost, Quantity string) (float32, float32, uint16) {
	newPrice, _ := strconv.ParseFloat(Price, 64)
	newCost, _ := strconv.ParseFloat(Cost, 64)
	newQuantity, _ := strconv.Atoi(Quantity)
	return float32(newPrice), float32(newCost), uint16(newQuantity)
}

func ConvertSale(item Sale) []string {
	p := fmt.Sprint(item.Price)
	c := fmt.Sprint(item.Cost)
	inven := strconv.Itoa(int(item.Quantity))

	return []string{
		p,
		c,
		inven,
	}
}
