package Database

import (
	"strconv"
	"strings"
)

func CreateSale(ID uint16, Date, priceTxt, costTxt, stockTxt string, customer int) int {
	raw := strings.SplitN(Date, "-", 3)
	if len(raw) < 3 {
		return 3
	}

	year, e := strconv.Atoi(raw[0][1:])
	if e != nil {
		return 0
	}

	month, e := strconv.Atoi(raw[1])
	if e != nil {
		return 0
	}

	day, e := strconv.Atoi(raw[2])
	if e != nil {
		return 0
	}

	quan, err := ProcessQuantity(stockTxt)
	if err != -1 {
		return err
	}

	price, cost, _, err := ConvertString(priceTxt, costTxt, "")
	if err != -1 {
		return err
	}

	Sales = append(Sales, Sale{
		ID:       ID,
		Price:    price,
		Cost:     cost,
		Quantity: quan,
		Customer: uint8(customer),
		Usr:      Current_User,
		Day:      uint8(day),
		Month:    uint8(month),
		Year:     uint8(year),
	})

	return -1
}

func DeleteEverything() {
	Items = map[uint16]*Entry{}
	Sales = []Sale{}
	Current_User = 0
	Users = []string{}
	Customers = []string{}
}
