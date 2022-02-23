package Test

import (
	"BronzeHermes/Cam"
	"BronzeHermes/Database"
	"encoding/binary"
	"encoding/json"
	"image"
	_ "image/jpeg"
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

	red := []Database.Sale{
		{ID: 012345, Year: 100, Day: 12, Month: 6, Price: 20.44},
	}

	save, err := os.OpenFile("test_save.red", os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer save.Close()

	bs := make([]byte, 17*len(red))
	for i, x := range red {
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

	black := make([]Database.Sale, len(buf)/17)

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

	t.Log(black)
}

func TestBlue(t *testing.T) {
	file, _ := os.Open("TestCode.jpg")
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		t.Error(err)
	}

	res := Cam.ReadImage(img)
	//Change to log the uint
	t.Log(res.GetText())
}
