package Database

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
)

var NameKeys map[uint32]string

//Database for items
var Items []Sale

//Database for Reports
var ReportData []Sale

//Database for Log
var PriceLog []Sale

type Sale struct {
	ID       uint32
	Quantity uint16
	Price    float64
	Cost     float64
}

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
	err = encoder.Decode(NameKeys)
	if err != nil {
		return err
	}
	return nil
}

func SaveData(variant int) error {
	var database []Sale
	var file string

	switch variant {
	case 0:
		database = Items
		file = "Item_Reference.red"
	case 1:
		database = ReportData
		file = "ReportData.red"
	case 2:
		database = PriceLog
		file = "PriceLog.red"
	}

	order := binary.BigEndian

	save, err := os.OpenFile(file, os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer save.Close()

	bs := make([]byte, 22*len(database))
	for i, x := range database {
		c := (22 * i)

		order.PutUint32(bs[c:c+4], x.ID)
		order.PutUint16(bs[c+4:c+6], x.Quantity)
		order.PutUint64(bs[c+6:c+14], math.Float64bits(x.Price))
		order.PutUint64(bs[c+14:c+22], math.Float64bits(x.Cost))
	}
	_, err = save.Write(bs)
	return err
}

func LoadData(variant int) (error, []Sale) {
	order := binary.BigEndian

	buf, err := ioutil.ReadFile("test_save.red")
	if err != nil {
		return err, []Sale{}
	}

	black := make([]Sale, len(buf)/22)

	for i := range black {
		c := 22 * i
		black[i].ID = order.Uint32(buf[c : c+4])
		black[i].Quantity = order.Uint16(buf[c+4 : c+6])
		black[i].Price = math.Float64frombits(order.Uint64(buf[c+6 : c+14]))
		black[i].Cost = math.Float64frombits(order.Uint64(buf[c+14 : c+22]))
	}

	return err, black
}
