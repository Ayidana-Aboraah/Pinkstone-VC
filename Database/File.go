package Database

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"
	"math"

	"fyne.io/fyne/v2"
)

var NameKeys = map[uint64]string{}

var Databases [3][]Sale
var Expenses []Expense

type Sale struct {
	Year, Month, Day uint8
	Quantity         uint16
	Price, Cost      float32
	ID               uint64
}

type Expense struct { // - for expense, + for gift
	Frequency uint8
	Date      [3]uint8 //Year, Month, Day,
	Amount    float32
	Name      string
}

const (
	ITEMS uint8 = iota
	REPORT
	LOG
)

const (
	ONCE uint8 = iota
	MONTHLY
	YEARLY
)
const DATA_SIZE = 19

func PutUint40(b []byte, v uint64) {
	_ = b[4]
	b[0] = byte(v >> 32)
	b[1] = byte(v >> 24)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 8)
	b[4] = byte(v)
}

func FromUint40(b []byte) uint64 {
	return uint64(b[4]) | uint64(b[3])<<8 | uint64(b[2])<<16 | uint64(b[1])<<24 | uint64(b[0])<<32
}

func DataInit() {
	for i, file := 0, ""; i < 6; i++ {
		switch i {
		case 0:
			file = "Item_Reference.red"
		case 1:
			file = "Report_Data.red"
		case 2:
			file = "Price_Log.red"
		case 3:
			file = "name_keys.json"
		case 4:
			file = "BackUp.red"
		case 5:
			file = "BackUp_Keys.red"
		}

		save, err := fyne.CurrentApp().Storage().Create(file)
		if err == nil {
			save.Close()
		}
	}
}

func SaveData() error {
	order := binary.BigEndian
	var file string

	for idx, database := range Databases {
		switch idx {
		case 0:
			file = "Item_Reference.red"
		case 1:
			file = "Report_Data.red"
		case 2:
			file = "Price_Log.red"
		}

		save, err := fyne.CurrentApp().Storage().Save(file)
		if err != nil {
			return err
		}

		bs := make([]byte, DATA_SIZE*len(database))
		for i, x := range database {
			c := (DATA_SIZE * i)

			bs[c] = x.Year
			bs[c+1] = x.Month
			bs[c+2] = x.Day

			order.PutUint16(bs[c+3:c+5], x.Quantity)
			order.PutUint32(bs[c+5:c+9], math.Float32bits(x.Price))
			order.PutUint32(bs[c+9:c+13], math.Float32bits(x.Cost))
			PutUint40(bs[c+13:c+19], x.ID)
		}

		_, err = save.Write(bs)
		save.Close()

		if err != nil {
			return err
		}
	}

	names, err := fyne.CurrentApp().Storage().Save("name_keys.json")

	if err != nil {
		return err
	}
	defer names.Close()

	encoder := json.NewEncoder(names)
	encoder.Encode(NameKeys)
	return nil
}

func LoadData() error {
	order := binary.BigEndian
	var file string

	for idx := range Databases {
		switch idx {
		case 0:
			file = "Item_Reference.red"
		case 1:
			file = "Report_Data.red"
		case 2:
			file = "Price_Log.red"
		}

		file, err := fyne.CurrentApp().Storage().Open(file)
		if err != nil {
			return err
		}

		buf, err := io.ReadAll(file)
		file.Close()

		if err != nil {
			return err
		}

		black := make([]Sale, len(buf)/DATA_SIZE)

		for i := range black {
			c := DATA_SIZE * i

			black[i].Year = uint8(buf[c])
			black[i].Month = uint8(buf[c+1])
			black[i].Day = uint8(buf[c+2])

			black[i].Quantity = order.Uint16(buf[c+3 : c+5])
			black[i].Price = math.Float32frombits(order.Uint32(buf[c+5 : c+9]))
			black[i].Cost = math.Float32frombits(order.Uint32(buf[c+9 : c+13]))
			black[i].ID = FromUint40(buf[c+13 : c+19])
		}

		Databases[idx] = black
	}

	names, err := fyne.CurrentApp().Storage().Open("name_keys.json")
	if err != nil && err != io.EOF {
		return err
	}
	defer names.Close()

	encoder := json.NewDecoder(names)
	err = encoder.Decode(&NameKeys)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func BackUpAllData() error {
	order := binary.BigEndian

	save, err := fyne.CurrentApp().Storage().Save("BackUp.red") //BackUp Loading
	if err != nil {
		return err
	}
	defer save.Close()

	bs := make([]byte, ((DATA_SIZE * len(Databases[0])) + (DATA_SIZE * len(Databases[1])) + (DATA_SIZE * len(Databases[2]))))

	previousLength := 0
	for _, database := range Databases {
		for i, x := range database {
			initial := (previousLength * DATA_SIZE) + (DATA_SIZE * i)

			bs[initial] = x.Year
			bs[initial+1] = x.Month
			bs[initial+2] = x.Day

			order.PutUint16(bs[initial+3:initial+5], x.Quantity)
			order.PutUint32(bs[initial+5:initial+9], math.Float32bits(x.Price))
			order.PutUint32(bs[initial+9:initial+13], math.Float32bits(x.Cost))
			PutUint40(bs[initial+13:initial+19], x.ID)
		}
		previousLength += len(database)
	}

	_, err = save.Write(bs)
	if err != nil {
		return err
	}

	names, err := fyne.CurrentApp().Storage().Save("BackUp_Keys.red")
	if err != nil {
		return err
	}
	defer names.Close()

	encoder := json.NewEncoder(names)
	err = encoder.Encode(NameKeys)
	if err != nil {
		return err
	}
	return nil
}

func LoadBackUp() error {
	order := binary.BigEndian

	file, err := fyne.CurrentApp().Storage().Open("BackUp.red")
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(file)
	file.Close()

	if err != nil {
		return err
	}

	black := make([]Sale, len(buf)/DATA_SIZE)

	for i, v := range black {
		c := DATA_SIZE * i

		black[i].Year = uint8(buf[c])
		black[i].Month = uint8(buf[c+1])
		black[i].Day = uint8(buf[c+2])

		black[i].Quantity = order.Uint16(buf[c+3 : c+5])
		black[i].Price = math.Float32frombits(order.Uint32(buf[c+5 : c+9]))
		black[i].Cost = math.Float32frombits(order.Uint32(buf[c+9 : c+13]))
		black[i].ID = FromUint40(buf[c+13 : c+19])

		if v.Year == 0 {
			Databases[0] = append(Databases[0], v)
		} else if v.Quantity == 0 {
			Databases[2] = append(Databases[2], v)
		} else {
			Databases[1] = append(Databases[1], v)
		}
	}

	names, err := fyne.CurrentApp().Storage().Open("BackUp_Keys.red")
	if err != nil {
		return err
	}
	defer names.Close()

	decoder := json.NewDecoder(names)
	err = decoder.Decode(&NameKeys)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
