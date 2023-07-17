package Database

import (
	"fmt"
	"strconv"
	"strings"
)

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

func ProcessQuantity(n string) (quantity float32, errID int) {
	raw := strings.SplitN(n, " ", 2)
	a, b, found := strings.Cut(n, "/")
	if len(raw) == 2 && found {
		pop := strings.SplitN(raw[1], "/", 2)
		numerator, denominator, whole, err := ConvertString(pop[0], pop[1], raw[0])
		if len(pop) != 2 || err != -1 {
			errID = 0
			return
		}

		quantity = whole + (numerator / denominator)
	} else if len(raw) == 1 && found {
		num, den, _, err := ConvertString(a, b, "")
		if err != -1 {
			errID = 0
			return
		}
		quantity = num / den
	} else {
		v, err := strconv.ParseFloat(raw[0], 32)
		if err != nil {
			errID = 0
			return
		}
		quantity = float32(v)
	}
	return quantity, -1
}

func CreateItem(name, priceTxt, costTxt, stockTxt string) (ID uint16, errID int) {
	quantity, err := ProcessQuantity(stockTxt)
	if err != -1 {
		errID = err
		return
	}

	price, cost, _, err := ConvertString(priceTxt, costTxt, "")
	if err != -1 {
		return 0, err
	}

	// Check for an open slot
	ID = uint16(len(Items))
	v, found := Items[ID]

	for found && v != nil {
		v, found = Items[ID]
		if !found || v == nil {
			break
		}
		ID += 1
	}

	Items[ID] = &Entry{Price: price, Name: name, Quantity: [3]float32{quantity, 0, 0}, Cost: [3]float32{cost, 0, 0}}
	fmt.Println("!Found, Adding: ", Items[ID])
	return ID, -1
}

func AddItem(target uint16, priceTxt, costTxt, stockTxt string) (errID int) {
	quan, err := ProcessQuantity(stockTxt)
	if err != -1 {
		return err
	}

	price, cost, _, err := ConvertString(priceTxt, costTxt, "")
	if err != -1 {
		return err
	}

	Items[target].Price = price

	rejections := 0
	for i := 0; i < 3; i++ {
		if Items[target].Cost[i] == cost {
			Items[target].Quantity[i] += quan
			break
		}
		if Items[target].Quantity[i] == 0 {
			Items[target].Quantity[i] = quan
			Items[target].Cost[i] = cost
			break
		}
		rejections += 1
	}

	if rejections == 3 {
		return 2
	}

	return -1
}
