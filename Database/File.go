package Database

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
)

var NameKeys map[uint64]string

var Items []Sale
var ReportData []Sale
var PriceLog []Sale

type Sale struct {
	Year     uint8
	Month    uint8
	Day      uint8
	Quantity uint16
	Price    float32
	Cost     float32
	ID       uint64
}

func SaveKeys() error {
	names, err := os.OpenFile("name_keys.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer names.Close()

	encoder := json.NewEncoder(names)
	encoder.Encode(NameKeys)
	return nil
}

func LoadKeys() error {
	names, err := os.OpenFile("name_keys.json", os.O_CREATE, os.ModePerm)
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

	bs := make([]byte, 21*len(database))
	for i, x := range database {
		c := (21 * i)

		bs[c] = x.Year
		bs[c+1] = x.Month
		bs[c+2] = x.Day

		order.PutUint16(bs[c+3:c+5], x.Quantity)
		order.PutUint32(bs[c+5:c+9], math.Float32bits(x.Price))
		order.PutUint32(bs[c+9:c+13], math.Float32bits(x.Cost))
		order.PutUint64(bs[c+13:c+21], x.ID)
	}
	_, err = save.Write(bs)
	return err
}

func LoadData(variant int) ([]Sale, error) {
	var file string

	switch variant {
	case 0:
		file = "Item_Reference.red"
	case 1:
		file = "ReportData.red"
	case 2:
		file = "PriceLog.red"
	}

	order := binary.BigEndian

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return []Sale{}, err
	}

	black := make([]Sale, len(buf)/21)

	for i := range black {
		c := 21 * i

		black[i].Year = uint8(buf[c])
		black[i].Month = uint8(buf[c+1])
		black[i].Day = uint8(buf[c+2])

		black[i].Quantity = order.Uint16(buf[c+3 : c+5])
		black[i].Price = math.Float32frombits(order.Uint32(buf[c+5 : c+9]))
		black[i].Cost = math.Float32frombits(order.Uint32(buf[c+9 : c+13]))
		black[i].ID = order.Uint64(buf[c+13 : c+21])
	}

	return black, nil
}

func BackUpAllData() error {
	names, err := os.OpenFile("Backup_Keys.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer names.Close()

	encoder := json.NewEncoder(names)
	err = encoder.Encode(NameKeys)
	if err != nil {
		return err
	}

	order := binary.BigEndian

	save, err := os.OpenFile("BackUp.red", os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer save.Close()

	databases := [][]Sale{Items, ReportData, PriceLog}
	bs := make([]byte, ((21 * len(databases[0])) + (21 * len(databases[1])) + (21 * len(databases[2]))))

	previousLength := 0
	for _, database := range databases {
		for i, x := range database {
			initial := (previousLength * 21) + (21 * i)

			bs[initial] = x.Year
			bs[initial+1] = x.Month
			bs[initial+2] = x.Day

			order.PutUint16(bs[initial+3:initial+5], x.Quantity)
			order.PutUint32(bs[initial+5:initial+9], math.Float32bits(x.Price))
			order.PutUint32(bs[initial+9:initial+13], math.Float32bits(x.Cost))
			order.PutUint64(bs[initial+13:initial+21], x.ID)
		}
		previousLength = len(database)
	}

	_, err = save.Write(bs)
	return err
}

func LoadBackUp() error {
	order := binary.BigEndian

	buf, err := ioutil.ReadFile("BackUp.red")
	if err != nil {
		return err
	}

	black := make([]Sale, len(buf)/21)

	for i, v := range black {
		c := 21 * i

		black[i].Year = uint8(buf[c])
		black[i].Month = uint8(buf[c+1])
		black[i].Day = uint8(buf[c+2])

		black[i].Quantity = order.Uint16(buf[c+3 : c+5])
		black[i].Price = math.Float32frombits(order.Uint32(buf[c+5 : c+9]))
		black[i].Cost = math.Float32frombits(order.Uint32(buf[c+9 : c+13]))
		black[i].ID = order.Uint64(buf[c+13 : c+21])

		if v.Year == 0 {
			Items = append(Items, v)
		} else if v.Quantity == 0 {
			PriceLog = append(PriceLog, v)
		} else {
			ReportData = append(ReportData, v)
		}
	}

	return nil
}
