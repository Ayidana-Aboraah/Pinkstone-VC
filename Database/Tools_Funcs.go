package Database

import (
	"strconv"
	"strings"
)

func ConvertDateForPie(date string) ([]uint8, error) {
	raw := strings.Split(date, "/")

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil, err
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, err
	}

	day, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, err
	}

	return []uint8{
		uint8(year),
		uint8(month),
		uint8(day),
	}, nil
}

func ConvertDateForLine(date string) ([]uint8, error) {
	raw := strings.Split(date, "/")

	year, err := strconv.Atoi(raw[0][1:])
	if err != nil {
		return nil, err
	}

	month, err := strconv.Atoi(raw[1])
	if err != nil {
		return nil, err
	}

	return []uint8{
		uint8(year),
		uint8(month),
	}, nil

}
