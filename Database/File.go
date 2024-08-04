package Database

import (
	"io"
	"io/ioutil"

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

var Items = map[uint16]*Entry{}
var Sales []Sale

var Customers = []string{}

const (
	ONCE uint8 = iota
	MONTHLY
	YEARLY
)

var file_names = [...]string{
	"Item_Reference.red",
	"Report_Data.red",
	"Customers.red",
	"Usrs.red",
	"BackUp_Map.red",
	"BackUp_Sales.red",
	"BackUp.red",
}

const (
	item_ref = iota
	report
	customers
	users
	backMap
	backSales
	backup
)

func DataInit() error {
	for i := 0; i < 7; i++ {
		save, err := fyne.CurrentApp().Storage().Create(file_names[i])
		if err == nil {
			save.Close()
		}
	}
	return nil
}

func SaveData() error {
	db, err := fyne.CurrentApp().Storage().Save(file_names[item_ref])
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Write(save_kv())
	if err != nil {
		return err
	}

	save, err := fyne.CurrentApp().Storage().Save(file_names[report])
	if err != nil {
		return err
	}

	_, err = save.Write(save_sales())
	save.Close()
	if err != nil {
		return err
	}

	customerFile, err := fyne.CurrentApp().Storage().Save(file_names[customers])
	if err != nil && err != io.EOF {
		return err
	}

	_, err = customerFile.Write(save_customers())
	customerFile.Close()
	if err != nil {
		return err
	}

	return nil
}

func LoadData() error {
	f, err := fyne.CurrentApp().Storage().Open(file_names[item_ref])
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	f.Close()

	load_kv(buf)

	reportFile, err := fyne.CurrentApp().Storage().Open(file_names[report])
	if err != nil {
		return err
	}

	buf, err = io.ReadAll(reportFile)
	if err != nil {
		return err
	}
	reportFile.Close()

	load_sales(buf)

	customerFile, err := fyne.CurrentApp().Storage().Open(file_names[customers])
	if err != nil && err != io.EOF {
		return err
	}

	custmBuf, err := io.ReadAll(customerFile)
	customerFile.Close()
	if err != nil && err != io.EOF {
		return err
	}

	load_customers(custmBuf)

	return nil
}

func SaveBackUp() error {
	save, err := fyne.CurrentApp().Storage().Save(file_names[backup])
	if err != nil {
		return err
	}

	_, err = save.Write(save_kv())
	if err != nil {
		return err
	}
	save.Close()

	save, err = fyne.CurrentApp().Storage().Save(file_names[backSales])
	if err != nil {
		return err
	}

	_, err = save.Write(save_sales())
	if err != nil {
		return err
	}
	save.Close()

	names, err := fyne.CurrentApp().Storage().Save(file_names[backMap])
	if err != nil {
		return err
	}
	defer names.Close()

	_, err = names.Write(saveBackupMap())
	return err
}

func LoadBackUp() error {
	file, err := fyne.CurrentApp().Storage().Open(file_names[backup])
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(file)
	file.Close()

	if err != nil {
		return err
	}

	load_kv(buf)

	file, err = fyne.CurrentApp().Storage().Open(file_names[backSales])
	if err != nil {
		return err
	}

	buf, err = ioutil.ReadAll(file)
	file.Close()

	if err != nil {
		return err
	}

	load_sales(buf)

	names, err := fyne.CurrentApp().Storage().Open(file_names[backMap])
	if err != nil {
		return err
	}
	defer names.Close()

	raw, err := io.ReadAll(names)
	if err != nil && err != io.EOF {
		return err
	}

	loadBackUpMap(raw)

	return nil
}
