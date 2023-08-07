package unknown

import (
	"BronzeHermes/Debug"
	"strconv"
	"strings"
)

func RemoveAndUpdate(name *string, update func()) {
	//TODO: Add Logging
	*name = string([]byte{216}) + *name
	update()
}

func AddToNames(names []string, name string) ([]string, uint8) {
	// TODO: Logging
	return append(names, name), uint8(len(names))
}

func ProcessDate(in string) (out [3]uint8, errID int) {
	pop := strings.SplitN(in, "-", 3)
	if len(pop) < 3 {
		return out, Debug.Need_More_Info
	}

	year, err := strconv.Atoi(pop[0][1:])
	if err != nil {
		return out, Debug.Invalid_Input
	}
	out[2] = uint8(year)

	month, err := strconv.Atoi(pop[1])
	if err != nil {
		return out, Debug.Invalid_Input
	}
	out[1] = uint8(month)

	day, err := strconv.Atoi(pop[2])
	if err != nil {
		return out, Debug.Invalid_Input
	}
	out[0] = uint8(day)

	return out, Debug.Success
}

func ProcessDate2(in string) (out [3]uint8, errID int) {
	pop := strings.SplitN(in, "-", 3)

	if len(pop) == 0 {
		return out, Debug.Invalid_Input
	}

	year, err := strconv.Atoi(pop[0][1:])
	if err != nil {
		return out, Debug.Invalid_Input
	}
	out[0] = uint8(year)

	if len(pop) > 1 {
		month, err := strconv.Atoi(pop[1])
		if err != nil {
			return out, Debug.Invalid_Input
		}
		out[1] = uint8(month)
	}

	if len(pop) > 2 {
		day, err := strconv.Atoi(pop[2])
		if err != nil {
			return out, Debug.Invalid_Input
		}
		out[2] = uint8(day)
	}

	return out, Debug.Success
}
