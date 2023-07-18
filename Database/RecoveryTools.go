package Database

import (
	"log"
	"strconv"
	"strings"
)

func RecoverTxt(n string, kv map[string]uint16) {
	CurrentDay := 1

	places := strings.Split(n, "\n")

	for _, entry := range places {
		if entry[0] == ' ' {
			v, err := strconv.Atoi(entry[1:])
			if err != nil {
				log.Println(entry, err)
			}
			CurrentDay = v
		}
		pops := strings.Split(entry, ",")
		if len(pops) != 4 {
			log.Println(len(pops), " | ", pops)
		}

		// use the map to get the ID to set for the new sale
		// convert the second value to be the price
		// convert the thrid to be cost
		quantity, err := ProcessQuantity(pops[3])
		if err != -1 {
			log.Println(err)
		}
		price, cost, _, err := ConvertString(pops[1], pops[2], "")
		if err != -1 {
			log.Println(err)
		}

		Sales = append(Sales, Sale{
			ID:       kv[pops[0]],
			Price:    price,
			Cost:     cost,
			Quantity: quantity,
			Day:      uint8(CurrentDay),
			Month:    7,
			Year:     23,
			Usr:      Current_User,
			Customer: 0,
		})
	}

	log.Println(Sales)
}
