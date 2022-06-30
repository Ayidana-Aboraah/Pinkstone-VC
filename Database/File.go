package Database

import (
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"math"
	"strings"

	"fyne.io/fyne/v2"
)

var NameKeys = map[uint64]string{}
var ItemKeys = map[uint64]struct {
	Price float32
	idxes []uint16
	Name  string
}{}

var Items []struct {
	Quantity uint16
	Cost     float32
	ID       uint64
}

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
	Date      [3]uint8 // Day, Month Year
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
const DATA_SIZE = 18

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

func DataInit(remove bool) {
	for i, file := 0, ""; i < 7; i++ {
		switch i {
		case 0:
			file = "Item_Reference.red"
		case 1:
			file = "Report_Data.red"
		case 2:
			file = "Price_Log.red"
		case 3:
			file = "Name_Map.red"
		case 4:
			file = "Expenses.red"
		case 5:
			file = "BackUp.red"
		case 6:
			file = "BackUp_Map.red"
		}

		if !remove {
			save, err := fyne.CurrentApp().Storage().Create(file)
			if err == nil {
				save.Close()
			}
		} else {
			err := fyne.CurrentApp().Storage().Remove(file)
			if err != nil {
				panic(err) //IDK what else to do
			}
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
			PutUint40(bs[c+13:c+DATA_SIZE], x.ID)
		}

		_, err = save.Write(bs)
		save.Close()

		if err != nil {
			return err
		}
	}

	expenses, err := fyne.CurrentApp().Storage().Save("Expenses.red")
	if err != nil && err != io.EOF {
		return err
	}
	defer expenses.Close()

	var res []byte
	for _, v := range Expenses {
		buff := make([]byte, 8)
		buff[0] = v.Frequency
		buff[1] = v.Date[0]
		buff[2] = v.Date[1]
		buff[3] = v.Date[2]
		order.PutUint32(buff[4:], math.Float32bits(v.Amount))
		if v.Name == "" {
			v.Name = "_"
		}
		buff = append(buff, []byte(v.Name+"\n")...)
		res = append(res, buff...)
	}

	_, err = expenses.Write(res)
	if err != nil {
		return err
	}

	// New Key saving
	names, err := fyne.CurrentApp().Storage().Save("Name_Map.red")
	if err != nil {
		return err
	}
	defer names.Close()
	var result []byte
	for k, v := range NameKeys {
		mine := make([]byte, 5)
		PutUint40(mine, k)
		mine = append(mine, []byte(v+"\n")...)
		result = append(result, mine...)
	}
	_, err = names.Write(result)

	return err
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
			black[i].ID = FromUint40(buf[c+13 : c+DATA_SIZE])
		}

		Databases[idx] = black
	}

	expenses, err := fyne.CurrentApp().Storage().Open("Expenses.red")
	if err != nil && err != io.EOF {
		return err
	}
	defer expenses.Close()

	exp_bytes, err := io.ReadAll(expenses)
	if err != nil && err != io.EOF {
		return err
	}

	raw := strings.Split(string(exp_bytes), "\n")
	for _, v := range raw[:len(raw)-1] {
		var expense Expense
		process := []byte(v)
		expense.Frequency = process[0]
		expense.Date[0] = process[1]
		expense.Date[1] = process[2]
		expense.Date[2] = process[3]
		expense.Amount = math.Float32frombits(order.Uint32(process[4:8])) //TODO: Look into using a float16 to save space
		expense.Name = v[8:]
		if expense.Name == "_" {
			expense.Name = ""
		}
		Expenses = append(Expenses, expense)
	}

	names, err := fyne.CurrentApp().Storage().Open("Name_Map.red")
	if err != nil && err != io.EOF {
		return err
	}
	defer names.Close()

	// New Key Loading
	arr, err := io.ReadAll(names)
	if err != nil && err != io.EOF {
		return err
	}

	entries := strings.Split(string(arr), "\n")

	for _, v := range entries[:len(entries)-1] {
		NameKeys[FromUint40([]byte(v[:5]))] = v[5:]
	}

	return nil
}

func BackUpAllData() error { // TODO: Fix this function
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
			order.PutUint32(bs[initial+5:initial+9], math.Float32bits(x.Price)) //TODO: Look into using a float16 to save space
			order.PutUint32(bs[initial+9:initial+13], math.Float32bits(x.Cost)) //TODO: Look into using a float16 to save space
			PutUint40(bs[initial+13:initial+DATA_SIZE], x.ID)
		}
		previousLength += len(database)
	}

	_, err = save.Write(bs)
	if err != nil {
		return err
	}

	names, err := fyne.CurrentApp().Storage().Save("BackUp_Map.red")
	if err != nil {
		return err
	}
	defer names.Close()

	// TODO: Add expenses to this save
	var res []byte
	for _, v := range Expenses {
		buff := make([]byte, 8)
		buff[0] = v.Frequency
		buff[1] = v.Date[0]
		buff[2] = v.Date[1]
		buff[3] = v.Date[2]
		order.PutUint32(buff[4:], math.Float32bits(v.Amount))
		if v.Name == "" {
			v.Name = "_"
		}
		buff = append(buff, []byte(v.Name+"\n")...)
		res = append(res, buff...)
	}

	res = append(res, []byte("\n\n")...) // Just append a break character

	var result []byte
	for k, v := range NameKeys {
		mine := make([]byte, 6)
		PutUint40(mine, k)
		mine = append(mine, []byte(v+"\n")...)
		result = append(result, mine...)
	}
	result = append(res, result...)
	_, err = names.Write(result)
	return err
}

func LoadBackUp() error { //TODO: Fix this function
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

	for i := range black {
		c := DATA_SIZE * i

		black[i].Year = uint8(buf[c])
		black[i].Month = uint8(buf[c+1])
		black[i].Day = uint8(buf[c+2])

		black[i].Quantity = order.Uint16(buf[c+3 : c+5])
		black[i].Price = math.Float32frombits(order.Uint32(buf[c+5 : c+9]))
		black[i].Cost = math.Float32frombits(order.Uint32(buf[c+9 : c+13]))
		black[i].ID = FromUint40(buf[c+13 : c+DATA_SIZE])

		if black[i].Year == 0 {
			Databases[0] = append(Databases[0], black[i])
		} else if black[i].Quantity == 0 {
			Databases[2] = append(Databases[2], black[i])
		} else {
			Databases[1] = append(Databases[1], black[i])
		}
	}

	names, err := fyne.CurrentApp().Storage().Open("BackUp_Map.red")
	if err != nil {
		return err
	}
	defer names.Close()

	// New Key Loading
	raw, err := io.ReadAll(names)
	if err != nil && err != io.EOF {
		return err
	}

	exp, NameKV, found := strings.Cut(string(raw), "\n\n")
	if !found {
		return errors.New("Data not found in BackUp_Map.red")
	}

	exp_data := strings.Split(exp, "\n")
	for _, v := range exp_data {
		var expense Expense
		process := []byte(v)
		expense.Frequency = process[0]
		expense.Date[0] = process[1]
		expense.Date[1] = process[2]
		expense.Date[2] = process[3]
		expense.Amount = math.Float32frombits(order.Uint32(process[4:8])) //TODO: Look into using a float16 to save space
		expense.Name = v[8:]
		if expense.Name == "_" {
			expense.Name = ""
		}
		Expenses = append(Expenses, expense)
	}

	entries := strings.Split(NameKV[2:], "\n")

	for _, v := range entries[:len(entries)-1] {
		NameKeys[FromUint40([]byte(v[:5]))] = v[5:]
	}

	return nil
}
