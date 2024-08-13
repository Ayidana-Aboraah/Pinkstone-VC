package Database

import (
	"BronzeHermes/Debug"
	"cmp"
	"slices"
	"strconv"
	"strings"
)

func CleanUpDeadItems() {
	for k := range Items {
		if Items[k].Name[0] == byte(216) {
			found := false
			for _, x := range Sales {
				if x.ID == uint16(k) {
					found = true
					break
				}
			}
			if !found {
				Items[k].Name = strings.Replace(Items[k].Name, string([]byte{216}), string([]byte{171}), -1)
			}
		}
	}
}

func SearchInventory(input string) (Names []string, IDs []uint16) {
	if input == "" {
		for i, v := range Items {
			IDs = append(IDs, uint16(i))
			Names = append(Names, v.Name)
		}
		return
	}

	input = strings.TrimSpace(input)

	for id, iv := range Items {
		if strings.Contains(strings.ToLower(iv.Name), strings.ToLower(input)) && iv.Name[0] != byte(216) {
			IDs = append(IDs, uint16(id))
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
}

// func orderItemKeys(list []int) (out []int) {
// 	for _, v := range list {
// 		init_out_len := len(out)

// 		for i, o := range out {
// 			if v > o {
// 				// insert into array
// 				out = append(out[:i+1], out[i:]...)
// 				out[i] = v
// 				break
// 			}
// 		}

// 		if init_out_len == len(out) {
// 			out = append(out, v)
// 		}
// 	}
// 	return
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

func CreateItem(name, priceTxt, costTxt, stockTxt string) (ID uint16, errID int) {
	quantity, err := ProcessQuantity(stockTxt)
	if err != Debug.Success {
		errID = err
		return
	}

	price, cost, _, err := ConvertString(priceTxt, costTxt, "")
	if err != Debug.Success {
		return 0, Debug.Invalid_Input
	}

	// TODO: Maybe Check If we're creating a new itemthat  already exists

	replace := -1

	for i := 0; i < len(Items) && replace == -1; i++ {
		if Items[i].Name[0] == byte(171) {
			replace = i
		}
	}

	if replace == -1 {
		Items = append(Items, Item{Price: price, Name: name, Quantity: [3]float32{quantity, 0, 0}, Cost: [3]float32{cost, 0, 0}})
	} else {
		Items[replace] = Item{Price: price, Name: name, Quantity: [3]float32{quantity, 0, 0}, Cost: [3]float32{cost, 0, 0}}
	}

	// fmt.Println("!Found, Adding: ", Items[ID]) // LOG
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

	i := 0
	for ; i < len(Items[target].Cost); i++ {
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

	if i == len(Items[target].Cost) {
		return Debug.Maxed_Out_Stocks
	}

	return Debug.Success
}
