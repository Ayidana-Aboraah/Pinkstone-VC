package Database

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func SearchCustomers(input string) (Names []string, IDs []uint16) {
	if input == "" {
		Names = Customers
		for i := range Customers {
			IDs = append(IDs, uint16(i))
		}
		return
	}

	input = strings.TrimSpace(input)

	for i := 0; i < len(Customers); i++ {
		if strings.Contains(strings.ToLower(Customers[i]), strings.ToLower(input)) && Customers[i][0] != byte(216) {
			Names = append(Names, Customers[i])
			IDs = append(IDs, uint16(i))
		}
	}
	return
}

func AddDamages(target uint16, quantityTxt string) int {
	quantity, err := ProcessQuantity(quantityTxt)
	if err != -1 {
		return 0
	}

	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

	s := Sale{
		ID:       target,
		Price:    0,
		Cost:     Items[target].Cost[0],
		Quantity: float32(quantity),
		Usr:      255,
		Day:      uint8(day),
		Month:    uint8(month),
		Year:     uint8(year),
	}

	BuyCart([]Sale{s}, 0)
	return -1
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

func ProcessNewItemData(bargin, pieceTxt, totalTxt string, s *Sale) int {
	if pieceTxt != "" || totalTxt != "" {
		if pieceTxt == "" || totalTxt == "" {
			return 1
		} else {
			piece, err := strconv.ParseFloat(pieceTxt, 32)

			if err != nil || piece < 0 {
				return 0
			}

			total, err := strconv.ParseFloat(totalTxt, 32)
			if err != nil || total < 0 {
				return 0
			}

			s.Quantity = float32(piece / total)
		}
	}

	if bargin != "" {
		f, err := strconv.ParseFloat(bargin, 32)
		if err != nil {
			return 0
		}
		s.Price = float32(f) / s.Quantity
	}

	return -1
}

func ShiftQuantity(ID uint16) {
	Items[ID].Quantity[0] = Items[ID].Quantity[1]
	Items[ID].Quantity[1] = Items[ID].Quantity[2]
	Items[ID].Quantity[2] = 0

	if Items[ID].Cost[1] > 0 {
		Items[ID].Cost[0] = Items[ID].Cost[1]
		Items[ID].Cost[1] = Items[ID].Cost[2]
		Items[ID].Cost[2] = 0
	}
}

func AdjustStock(ID uint16) {
	for Items[ID].Quantity[0] <= 0 && Items[ID].Cost[1] > 0 {
		ShiftQuantity(ID)
	}
}

func BuyCart(ShoppingCart []Sale, customer int) []Sale {
	y, month, day := time.Now().Date()
	year, _ := strconv.Atoi(strconv.Itoa(y)[1:])

	for _, v := range ShoppingCart {

		v.Day = uint8(day)
		v.Month = uint8(month)
		v.Year = uint8(year)
		v.Customer = uint8(customer)

		for i := range Items[v.ID].Quantity {
			v.Cost = Items[v.ID].Cost[i]
			Items[v.ID].Quantity[i] -= v.Quantity

			if Items[v.ID].Quantity[i] >= 0 || i == 2 || Items[v.ID].Cost[i+1] == 0 {
				Sales = append(Sales, v)
				break
			}

			diff := Items[v.ID].Quantity[i] * -1

			v.Quantity -= diff

			Sales = append(Sales, v)

			v.Quantity = diff
		}

		AdjustStock(v.ID)
	}
	return ShoppingCart[:0]
}

func AddToCart(item Sale, ShoppingCart []Sale) (out []Sale, errID int) {
	errID = -1

	if Items[item.ID].Quantity[0] <= 0 {
		errID = 4
	}

	for i, v := range ShoppingCart {
		if v.ID == item.ID && v.Price == item.Price {
			ShoppingCart[i].Quantity += item.Quantity
			out = ShoppingCart
			return ShoppingCart, errID
		}
	}

	out = append(ShoppingCart, item)
	return out, errID
}

func DecreaseFromCart(item int, ShoppingCart []Sale) []Sale {

	if ShoppingCart[item].Quantity-1 > 0 {
		ShoppingCart[item].Quantity -= 1
	} else {
		ShoppingCart[item] = ShoppingCart[len(ShoppingCart)-1] // Copy last element to index i.
		ShoppingCart = ShoppingCart[:len(ShoppingCart)-1]      // Truncate slice.
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

func ConvertString(Price, Cost, Quantity string) (float32, float32, float32, int) {
	newPrice, errA := strconv.ParseFloat(Price, 32)
	newCost, errB := strconv.ParseFloat(Cost, 32)
	newQuantity, errC := strconv.ParseFloat(Quantity, 32)
	if errA != nil || errB != nil || (errC != nil && Quantity != "") {
		return 0.0, 0.0, 0.0, 0
	}
	return float32(newPrice), float32(newCost), float32(newQuantity), -1
}

func NewItem(id uint16) (result Sale) {
	vals := Items[id]

	result.ID = id
	result.Price = vals.Price
	result.Cost = vals.Cost[0]
	result.Quantity = 1
	result.Usr = uint8(Current_User)
	return
}
