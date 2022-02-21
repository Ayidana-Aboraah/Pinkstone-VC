package Database

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"testing"
)

var red []Sale
var blue map[int16]string

type Sale struct {
	ID       uint16
	Quantity uint16
	Price    float64
	Cost     float64
}

func SaveData() error {
	order := binary.BigEndian
	names, err := os.OpenFile("test_name_encode.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer names.Close()

	encoder := json.NewEncoder(names)
	encoder.Encode(blue)

	save, err := os.OpenFile("test_save.red", os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer save.Close()

	bs := make([]byte, 20*len(red))
	for i, x := range red {
		c := (20 * i)

		order.PutUint16(bs[c:c+2], x.ID)
		order.PutUint16(bs[c+2:c+4], x.Quantity)
		order.PutUint64(bs[c+4:c+12], math.Float64bits(x.Price))
		order.PutUint64(bs[c+12:c+20], math.Float64bits(x.Cost))
	}

	_, err = save.Write(bs)
	return err
}

func LoadData() error {
	order := binary.BigEndian
	names, err := os.OpenFile("test_name_encode.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer names.Close()

	encoder := json.NewDecoder(names)
	err = encoder.Decode(blue)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadFile("test_save.red")
	if err != nil {
		return err
	}

	black := make([]Sale, len(buf)/20)

	for i, _ := range black {
		c := (20 * i)
		black[i].ID = order.Uint16(buf[c : c+2])
		black[i].Quantity = order.Uint16(buf[c+2 : c+4])
		black[i].Price = math.Float64frombits(order.Uint64(buf[c+4 : c+12]))
		black[i].Cost = math.Float64frombits(order.Uint64(buf[c+12 : c+20]))
	}

	red = black
	return err
}

func TestBlue(t *testing.T) {
}
