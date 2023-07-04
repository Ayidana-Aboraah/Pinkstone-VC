package Database

import (
	"BronzeHermes/UI"
	"fmt"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

func SearchInventory(input string) (Names []string, IDs []uint16) {
	if input == "" {
		for i, v := range Items {
			IDs = append(IDs, i)
			Names = append(Names, v.Name)
		}
		return
	}

	for id, iv := range Items {
		if strings.Contains(strings.ToLower(iv.Name), strings.ToLower(input)) && iv.Name[0] != byte(216) {
			IDs = append(IDs, id)
			Names = append(Names, iv.Name)
		}
	}
	return
}

func SearchCustomers(input string) (Names []string, IDs []uint16) {
	if input == "" {
		Names = Customers
		for i := range Customers {
			IDs = append(IDs, uint16(i))
		}
		return
	}

	for i := 0; i < len(Customers); i++ {
		if strings.Contains(strings.ToLower(Customers[i]), strings.ToLower(input)) && Customers[i][0] != byte(216) {
			Names = append(Names, Customers[i])
			IDs = append(IDs, uint16(i))
		}
	}
	return
}

func CleanUpDeadItems() {
	for k := range Items {
		if Items[k].Name[0] == byte(216) {
			found := false
			for _, x := range Sales {
				if x.ID == k {
					found = true
					break
				}
			}
			if !found {
				delete(Items, k)
			}
		}
	}
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
		out += fmt.Sprintf("\n%s x%1.2f for ₵%1.2f\n", Items[v.ID].Name, v.Quantity, v.Price)
	}
	out += fmt.Sprintf("Total: %1.1f\n\n Cashier: %s\n", GetCartTotal(cart), Users[Current_User])
	out += "ALL SALES ARE FINAL\nThank you, please do come again\nSoftware Developed By Ayidana Aboraah\nTEL: +1 571-697-9347\nredstonegameraa@gmail.com\n"
	return
}

func ShiftQuantity(ID uint16) {
	Items[ID].Quantity[0] = Items[ID].Quantity[1]
	Items[ID].Quantity[1] = Items[ID].Quantity[2]
	Items[ID].Quantity[2] = 0

	Items[ID].Cost[0] = Items[ID].Cost[1]
	Items[ID].Cost[1] = Items[ID].Cost[2]
	Items[ID].Cost[2] = 0
}

func BuyCart(ShoppingCart []Sale, customer int) []Sale {
	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])
	for _, v := range ShoppingCart {

		v.Day = uint8(day)
		v.Month = uint8(month)
		v.Year = uint8(year)
		v.Customer = uint8(customer)

		if Items[v.ID].Quantity[0]-v.Quantity <= 0 {
			newbie := v
			newbie.Cost = Items[v.ID].Cost[1]
			newbie.Quantity = (Items[v.ID].Quantity[0] - v.Quantity) * -1

			if newbie.Quantity != 0 {
				v.Quantity -= newbie.Quantity

				if Items[v.ID].Quantity[1]-newbie.Quantity < 0 {
					newbie2 := newbie
					newbie2.Cost = Items[v.ID].Cost[2]
					newbie2.Quantity = (Items[v.ID].Quantity[1] - newbie.Quantity) * -1

					newbie.Quantity -= newbie2.Quantity

					Items[v.ID].Quantity[1] -= newbie.Quantity
					Items[v.ID].Quantity[2] -= newbie2.Quantity

					Sales = append(Sales, newbie2)
				}

				Sales = append(Sales, newbie)
			}

			ShiftQuantity(v.ID)
			if Items[v.ID].Quantity[0] == 0 {
				ShiftQuantity(v.ID)
			}
		} else {
			Items[v.ID].Quantity[0] -= v.Quantity
		}

		Sales = append(Sales, v)
	}
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
	// var pop []int
	for k := range Items {
		if Items[k].Name[0] != byte(216) {
			inter = append(inter, int(k))
		}
	}
	// for inv := len(pop); inv > 0; inv-- {
	// 	for i := 0; i < len(pop); i++ { // || I > val[i]

	// 	}
	// }

	return
}

func ProcessQuantity(n string, w fyne.Window) (quantity float32) {
	raw := strings.SplitN(n, " ", 2)
	if len(raw) == 2 {
		pop := strings.SplitN(raw[1], "/", 2)
		numerator, denominator, whole := ConvertString(pop[0], pop[1], raw[0])
		if UI.HandleKnownError(0, len(pop) != 2, w) {
			return
		}
		quantity = whole + (numerator / denominator)
	} else {
		v, err := strconv.ParseFloat(raw[0], 32)
		if UI.HandleKnownError(0, err != nil, w) {
			return
		}
		quantity = float32(v)
	}
	return
}

func RemoveReportEntry(index int) {

	for i, v := range Items[Sales[index].ID].Cost {
		if v == 0 || v == Sales[index].Cost {
			Items[Sales[index].ID].Cost[i] = Sales[index].Cost
			Items[Sales[index].ID].Quantity[i] += Sales[index].Quantity
			break
		}
	}

	Sales[index] = Sales[len(Sales)-1]
	Sales = Sales[:len(Sales)-1]
}

func ConvertString(Price, Cost, Quantity string) (float32, float32, float32) {
	newPrice, _ := strconv.ParseFloat(Price, 32)
	newCost, _ := strconv.ParseFloat(Cost, 32)
	newQuantity, _ := strconv.ParseFloat(Quantity, 32)
	return float32(newPrice), float32(newCost), float32(newQuantity)
}

func ConvertItem(id uint16) (result Sale) {
	vals := Items[id]

	result.ID = id
	result.Price = vals.Price
	result.Cost = vals.Cost[0]
	result.Quantity = 1
	result.Usr = uint8(Current_User)
	return
}
