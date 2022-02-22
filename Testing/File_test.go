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

//UNIT TESTS
func TestSaveKeys(t *testing.T) {
	blue := map[uint32]string{
		10: "12",
		12: "green",
	}

	names, err := os.OpenFile("test_name_encode.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Error(err)
	}
	defer names.Close()

	encoder := json.NewEncoder(names)
	err = encoder.Encode(blue)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadKeys(t *testing.T) {
	var blue map[uint32]string
	names, err := os.OpenFile("test_name_encode.json", os.O_CREATE, os.ModePerm)
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
}

func TestSaveData(t *testing.T) {
	order := binary.BigEndian

	red := []Database.Sale{}

	save, err := os.OpenFile("test_save.red", os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer save.Close()

	bs := make([]byte, 22*len(red))
	for i, x := range red {
		c := (22 * i)

		order.PutUint32(bs[c:c+4], x.ID)
		order.PutUint16(bs[c+4:c+6], x.Quantity)
		order.PutUint64(bs[c+6:c+14], math.Float64bits(x.Price))
		order.PutUint64(bs[c+14:c+22], math.Float64bits(x.Cost))
	}

	_, err = save.Write(bs)
	if err != nil {
		t.Error(err)
	}
}

func TestLoadData(t *testing.T) {
	order := binary.BigEndian

	buf, err := ioutil.ReadFile("test_save.red")
	if err != nil {
		t.Error(err)
	}

	black := make([]Database.Sale, len(buf)/22)

	for i := range black {
		c := 22 * i
		black[i].ID = order.Uint32(buf[c : c+4])
		black[i].Quantity = order.Uint16(buf[c+4 : c+6])
		black[i].Price = math.Float64frombits(order.Uint64(buf[c+6 : c+14]))
		black[i].Cost = math.Float64frombits(order.Uint64(buf[c+14 : c+22]))
	}

	t.Log(black)
}

//INTEGRATION TESTS
func GenerateTestData(t *testing.T) {
	//Create Sale List
	database := []Database.Sale{
		{ID: 121},
		{ID: 491289},
		{ID: 101983},
		{ID: 128001},
		{ID: 410283},
		{ID: 184098},
		{ID: 194820},
		{ID: 222222},
		{ID: 492780},
	}

	//Create Keys from Sale List
	keys := map[uint32]string{}

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

	bs := make([]byte, 22*len(database))
	for i, x := range database {
		c := (22 * i)

		order.PutUint32(bs[c:c+4], x.ID)
		order.PutUint16(bs[c+4:c+6], x.Quantity)
		order.PutUint64(bs[c+6:c+14], math.Float64bits(x.Price))
		order.PutUint64(bs[c+14:c+22], math.Float64bits(x.Cost))
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

func ReadTestData(t *testing.T) ([]Database.Sale, map[uint32]string) {
	//Read Sales
	order := binary.BigEndian

	buf, err := ioutil.ReadFile("test_integration.red")
	if err != nil {
		t.Error(err)
	}

	black := make([]Database.Sale, len(buf)/22)

	for i := range black {
		c := 22 * i
		black[i].ID = order.Uint32(buf[c : c+4])
		black[i].Quantity = order.Uint16(buf[c+4 : c+6])
		black[i].Price = math.Float64frombits(order.Uint64(buf[c+6 : c+14]))
		black[i].Cost = math.Float64frombits(order.Uint64(buf[c+14 : c+22]))
	}

	t.Log(black)

	//Read Keys
	var blue map[uint32]string
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

	//Compare Keys and sales
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
