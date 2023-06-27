package Database

import (
	"encoding/binary"
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
	Year, Month, Day, Usr uint8
	ID                    uint16
	Price, Cost, Quantity float32
}

var Item = map[uint16]*Entry{}
var Reports [2][]Sale
var Free_Spaces []int

var Users = []string{}
var Current_User uint8

const (
	ONCE uint8 = iota
	MONTHLY
	YEARLY
)

func PutUint40(b []byte, v uint64) {
	_ = b[4]
	b[0] = byte(v >> 32)
	b[1] = byte(v >> 24)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 8)
	b[4] = byte(v)
}

func FromUint40(b []byte) uint64 {
	return uint64(b[4]) | uint64(b[3])<<8 | uint64(b[2])<<16 | uint64(b[1])<<24 | uint64(b[0])<<32
}

func DataInit(remove bool) error {
	for i, file := 0, ""; i < 6; i++ {
		switch i {
		case 0:
			file = "Item_Reference.red"
		case 1:
			file = "Report_Data.red"
		case 2:
			file = "Price_Log.red"
		case 3:
			file = "BackUp.red"
		case 4:
			file = "BackUp_Map.red"
		case 5:
			file = "Usrs.red"
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
	order := binary.BigEndian
	file := "Item_Reference.red"

	db, err := fyne.CurrentApp().Storage().Save(file)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Write(save_kv(order))
	if err != nil {
		return err
	}

	for idx, database := range Reports {
		switch idx {
		case 0:
			file = "Report_Data.red"
		case 1:
			file = "Price_Log.red"
		}

		save, err := fyne.CurrentApp().Storage().Save(file)
		if err != nil {
			return err
		}

		_, err = save.Write(save_report(database, order))
		save.Close()
		if err != nil {
			return err
		}
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
	order := binary.BigEndian
	file := "Item_Reference.red"

	f, err := fyne.CurrentApp().Storage().Open(file)
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	f.Close()

	load_kv(buf, order)

	for idx := range Reports {
		switch idx {
		case 0:
			file = "Report_Data.red"
		case 1:
			file = "Price_Log.red"
		}

		file, err := fyne.CurrentApp().Storage().Open(file)
		if err != nil {
			return err
		}

		buf, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		file.Close()

		Reports[idx] = load_report(buf, order)
	}

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
	order := binary.BigEndian

	save, err := fyne.CurrentApp().Storage().Save("BackUp.red")
	if err != nil {
		return err
	}

	var BackUp_Buff []byte

	BackUp_Buff = append(BackUp_Buff, save_kv(order)...)
	BackUp_Buff = append(BackUp_Buff, []byte{10, 10, 10}...)

	BackUp_Buff = append(BackUp_Buff, save_report(Reports[0], order)...)
	BackUp_Buff = append(BackUp_Buff, []byte{10, 10, 10}...)

	BackUp_Buff = append(BackUp_Buff, save_report(Reports[1], order)...)

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

	_, err = names.Write(save_kv(order))
	return err
}

func LoadBackUp() error {
	order := binary.BigEndian

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

	load_kv([]byte(black[0]), order)
	Reports[0] = load_report([]byte(black[1]), order)
	Reports[1] = load_report([]byte(black[2]), order)

	names, err := fyne.CurrentApp().Storage().Open("BackUp_Map.red")
	if err != nil {
		return err
	}
	defer names.Close()

	raw, err := io.ReadAll(names)
	if err != nil && err != io.EOF {
		return err
	}

	load_kv([]byte(raw), order) // NOTE: Watch for odd activity | [2:] works the same as [], so may be something ups

	return nil
}
