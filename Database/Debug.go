package Database

import (
	"BronzeHermes/Debug"
	"time"
)

func CreateSale(ID uint16, priceTxt, costTxt, stockTxt string, customer int) int {

	quan, err := ProcessQuantity(stockTxt)
	if err != Debug.Success {
		return err
	}

	price, cost, _, err := ConvertString(priceTxt, costTxt, "")
	if err != Debug.Success {
		return err
	}

	Sales = append(Sales, Sale{
		ID:        ID,
		Price:     price,
		Cost:      cost,
		Quantity:  quan,
		Customer:  uint16(customer),
		Timestamp: time.Now().Local().Unix(),
	})

	return Debug.Success
}

func DeleteEverything() {
	Items = []Item{}
	Sales = []Sale{}
	Customers = []string{}
}
