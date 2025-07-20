package Database

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"

	"fyne.io/fyne/v2"
)

type Entry struct {
	Price    float32
	Cost     [3]float32
	Quantity [3]float32
	Name     string
}

type Sale struct {
	Customer              uint16 // 0 is reserved for Damages
	ID                    uint16
	Price, Cost, Quantity float32
	Timestamp             int64
}

// var Items = map[uint16]*Entry{}
var Items []Entry
var DeadIDsQueue []int
var Sales []Sale

var Customers = []string{}

const (
	ONCE uint8 = iota
	MONTHLY
	YEARLY
)

var file_names = [...]string{
	"Items.red",
	"Report.red",
	"Customers.red",
	"Backup_Items.red",
	"Backup_Report.red",
	"Backup_Customers.red",
}

// func Init() {

// }

func Init(test bool) {
	// Shared mutable state of doom and destruction
	encoder = gob.NewEncoder(&buffer)
	decoder = gob.NewDecoder(&buffer)

	if !test {
		for _, name := range file_names {
			file, err := fyne.CurrentApp().Storage().Create(name)
			if err == nil {
				file.Close()
			}
		}
	}

}

var buffer bytes.Buffer //Items, Reports, Customers
var encoder *gob.Encoder
var decoder *gob.Decoder

func Save(writer io.Writer, data any) error {
	err := encoder.Encode(data)
	if err != nil {
		return err
	}

	fmt.Println(buffer.Bytes())
	_, err = buffer.WriteTo(writer)
	if err != nil {
		return err
	}

	fmt.Println(buffer.Bytes())
	return nil
}

func Load(reader io.Reader, output any) error {
	log.Println(buffer.Bytes())
	_, err := buffer.ReadFrom(reader)
	if err != nil {
		return err
	}
	log.Println(buffer.Bytes())

	err = decoder.Decode(output)
	if err != nil {
		return err
	}
	return nil
}

func SaveData(backup bool) error {

	var file fyne.URIWriteCloser
	var err error

	var fileNames []string
	if backup {
		fileNames = file_names[3:]
	} else {
		fileNames = file_names[:3]
	}

	data := [6]any{
		Items,
		Sales,
		Customers,
		Items,
		Sales,
		Customers,
	}

	for i, name := range fileNames {
		fmt.Println(name)
		fmt.Println(data[i])

		file, err = fyne.CurrentApp().Storage().Save(name)
		if err != nil && err != io.EOF {
			return err
		}
		defer file.Close()

		err = Save(file, data[i])
		if err != nil {
			return err
		}
	}

	fmt.Println("\n")

	return nil
}

func LoadData(backup bool) error {
	var file fyne.URIReadCloser
	var err error
	var fileNames []string
	if backup {
		fileNames = file_names[3:]
	} else {
		fileNames = file_names[:3]
	}

	data := [6]any{
		&Items,
		&Sales,
		&Customers,
		&Items,
		&Sales,
		&Customers,
	}

	for i, name := range fileNames {
		fmt.Println(name)

		file, err = fyne.CurrentApp().Storage().Open(name)
		if err != nil && err != io.EOF {
			return err
		}
		defer file.Close()

		err = Load(file, data[i])
		if err != nil {
			return err
		}
	}
	fmt.Println("\n")

	return nil
}
