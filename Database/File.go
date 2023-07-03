package Database

import (
	"io"
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
)

type Entry struct {
	Price    float32
	Cost     [3]float32
	Quantity [3]float32
	Name     string
}

type Sale struct {
	Year, Month, Day, Usr, Customer uint8
	ID                              uint16
	Price, Cost, Quantity           float32
}

var Item = map[uint16]*Entry{}
var Report []Sale

var Customers = []string{}
var Users = []string{}
var Current_User uint8

const (
	ONCE uint8 = iota
	MONTHLY
	YEARLY
)

func DataInit(remove bool) error {
	for i, file := 0, ""; i < 6; i++ {
		switch i {
		case 0:
			file = "Item_Reference.red"
		case 1:
			file = "Report_Data.red"
		case 2:
			file = "Customers.red"
		case 3:
			file = "Usrs.red"
		case 4:
			file = "BackUp_Map.red"
		case 5:
			file = "BackUp.red"
		}

		if !remove {
			save, err := fyne.CurrentApp().Storage().Create(file)
			if err == nil {
				save.Close()
			}
		} else {
			err := fyne.CurrentApp().Storage().Remove(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func SaveData() error {
	db, err := fyne.CurrentApp().Storage().Save("Item_Reference.red")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Write(save_kv())
	if err != nil {
		return err
	}

	save, err := fyne.CurrentApp().Storage().Save("Report_Data.red")
	if err != nil {
		return err
	}

	_, err = save.Write(save_report(Report))
	save.Close()
	if err != nil {
		return err
	}

	customerFile, err := fyne.CurrentApp().Storage().Save("Customers.red")
	if err != nil && err != io.EOF {
		return err
	}

	_, err = customerFile.Write(save_customers())
	customerFile.Close()
	if err != nil {
		return err
	}

	usrFile, err := fyne.CurrentApp().Storage().Save("Usrs.red")
	if err != nil && err != io.EOF {
		return err
	}

	_, err = usrFile.Write(save_users())
	usrFile.Close()
	if err != nil {
		return err
	}

	return err
}

func LoadData() error {
	f, err := fyne.CurrentApp().Storage().Open("Item_Reference.red")
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	f.Close()

	load_kv(buf)

	reportFile, err := fyne.CurrentApp().Storage().Open("Report_Data.red")
	if err != nil {
		return err
	}

	buf, err = io.ReadAll(reportFile)
	if err != nil {
		return err
	}
	reportFile.Close()

	load_report(buf)

	customerFile, err := fyne.CurrentApp().Storage().Open("Customers.red")
	if err != nil && err != io.EOF {
		return err
	}

	custmBuf, err := io.ReadAll(customerFile)
	customerFile.Close()
	if err != nil && err != io.EOF {
		return err
	}

	load_customers(custmBuf)

	usrFile, err := fyne.CurrentApp().Storage().Open("Usrs.red")
	if err != nil && err != io.EOF {
		return err
	}

	usrBytes, err := io.ReadAll(usrFile)
	usrFile.Close()
	if err != nil && err != io.EOF {
		return err
	}

	load_users(usrBytes)

	return nil
}

func SaveBackUp() error {
	save, err := fyne.CurrentApp().Storage().Save("BackUp.red")
	if err != nil {
		return err
	}

	var BackUp_Buff []byte

	BackUp_Buff = append(BackUp_Buff, save_kv()...)
	BackUp_Buff = append(BackUp_Buff, []byte{10, 10, 10}...)

	BackUp_Buff = append(BackUp_Buff, save_report(Report)...)
	BackUp_Buff = append(BackUp_Buff, []byte{10, 10, 10}...)

	_, err = save.Write(BackUp_Buff)
	if err != nil {
		return err
	}
	save.Close()

	names, err := fyne.CurrentApp().Storage().Save("BackUp_Map.red")
	if err != nil {
		return err
	}
	defer names.Close()

	buf := append(save_kv(), "\n\n"...)
	buf = append(buf, save_customers()...)

	_, err = names.Write(buf)
	return err
}

func LoadBackUp() error {
	file, err := fyne.CurrentApp().Storage().Open("BackUp.red")
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(file)
	file.Close()

	if err != nil {
		return err
	}

	black := strings.Split(string(buf), "\n\n\n")

	load_kv([]byte(black[0]))
	load_report([]byte(black[1]))

	names, err := fyne.CurrentApp().Storage().Open("BackUp_Map.red")
	if err != nil {
		return err
	}
	defer names.Close()

	raw, err := io.ReadAll(names)
	if err != nil && err != io.EOF {
		return err
	}

	maps := strings.SplitN(string(raw), "\n\n", 2)

	load_users([]byte(maps[0])) // NOTE: Watch for odd activity | [2:] works the same as [], so may be something ups
	load_customers([]byte(maps[1]))

	return nil
}
