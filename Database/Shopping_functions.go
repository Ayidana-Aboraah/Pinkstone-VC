package Database

import (
	"BronzeHermes/Helpers"
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

func AddDamages(target uint, quantityTxt string) int {
	quantity, err := ProcessQuantity(quantityTxt)
	if err != Helpers.Success {
		return Helpers.Invalid_Input
	}

	s := Sale{
		ID:        target,
		Price:     0,
		Cost:      Items[target].Cost[0],
		Quantity:  quantity,
		Timestamp: time.Now().Local().Unix(),
	}

	Buy([]Sale{s}, 0)
	return Helpers.Success
}

var VodaPhoneNumber string
var MOMO_Num string

// var MerchantID string

const recieptTemplate string = "%s\n" + "\nLoc: Santasi\nTel/Vodacash: %s\nTel/MOMO: %s\nMerchant ID: 868954\nCustomer: %s\n" + "\n%s" + "Total: %1.1f\n" + "ALL SALES ARE FINAL\nThank you, please do come again\nSoftware Developed By Ayidana Aboraah\naboraahayidana@gmail.com\n"

func MakeReceipt(cart []Sale, customer string) string {

	var sales string

	for _, v := range cart {
		sales += fmt.Sprintf("\n%s x%1.2f for ₵%1.2f\n", Items[v.ID].Name, v.Quantity, v.Price)
	}

	return fmt.Sprintf(recieptTemplate, time.Now().Local().String(), customer, VodaPhoneNumber, MOMO_Num, sales, Total(cart))
}

func ProcessNewItemData(bargin, stockTxt string, id uint) (s Sale, err int) {
	if stockTxt != "" {
		quantity, errID := ProcessQuantity(stockTxt)
		if errID != Helpers.Success {
			return s, errID
		}
		s.Quantity = quantity
	}

	if bargin != "" {
		f, err := strconv.ParseFloat(bargin, 32)
		if err != nil {
			return s, Helpers.Invalid_Input
		}
		s.Price = f / s.Quantity
		s.Cost *= s.Quantity
	}

	return s, Helpers.Success
}

func Adjust(ID uint) {
	if Items[ID].Quantity[0] <= 0 && len(Items[ID].Quantity) > 1 {
		Items[ID].Quantity = Items[ID].Quantity[1:]
		Items[ID].Cost = Items[ID].Cost[1:]
	}
}

func Buy(ShoppingCart []Sale, customer uint) []Sale {
	for _, sale := range ShoppingCart {

		sale.Timestamp = time.Now().Local().Unix()
		sale.Customer = customer
		sale.Cost = Items[sale.ID].Cost[0]

		for Items[sale.ID].Quantity[0]-sale.Quantity <= 0 && len(Items[sale.ID].Quantity) > 1 {
			sCopy := sale

			sCopy.Quantity = Items[sale.ID].Quantity[0]
			Sales = append(Sales, sCopy)

			Adjust(sale.ID)

			sale.Quantity -= Items[sale.ID].Quantity[0]
			sale.Cost = Items[sale.ID].Cost[0]
		}

		Sales = append(Sales, sale)
		Items[sale.ID].Quantity[0] -= sale.Quantity
	}
	return ShoppingCart[:0]
}

func AddToCart(item Sale, ShoppingCart []Sale) (out []Sale, errID int) {
	errID = Helpers.Success

	if Items[item.ID].Quantity[0] <= 0 {
		errID = Helpers.Empty_Quantity_Warning
	}

	for i, cartItem := range ShoppingCart {
		if cartItem.ID == item.ID && cartItem.Price == item.Price {
			ShoppingCart[i].Quantity += item.Quantity
			return ShoppingCart, errID
		}
	}

	return append(ShoppingCart, item), errID
}

func DecreaseCart(item int, ShoppingCart []Sale) []Sale {

	if ShoppingCart[item].Quantity-1 > 0 {
		ShoppingCart[item].Quantity -= 1
	} else {
		ShoppingCart[item] = ShoppingCart[len(ShoppingCart)-1] // Copy last element to index i.
		ShoppingCart = ShoppingCart[:len(ShoppingCart)-1]      // Truncate slice.
	}

	return ShoppingCart
}

func Total(ShoppingCart []Sale) (total float64) {
	for _, v := range ShoppingCart {
		total += v.Price * v.Quantity
	}
	return
}

func Convert(shoppingCart []Sale) (intercart []interface{}) {
	for i := range shoppingCart {
		intercart = append(intercart, shoppingCart[i])
	}
	return
}

func SaleItem(id uint) (result Sale) {
	vals := Items[id]

	result.ID = id
	result.Price = vals.Price
	result.Cost = vals.Cost[0]
	result.Quantity = 1
	return
}
