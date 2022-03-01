package Test

import (
	"BronzeHermes/Database"
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"testing"
)

//INTEGRATION TESTS
func GenerateTestData(t *testing.T) {
	//Create Sale List
	database := []Database.Sale{
		{ID: 108499092, Year: 20, Day: 30, Month: 10},
		{ID: 491289142},
		{ID: 101984123},
		{ID: 128004121},
		{ID: 410283561},
		{ID: 184098561},
		{ID: 194820512},
		{ID: 222222313},
		{ID: 492780612},
	}

	//Create Keys from Sale List
	keys := map[uint64]string{}

	for i := 0; i < len(database)-2; i++ {
		keys[database[i].ID] = string(rune(database[i].ID))
	}

	//Save Sale List
	order := binary.BigEndian

	save, err := os.OpenFile("test_integration.red", os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Error(err)
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

	if err != nil {
		t.Error(err)
	}

	//Save Keys
	names, err := os.OpenFile("test_intgration_names_encoded.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	defer names.Close()

	encoder := json.NewEncoder(names)
	err = encoder.Encode(keys)
	if err != nil {
		t.Error(err)
	}
}

func ReadTestData(t *testing.T) ([]Database.Sale, map[uint64]string) {
	//Read Sales
	order := binary.BigEndian

	buf, err := ioutil.ReadFile("test_integration.red")
	if err != nil {
		t.Error(err)
	}

	black := make([]Database.Sale, len(buf)/21)

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

	t.Log(black)

	//Read Keys
	var blue map[uint64]string
	names, err := os.OpenFile("test_intgration_names_encoded.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer names.Close()

	encoder := json.NewDecoder(names)
	err = encoder.Decode(&blue)
	if err != nil {
		t.Error(err)
	}
	t.Log(blue)
	return black, blue
}

func TestKeyCreationFromDatabases(t *testing.T) {
	GenerateTestData(t)
	black, blue := ReadTestData(t)

	keyTotal := len(black) - 2
	currentKeys := 0

	for _, v := range black {
		if _, found := blue[v.ID]; found {
			currentKeys++
		}
	}

	//If the total keys don't match total sales than it's inncorrect
	if currentKeys != keyTotal {
		t.Errorf("Current Keys: %v, when the total keys are %v", currentKeys, keyTotal)
	}
}
