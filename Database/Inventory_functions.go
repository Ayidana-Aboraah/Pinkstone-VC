package Database

import (
	"BronzeHermes/UI"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
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

func ProcessQuantity(n string, w fyne.Window) (quantity float32) {
	raw := strings.SplitN(n, " ", 2)
	if len(raw) == 2 {
		pop := strings.SplitN(raw[1], "/", 2)
		numerator, denominator, whole := ConvertString(pop[0], pop[1], raw[0])
		if UI.HandleKnownError(0, len(pop) != 2, w) {
			return
		}
		quantity = whole + (numerator / denominator)
	} else {
		v, err := strconv.ParseFloat(raw[0], 32)
		if UI.HandleKnownError(0, err != nil, w) {
			return
		}
		quantity = float32(v)
	}
	return
}
