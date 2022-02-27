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
	blue := map[uint64]string{
		674398202423: "Piano_Book",
		490810512013: "verrr",
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
	var blue map[uint64]string
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
		{ID: 674398202423, Year: 021, Day: 1, Month: 12, Price: 24, Cost: 43, Quantity: 2},
	}

	if _, err := os.ReadDir("Beep"); err != nil {
		os.Mkdir("Beep", os.ModeDir)
	}

	save, err := os.OpenFile("Beep/test_save.red", os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Error(err)
	}

	defer save.Close()

	bs := make([]byte, 21*len(red))
	for i, x := range red {
		c := (21 * i)

		bs[c] = x.Year
		bs[c+1] = x.Month
		bs[c+2] = x.Day

		order.PutUint16(bs[c+3:c+5], x.Quantity)
		order.PutUint32(bs[c+5:c+9], math.Float32bits(x.Price))
		order.PutUint32(bs[c+9:c+13], math.Float32bits(x.Cost))
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

func TestBackUp(t *testing.T) {
	Items := []Database.Sale{}
	PriceLog := []Database.Sale{}
	ReportData := []Database.Sale{}

	err := Database.BackUpAllData()
	if err != nil {
		t.Error(err)
	}

	initial := [][]Database.Sale{Items, ReportData, PriceLog}

	Database.LoadBackUp()

	for e, database := range initial {
		for i, _ := range database {
			var current []Database.Sale
			switch e {
			case 0:
				current = Items
			case 1:
				current = ReportData
			case 2:
				current = PriceLog
			}

			if database[i] == current[i] {
				t.Errorf("%v doesn't match %v", database, current)
			}
		}
	}
}
