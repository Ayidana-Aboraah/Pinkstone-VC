package Database

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func SearchInventory(input string) (Names []string, IDs []uint16) {
	for id, iv := range Item {
		if strings.Contains(strings.ToLower(iv.Name), strings.ToLower(input)) {
			IDs = append(IDs, id)
			Names = append(Names, iv.Name)
		}
	}
	return
}

func FilterUsers() (out []string) {
	for _, s := range Users {
		if s[0] != byte(216) {
			out = append(out, s)
		}
	}
	return
}

func MakeReceipt(cart []Sale, customer string) (out string) {
	y, m, d := time.Now().Date()
	hr, min, _ := time.Now().Clock()
	out = fmt.Sprintf("%d/%d/%d , %d:%2d\n", y, m, d, hr, min)
	out += "Loc: Santasi\nTel/Vodacash: 0506695927\nTel/MOMO: 0558324302\nMerchant ID: 868954\nCustomer: " + customer + "\n"

	for _, v := range cart {
		out += fmt.Sprintf("\n%s x%1.2f for ₵%1.2f\n", Item[v.ID].Name, v.Quantity, v.Price)
	}
	out += fmt.Sprintf("Total: %1.1f\n\n Cashier: %s\n", GetCartTotal(cart), Users[Current_User])
	out += "ALL SALES ARE FINAL\nThank you, please do come again\nSoftware Developed By Ayidana Aboraah\nTEL: +1 571-697-9347\nredstonegameraa@gmail.com\n"
	return
}

func ShiftQuantity(ID uint16) {

	Item[ID].Quantity[0] = Item[ID].Quantity[1]
	Item[ID].Quantity[1] = Item[ID].Quantity[2]
	Item[ID].Quantity[2] = 0

	Item[ID].Cost[0] = Item[ID].Cost[1]
	Item[ID].Cost[1] = Item[ID].Cost[2]
	Item[ID].Cost[2] = 0
}

func BuyCart(ShoppingCart []Sale) []Sale {
	for _, v := range ShoppingCart {
		y, month, day := time.Now().Date()
		year, _ := strconv.Atoi(strconv.Itoa(y)[1:])
		v.Day = uint8(day)
		v.Month = uint8(month)
		v.Year = uint8(year)

		if Item[v.ID].Quantity[0]-v.Quantity <= 0 {
			newbie := v
			newbie.Cost = Item[v.ID].Cost[1]
			newbie.Quantity = Item[v.ID].Quantity[0] - v.Quantity
			v.Quantity += newbie.Quantity

			Item[v.ID].Quantity[1] += newbie.Quantity

			Reports[0] = append(Reports[0], newbie)

			ShiftQuantity(v.ID)
			if Item[v.ID].Quantity[0] == 0 {
				ShiftQuantity(v.ID)
			}
		} else {
			Item[v.ID].Quantity[0] -= v.Quantity
		}

		Reports[0] = append(Reports[0], v)
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
		if shoppingCart[i].Usr != 255 {
			intercart = append(intercart, shoppingCart[i])
		}
	}
	return
}

func ConvertItemKeys() (inter []int) {
	for k := range Item {
		if Item[k].Name[0] != byte(216) {
			inter = append(inter, int(k))
		}
	}

	return
}

func RemoveReportEntry(report, index int) {
	Reports[report][index] = Reports[report][len(Reports[report])-1]
	Reports[report] = Reports[report][:len(Reports[report])-1]
}

func ConvertString(Price, Cost, Quantity string) (float32, float32, float32) {
	newPrice, _ := strconv.ParseFloat(Price, 32)
	newCost, _ := strconv.ParseFloat(Cost, 32)
	newQuantity, _ := strconv.ParseFloat(Quantity, 32)
	return float32(newPrice), float32(newCost), float32(newQuantity)
}

func ConvertItem(id uint16) (result Sale) {
	vals := Item[id]

	result.ID = id
	result.Price = vals.Price
	result.Cost = vals.Cost[0]
	result.Quantity = 1
	result.Usr = uint8(Current_User)
	return
}
