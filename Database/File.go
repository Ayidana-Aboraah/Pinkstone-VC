package Database

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
)

var NameKeys map[uint32]string

var Items []Sale
var ReportData []Sale
var PriceLog []Sale

type Sale struct {
	Year     uint8
	Month    uint8
	Day      uint8
	Quantity uint16
	ID       uint32
	Price    float32
	Cost     float32
}

const (
	ITEMS = iota
	REPORT
	LOG
)

func SaveKeys() error {
	names, err := os.OpenFile("name_to_keys.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer names.Close()

	encoder := json.NewEncoder(names)
	encoder.Encode(NameKeys)
	return nil
}

func LoadKeys() error {
	names, err := os.OpenFile("name_to_keys.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer names.Close()

	encoder := json.NewDecoder(names)
	err = encoder.Decode(&NameKeys)
	if err != nil {
		return err
	}
	return nil
}

func SaveData(variant int) error {
	var database []Sale
	var file string

	switch variant {
	case ITEMS:
		database = Items
		file = "Item_Reference.red"
	case REPORT:
		database = ReportData
		file = "ReportData.red"
	case LOG:
		database = PriceLog
		file = "PriceLog.red"
	}

	order := binary.BigEndian

	save, err := os.OpenFile(file, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer save.Close()

	bs := make([]byte, 17*len(database))
	for i, x := range database {
		c := (17 * i)

		bs[c] = x.Year
		bs[c+1] = x.Month
		bs[c+2] = x.Day

		order.PutUint16(bs[c+3:c+5], x.Quantity)
		order.PutUint32(bs[c+5:c+9], x.ID)
		order.PutUint32(bs[c+9:c+13], math.Float32bits(x.Price))
		order.PutUint32(bs[c+13:c+17], math.Float32bits(x.Cost))
	}
	_, err = save.Write(bs)
	return err
}

func LoadData(variant int) ([]Sale, error) {
	var file string

	switch variant {
	case ITEMS:
		file = "Item_Reference.red"
	case REPORT:
		file = "ReportData.red"
	case LOG:
		file = "PriceLog.red"
	}

	order := binary.BigEndian

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return []Sale{}, err
	}

	black := make([]Sale, len(buf)/17)

	for i := range black {
		c := 17 * i

		black[i].Year = uint8(buf[c])
		black[i].Month = uint8(buf[c+1])
		black[i].Day = uint8(buf[c+2])

		black[i].Quantity = order.Uint16(buf[c+3 : c+5])
		black[i].ID = order.Uint32(buf[c+5 : c+9])
		black[i].Price = math.Float32frombits(order.Uint32(buf[c+9 : c+13]))
		black[i].Cost = math.Float32frombits(order.Uint32(buf[c+13 : c+17]))
	}

	return black, nil
}
