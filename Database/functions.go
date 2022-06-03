package Database

import (
	"strconv"
	"strings"
)

func GetLine(selection string, dataType int, database []Sale) ([]string, [][]float32) {
	if selection == "" {
		return nil, nil
	}

	raw := strings.Split(selection, "/")

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil, nil
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, nil
	}

	date := []uint8{
		uint8(year),
		uint8(month),
	}

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
				case 3:
					total += float32(v.Quantity)
				}
			}

			totals = append(totals, total)
		}

		if totals != nil {
			names = append(names, name)
			sales = append(sales, totals)
		}
	}

	return names, sales
}

func GetPie(selection string, dataType int) ([]string, []float32) {
	if selection == "" {
		return nil, nil
	}

	raw := strings.Split(selection, "/")

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil, nil
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, nil
	}

	day, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, nil
	}

	date := []uint8{
		uint8(year),
		uint8(month),
		uint8(day),
	}

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
			case 3:
				total += float32(v.Quantity)
			}
		}

		if total > 0 {
			names = append(names, name)
			sales = append(sales, total)
		}
	}

	return names, sales
}

func FindItem(ID int) Sale { // Maybe implement a binary search
	for _, v := range Databases[0] {
		if int(v.ID) == ID {
			v.Quantity = 1 // If changing to for loop, make a direct copy when submitting
			return v
		}
	}
	return Sale{}
}
