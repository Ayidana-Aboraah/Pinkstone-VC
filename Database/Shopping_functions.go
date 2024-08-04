package Database

import (
	"BronzeHermes/Debug"
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
	if err != Debug.Success {
		return Debug.Invalid_Input
	}

	s := Sale{
		ID:        target,
		Price:     0,
		Cost:      Items[target].Cost[0],
		Quantity:  float32(quantity),
		Timestamp: time.Now().Local().Unix(),
	}

	BuyCart([]Sale{s}, 0)
	return Debug.Success
}

func MakeReceipt(cart []Sale, customer string) (out string) {

	out = time.Now().Local().String()
	out += "\nLoc: Santasi\nTel/Vodacash: 0506695927\nTel/MOMO: 0558324302\nMerchant ID: 868954\nCustomer: " + customer + "\n"
	// TODO: Maybe save the Voda & MOMO, along with the merchant ID as variable to be accessed by the controller
	for _, v := range cart {
		out += fmt.Sprintf("\n%s x%1.2f for ₵%1.2f\n", Items[v.ID].Name, v.Quantity, v.Price)
	}

	out += fmt.Sprintf("Total: %1.1f\n", GetCartTotal(cart))
	out += "ALL SALES ARE FINAL\nThank you, please do come again\nSoftware Developed By Ayidana Aboraah\naboraahayidana@gmail.com\n"

	return
}

func ProcessNewItemData(bargin, stockTxt string, s *Sale) int {
	if stockTxt != "" {
		quantity, errID := ProcessQuantity(stockTxt)
		if errID != Debug.Success {
			return errID
		}
		s.Quantity = quantity
	}

	if bargin != "" {
		f, err := strconv.ParseFloat(bargin, 32)
		if err != nil {
			return Debug.Invalid_Input
		}
		s.Price = float32(f) / s.Quantity
		s.Cost *= s.Quantity
	}

	return Debug.Success
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

func BuyCart(ShoppingCart []Sale, customer uint16) []Sale {
	for _, v := range ShoppingCart {

		v.Timestamp = time.Now().Local().Unix()
		v.Customer = customer

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
	errID = Debug.Success

	if Items[item.ID].Quantity[0] <= 0 {
		errID = Debug.Empty_Quantity_Warning
	}

	for i, v := range ShoppingCart {
		if v.ID == item.ID && v.Price == item.Price {
			ShoppingCart[i].Quantity += item.Quantity
			return ShoppingCart, errID
		}
	}

	return append(ShoppingCart, item), errID
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
		intercart = append(intercart, shoppingCart[i])
	}
	return
}

func ConvertString(Price, Cost, Quantity string) (price float32, cost float32, quantity float32, errID int) {
	p, errA := strconv.ParseFloat(Price, 32)
	c, errB := strconv.ParseFloat(Cost, 32)
	q, errC := strconv.ParseFloat(Quantity, 32)
	errID = Debug.Success

	if errA != nil || errB != nil || (errC != nil && Quantity != "") {
		errID = Debug.Invalid_Input
	}
	return float32(p), float32(c), float32(q), errID
}

func NewItem(id uint16) (result Sale) {
	vals := Items[id]

	result.ID = id
	result.Price = vals.Price
	result.Cost = vals.Cost[0]
	result.Quantity = 1
	return
}
