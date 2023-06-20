package Database

import (
	"fmt"
	"strconv"
	"time"
)

func MakeReceipt(cart []Sale) (out string) {

	y, m, d := time.Now().Date()
	out = fmt.Sprintf("\t Pinkstone Ltd. : %d/%d/%d\n", y, m, d)

	for _, v := range cart {
		out += fmt.Sprintf("%s x%d for ₵%1.1f\n", ItemKeys[v.ID].Name, v.Quantity, v.Price)
	}
	out += fmt.Sprintf("Total: %1.1f\n\n Cashier: "+Users[Current_User], GetCartTotal(cart))
	return
}

func BuyCart(ShoppingCart []Sale) []Sale {
	for _, v := range ShoppingCart {
		y, month, day := time.Now().Date()
		year, _ := strconv.Atoi(strconv.Itoa(y)[1:])
		v.Day = uint8(day)
		v.Month = uint8(month)
		v.Year = uint8(year)
		Reports[0] = append(Reports[0], v)
		ItemKeys[v.ID].Quantity -= v.Quantity
	}
	SaveData()
	return ShoppingCart[:0]
}

func AddToCart(item Sale, ShoppingCart []Sale) []Sale {
	for i, v := range ShoppingCart {
		if v.ID == item.ID && v.Price == item.Price {
			ShoppingCart[i].Quantity++
			return ShoppingCart
		}
	}
	return append(ShoppingCart, item)
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

func ConvertItemKeys() (inter []int) {
	for k := range ItemKeys {
		inter = append(inter, int(k))
	}

	return
}

// func ConvertItemIdxes(target uint64) (list []int) {
// 	list = append(list, ItemKeys[target].Idxes...)
// 	return
// }

// func RemoveItem(idx int, id uint64) {
// 	Free_Spaces = append(Free_Spaces, ItemKeys[id].Idxes[idx])
// 	ItemKeys[id].Idxes[idx] = ItemKeys[id].Idxes[len(ItemKeys[id].Idxes)-1]
// 	ItemKeys[id].Idxes = ItemKeys[id].Idxes[:len(ItemKeys[id].Idxes)-1]
// }

func ConvertExpenses() (inter []interface{}) {
	for i := range Expenses {
		inter = append(inter, Expenses[i])
	}
	return
}

func RemoveReportEntry(report, index int) {
	Reports[report][index] = Reports[report][len(Reports[report])-1]
	Reports[report] = Reports[report][:len(Reports[report])-1]
}

func RemoveExpense(index int) {
	Expenses[index] = Expenses[len(Expenses)-1]
	Expenses = Expenses[:len(Expenses)-1]
}

func ConvertString(Price, Quantity string) (float32, uint16) {
	newPrice, _ := strconv.ParseFloat(Price, 64)
	newQuantity, _ := strconv.Atoi(Quantity)
	return float32(newPrice), uint16(newQuantity)
}

func ConvertItem(id uint64) (result Sale) {
	vals := ItemKeys[id]

	result.ID = id
	result.Price = vals.Price
	result.Quantity = 1
	result.Usr = uint8(Current_User)
	return
}
