package Database

import (
	"BronzeHermes/Debug"
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

func FilterUsers() (out []int) {
	for i, s := range Users {
		if s[0] != byte(216) {
			out = append(out, i)
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

	return orderItemKeys(inter)
}

func orderItemKeys(list []int) (out []int) {
	for _, v := range list {
		init_out_len := len(out)

		for i, o := range out {
			if v > o {
				// insert into array
				out = append(out[:i+1], out[i:]...)
				out[i] = v
				break
			}
		}

		if init_out_len == len(out) {
			out = append(out, v)
		}
	}
	return
}

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
	for ; i < 3; i++ {
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

	if i == 3 {
		return Debug.Maxed_Out_Stocks
	}

	return Debug.Success
}
