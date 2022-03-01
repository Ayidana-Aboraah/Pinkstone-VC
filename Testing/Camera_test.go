package Test

import (
	"BronzeHermes/Cam"
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
