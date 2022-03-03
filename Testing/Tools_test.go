package Test

import (
	"BronzeHermes/Cam"
	"BronzeHermes/Database"
	"image"
	_ "image/jpeg"
	"os"
	"strconv"
	"testing"
)

func TestBarcodeReading(t *testing.T) {
	file, _ := os.Open("Assets/TestCode.jpg")
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		t.Error(err)
	}

	res := Cam.ReadImage(img)
	//Change to log the uint
	t.Log(res.GetText())
	a, err := strconv.Atoi(res.GetText())
	if err != nil {
		t.Error(a)
	}
}

func TestGenerateData(t *testing.T) {
	Database.Databases = [3][]Database.Sale{
		{
			{ID: 674398202423},
			{ID: 490810512013},
		},
		{
			{ID: 674398202423, Quantity: 2, Price: 12, Cost: 43, Year: 21, Month: 12, Day: 1},
			{ID: 674398202423, Quantity: 1, Price: 76, Cost: 43, Year: 21, Month: 12, Day: 1},
			{ID: 674398202423, Quantity: 2, Price: 0, Cost: 43, Year: 21, Month: 12, Day: 1},
			{ID: 490810512013, Quantity: 1, Price: 10, Cost: 1, Year: 21, Month: 12, Day: 2},
			{ID: 674398202423, Quantity: 2, Price: 5, Cost: 43, Year: 21, Month: 12, Day: 3},
			{ID: 674398202423, Quantity: 2, Price: 6, Cost: 43, Year: 21, Month: 12, Day: 3},
			{ID: 674398202423, Quantity: 1, Price: 0, Cost: 43, Year: 21, Month: 12, Day: 4},
			{ID: 490810512013, Quantity: 1, Price: 21, Cost: 0, Year: 21, Month: 12, Day: 4},
		},
		{},
	}

	Database.NameKeys = map[uint64]string{
		674398202423: "Piano_Book",
		490810512013: "Verando",
	}

	if err := Database.SaveData(); err != nil {
		t.Error(err)
	}
}
