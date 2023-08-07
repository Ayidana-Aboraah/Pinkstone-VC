package Database

import (
	"BronzeHermes/Debug"
	unknown "BronzeHermes/Unknown"
)

func CreateSale(ID uint16, dateStr, priceTxt, costTxt, stockTxt string, customer int) int {
	date, err := unknown.ProcessDate(dateStr)
	if err != Debug.Success {
		return err
	}

	quan, err := ProcessQuantity(stockTxt)
	if err != Debug.Success {
		return err
	}

	price, cost, _, err := ConvertString(priceTxt, costTxt, "")
	if err != Debug.Success {
		return err
	}

	Sales = append(Sales, Sale{
		ID:       ID,
		Price:    price,
		Cost:     cost,
		Quantity: quan,
		Customer: uint8(customer),
		Usr:      Current_User,
		Day:      date[2],
		Month:    date[1],
		Year:     date[0],
	})

	return Debug.Success
}

func DeleteEverything() {
	Items = map[uint16]*Entry{}
	Sales = []Sale{}
	Current_User = 0
	Users = []string{}
	Customers = []string{}
}
