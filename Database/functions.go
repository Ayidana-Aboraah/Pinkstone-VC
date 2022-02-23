package Database

import (
	"strconv"
	"strings"
)

func GetSalesLine(selection string) ([]string, [][]uint16) {
	raw := strings.Split(selection, "/")

	yr, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil, nil
	}

	mon, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, nil
	}

	year := uint8(yr)
	month := uint8(mon)

	//Change the error handling for this to show that you can't convert
	var sales [][]uint16
	var names []string

	for id, name := range NameKeys {
		var totals []uint16

		for i := uint8(1); i < 32; i++ {
			var total uint16
			for _, v := range ReportData {

				if v.ID != id || v.Day != i || v.Month != month || v.Year != year {
					continue
				}

				total += v.Quantity
			}

			totals = append(totals, total)
		}

		if totals != nil {
			continue
		}

		names = append(names, name)
		sales = append(sales, totals)
	}

	return names, sales
}

func GetPriceLine(selection string, dataType int) ([]string, [][]float32) {
	raw := strings.Split(selection, "/")

	yr, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil, nil
	}

	mon, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, nil
	}

	year := uint8(yr)
	month := uint8(mon)

	//Change the error handling for this to show that you can't convert
	var sales [][]float32
	var names []string

	for id, name := range NameKeys {
		var totals []float32

		for i := uint8(1); i < 32; i++ {
			var total float32
			for _, v := range ReportData {

				if v.ID != id || v.Day != i || v.Month != month || v.Year != year {
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

		if totals != nil {
			continue
		}

		names = append(names, name)
		sales = append(sales, totals)
	}

	return names, sales
}
