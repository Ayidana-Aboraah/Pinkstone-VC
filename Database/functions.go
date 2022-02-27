package Database

import (
	"strconv"
	"strings"
)

func GetSalesLine(selection string) ([]string, [][]float32) {
	date := ConvertDateForLine(selection)

	//Change the error handling for this to show that you can't convert
	var sales [][]float32
	var names []string

	for id, name := range NameKeys {
		var totals []float32

		for i := uint8(1); i < 32; i++ {
			var total uint16
			for _, v := range Databases[1] {

				if v.ID != id || v.Day != i || v.Month != date[1] || v.Year != date[0] {
					continue
				}

				total += v.Quantity
			}

			totals = append(totals, float32(total))
		}

		if totals == nil {
			continue
		}

		names = append(names, name)
		sales = append(sales, totals)
	}

	return names, sales
}

func GetLine(selection string, dataType int, database []Sale) ([]string, [][]float32) {
	date := ConvertDateForLine(selection)

	//Change the error handling for this to show that you can't convert
	var sales [][]float32
	var names []string

	for id, name := range NameKeys {
		var totals []float32

		for i := uint8(1); i < 32; i++ {
			var total float32
			for _, v := range database {

				if v.ID != id || v.Day != i || v.Month != date[1] || v.Year != date[0] {
					continue
				}

				switch dataType {
				case 0:
					total += v.Price
				case 1:
					total += v.Cost
				case 2:
					total += v.Price - v.Cost
				}

			}

			totals = append(totals, total)
		}

		if totals == nil {
			continue
		}

		names = append(names, name)
		sales = append(sales, totals)
	}

	return names, sales
}

func GetSalesPie(selection string) ([]string, []float32) {
	date := ConvertDateForPie(selection)

	//Change the error handling for this to show that you can't convert
	var sales []float32
	var names []string

	for id, name := range NameKeys {
		var total uint16

		for _, v := range Databases[1] {

			if v.ID != id || v.Day != date[2] || v.Month != date[1] || v.Year != date[0] {
				continue
			}

			total += v.Quantity
		}

		if total == 0 {
			continue
		}

		names = append(names, name)
		sales = append(sales, float32(total))
	}

	return names, sales
}

func GetPricePie(selection string, dataType int) ([]string, []float32) {
	date := ConvertDateForPie(selection)

	var sales []float32
	var names []string

	for id, name := range NameKeys {
		var total float32

		for _, v := range Databases[1] {

			if v.ID != id || v.Day != date[2] || v.Month != date[1] || v.Year != date[0] {
				continue
			}

			switch dataType {
			case 0:
				total += v.Price
			case 1:
				total += v.Cost
			case 2:
				total += v.Price - v.Cost
			}
		}

		if total == 0 {
			continue
		}

		names = append(names, name)
		sales = append(sales, total)
	}

	return names, sales
}

func ConvertDateForPie(date string) []uint8 {
	if date == "" {
		return nil
	}

	raw := strings.Split(date, "/")

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil
	}

	day, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil
	}

	return []uint8{
		uint8(year),
		uint8(month),
		uint8(day),
	}
}

func ConvertDateForLine(date string) []uint8 {
	if date == "" {
		return nil
	}

	raw := strings.Split(date, "/")

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil
	}

	return []uint8{
		uint8(year),
		uint8(month),
	}
}

func AddKey(id uint64, name string) {
	newKeys := make(map[uint64]string, len(NameKeys)+1)
	for idx, name := range NameKeys {
		newKeys[idx] = name
	}
	newKeys[id] = name
}

func FindItem(ID uint64) Sale {
	for _, v := range Databases[0] {
		if v.ID == ID {
			return v
		}
	}
	return Sale{}
}
