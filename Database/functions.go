package Database

import (
	"strconv"
	"strings"
)

func AddKey(id int, name string) {
	newKeys := make(map[uint64]string, len(NameKeys)+1)
	for idx, name := range NameKeys {
		newKeys[idx] = name
	}
	newKeys[uint64(id)] = name
	NameKeys = newKeys
}

func FindItem(ID int) Sale {
	for _, v := range Databases[0] {
		if int(v.ID) == ID {
			return v
		}
	}
	return Sale{}
}

func GetLine(selection string, dataType int, database []Sale) ([]string, [][]float32) {
	date := func() []uint8 {
		if selection == "" {
			return nil
		}

		raw := strings.Split(selection, "/")

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
	}()

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
	date := func() []uint8 {
		if selection == "" {
			return nil
		}

		raw := strings.Split(selection, "/")

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
	}()

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
