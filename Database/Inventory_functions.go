package Database

import (
	"BronzeHermes/Helpers"
	"cmp"
	"slices"
	"strings"
)

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

func ConvertItems() (inter []int) {
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

func ProcessQuantity(text string) (quantity float64, errID int) {
	const (
		whole = iota
		numerator
		denominator
	)

	// NOTE: Bunch of weird cases that can come from this, idk what to do about it
	a := strings.Replace(text, "/", " ", 1)
	b := strings.Split(a, " ")
	fraction := Helpers.ParseToFloat(b...)
	if fraction == nil {
		return 0, Helpers.Invalid_Input
	}

	switch len(fraction) {
	case 3:
		quantity = fraction[whole] + (fraction[numerator] / fraction[denominator])
	case 2:
		quantity = fraction[numerator] / fraction[denominator]
	case 1:
		quantity = fraction[whole]
	}

	return quantity, Helpers.Success
}

func CreateItem(name, priceTxt, costTxt, stockTxt string) (ID int, errID int) {

	res := Helpers.ParseToFloat(priceTxt, costTxt)
	if res == nil {
		return 0, Helpers.Invalid_Input
	}

	quantity, err := ProcessQuantity(stockTxt)
	if err != Helpers.Success {
		return 0, err
	}

	// Check for an open slot
	if len(DeadIDsQueue) == 0 {
		ID = len(Items)
		Items = append(Items, Entry{})
	} else {
		ID = DeadIDsQueue[0]
		DeadIDsQueue = DeadIDsQueue[1:]
	}

	Items[ID] = Entry{Price: res[0], Name: name, Quantity: []float64{quantity, 0, 0}, Cost: []float64{res[1], 0, 0}}
	return ID, Helpers.Success
}

func AddItem(target uint, priceTxt, costTxt, stockTxt string) (errID int) {
	quan, err := ProcessQuantity(stockTxt)
	if err != Helpers.Success {
		return err
	}

	res := Helpers.ParseToFloat(priceTxt, costTxt)
	if res == nil {
		return Helpers.Invalid_Input
	}

	Items[target].Price = res[0]

	for i := 0; i < len(Items[target].Quantity); i++ {
		if Items[target].Cost[i] == res[1] {
			Items[target].Quantity[i] += quan
			break
		}
		if Items[target].Quantity[i] <= 0 {
			Items[target].Quantity[i] = quan
			Items[target].Cost[i] = res[1]
			break
		}
	}

	return Helpers.Success
}
