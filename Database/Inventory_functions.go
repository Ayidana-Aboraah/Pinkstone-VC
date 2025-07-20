package Database

import (
	"BronzeHermes/Debug"
	"cmp"
	"slices"
	"strconv"
	"strings"
)

// func CleanUpDeadItems() {
// 	for ID := range Items {
// 		if Items[ID].Name[0] == byte(216) {
// 			found := false
// 			for _, x := range Sales {
// 				if x.ID == ID {
// 					found = true
// 					break
// 				}
// 			}
// 			if !found {
// 				DeadIDsQueue = append(DeadIDsQueue, ID)
// 			}
// 		}
// 	}
// }

func SearchInventory(input string) (Names []string, IDs []int) {
	if input == "" {
		for i, v := range Items {
			IDs = append(IDs, i)
			Names = append(Names, v.Name)
		}
		return
	}

	input = strings.TrimSpace(input)

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

	slices.SortStableFunc(inter, func(a, b int) int {
		return cmp.Compare(strings.ToLower(Items[uint16(a)].Name), strings.ToLower(Items[uint16(b)].Name))
	})

	return
} // 		if Items[ID].Name[0] == byte(216) {
// 			found := false
// 			for _, x := range Sales {
// 				if x.ID == ID {
// 					found = true
// 					break
// 				}
// 			}
// 			if !found {
// 				DeadIDsQueue = append(DeadIDsQueue, ID)
// 			}
// 		}
// 	}
// }

func ProcessQuantity(n string) (quantity float32, errID int) {

	raw := strings.SplitN(n, " ", 2)
	a, b, found := strings.Cut(raw[len(raw)-1], "/")

	if len(raw) == 2 && found {
		numerator, denominator, whole, err := ConvertString(a, b, raw[0])
		if err != Debug.Success {
			errID = Debug.Invalid_Input
			return
		}

		quantity = whole + (numerator / denominator)
	} else if len(raw) == 1 && found {
		num, den, _, err := ConvertString(a, b, "")
		if err != Debug.Success {
			errID = Debug.Invalid_Input

			return
		}
		quantity = num / den
	} else {
		v, err := strconv.ParseFloat(raw[0], 32)
		if err != nil {
			errID = Debug.Invalid_Input

			return
		}
		quantity = float32(v)
	}
	return quantity, Debug.Success
}

func CreateItem(name, priceTxt, costTxt, stockTxt string) (ID int, errID int) {
	quantity, err := ProcessQuantity(stockTxt)
	if err != Debug.Success {
		errID = err
		return
	}

	price, cost, _, err := ConvertString(priceTxt, costTxt, "")
	if err != Debug.Success {
		return 0, Debug.Invalid_Input
	}

	// Check for an open slot
	if len(DeadIDsQueue) == 0 {
		ID = len(Items)
	} else {
		ID = DeadIDsQueue[0]
		DeadIDsQueue = DeadIDsQueue[1:]
	}

	Items[ID] = Entry{Price: price, Name: name, Quantity: [3]float32{quantity, 0, 0}, Cost: [3]float32{cost, 0, 0}}
	return ID, Debug.Success
}

func AddItem(target uint16, priceTxt, costTxt, stockTxt string) (errID int) {
	quan, err := ProcessQuantity(stockTxt)
	if err != Debug.Success {
		return err
	}

	price, cost, _, err := ConvertString(priceTxt, costTxt, "")
	if err != Debug.Success {
		return err
	}

	Items[target].Price = price

	for i := 0; i < len(Items[target].Quantity); i++ {
		if Items[target].Cost[i] == cost {
			Items[target].Quantity[i] += quan
			break
		}
		if Items[target].Quantity[i] <= 0 {
			Items[target].Quantity[i] = quan
			Items[target].Cost[i] = cost
			break
		}
	}

	return Debug.Success
}
