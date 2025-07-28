package Database

import (
	"BronzeHermes/Helpers"
	"time"
)

func CreateSale(ID uint, priceTxt, costTxt, stockTxt string, customer uint) int {

	quan, err := ProcessQuantity(stockTxt)
	if err != Helpers.Success {
		return err
	}

	result := Helpers.ParseToFloat(priceTxt, costTxt)
	if result == nil {
		return Helpers.Invalid_Input
	}

	Sales = append(Sales, Sale{
		ID:        ID,
		Price:     result[0],
		Cost:      result[1],
		Quantity:  quan,
		Customer:  customer,
		Timestamp: time.Now().Local().Unix(),
	})

	return Helpers.Success
}

func DeleteEverything() {
	Items = []Entry{}
	Sales = []Sale{}
	Customers = []string{}
}
