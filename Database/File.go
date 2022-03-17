package Database

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"math"
	"os"

	"fyne.io/fyne/v2/storage"
)

var NameKeys map[uint64]string

// var Items []Sale
// var ReportData []Sale
// var PriceLog []Sale

//0 Items; 1 ReportData; 2 PriceLog
var Databases [3][]Sale

type Sale struct {
	Year     uint8
	Month    uint8
	Day      uint8
	Quantity uint16
	Price    float32
	Cost     float32
	ID       uint64
}

func InitCheck() error {
	uri, err := storage.ParseURI("file://com.android.externalstorage.documents/documents/primary%3AAndroid%2ABH_Saves")
	if err != nil {
		return err
	}

	if check, err := storage.Exists(uri); !check {
		if err != nil {
			return err
		}
		err = storage.CreateListable(uri)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveData() error {
	order := binary.BigEndian
	var file string

	for idx, database := range Databases {
		switch idx {
		case 0:
			file = "Item_Reference.red"
		case 1:
			file = "Report_Data.red"
		case 2:
			file = "Price_Log.red"
		}

		// func() bool {
		// 	uri, err := storage.ParseURI("Android/obb/com.redstoneagx.bronzehermes")
		// 	if err != nil {
		// 		fmt.Println(err)
		// 	}

		// 	if check, err := storage.Exists(uri); !check {
		// 		fmt.Println(check)
		// 		if err != nil {
		// 			fmt.Println(err)
		// 		}
		// 		return check
		// 	}
		// 	return true
		// }()

		uri, err := storage.ParseURI("file://com.android.externalstorage.documents/documents/primary%3AAndroid%2ABH_Saves%2F" + file)
		if err != nil {
			return err
		}

		save, err := storage.Writer(uri)
		if err != nil {
			return err
		}

		// save, err := os.OpenFile("Saves/"+file, os.O_CREATE, os.ModePerm)
		// if err != nil {
		// 	return err
		// }

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
		save.Close()

		if err != nil {
			return err
		}
	}

	err := func() error {
		uri, err := storage.ParseURI("file://com.android.externalstorage.documents/documents/primary%3AAndroid%2ABH_Saves%F2name_keys.json")
		if err != nil {
			return err
		}
		names, err := storage.Writer(uri)

		// names, err := os.OpenFile("Saves/name_keys.json", os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer names.Close()

		encoder := json.NewEncoder(names)
		encoder.Encode(NameKeys)
		return nil
	}()

	return err
}

func LoadData() error {
	var file string
	order := binary.BigEndian

	for idx := range Databases {
		switch idx {
		case 0:
			file = "Item_Reference.red"
		case 1:
			file = "Report_Data.red"
		case 2:
			file = "Price_Log.red"
		}

		uri, err := storage.ParseURI("file://com.android.externalstorage.documents/documents/primary%3AAndroid%2ABH_Saves%2F" + file)
		if err != nil {
			return err
		}

		file, err := storage.Reader(uri)
		if err != nil {
			return err
		}

		buf, err := ioutil.ReadAll(file)
		file.Close()

		if err != nil {
			return err
		}

		black := make([]Sale, len(buf)/21)

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

		Databases[idx] = black
	}

	err := func() error {
		uri, err := storage.ParseURI("file://com.android.externalstorage.documents/documents/primary%3AAndroid%2ABH_Saves%2Fname_keys.json")
		names, err := storage.Reader(uri)
		// names, err := os.OpenFile("Saves/name_keys.json", os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer names.Close()

		encoder := json.NewDecoder(names)
		err = encoder.Decode(&NameKeys)
		if err != nil {
			return err
		}
		return nil
	}()

	return err
}

func BackUpAllData() error {
	err := func() error {
		names, err := os.OpenFile("file://com.android.externalstorage.documents/documents/primary%3AAndroid%2ABH_Saves%2FBackup_Keys.json", os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer names.Close()

		encoder := json.NewEncoder(names)
		err = encoder.Encode(NameKeys)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return err
	}

	order := binary.BigEndian

	// save, err := os.OpenFile("Saves/BackUp.red", os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	return err
	// }
	uri, err := storage.ParseURI("file://com.android.externalstorage.documents/documents/primary%3AAndroid%2ABH_Saves%2FBackUp.red")
	if err != nil {
		return err
	}

	save, err := storage.Writer(uri)
	if err != nil {
		return err
	}

	defer save.Close()

	bs := make([]byte, ((21 * len(Databases[0])) + (21 * len(Databases[1])) + (21 * len(Databases[2]))))

	previousLength := 0
	for _, database := range Databases {
		for i, x := range database {
			initial := (previousLength * 21) + (21 * i)

			bs[initial] = x.Year
			bs[initial+1] = x.Month
			bs[initial+2] = x.Day

			order.PutUint16(bs[initial+3:initial+5], x.Quantity)
			order.PutUint32(bs[initial+5:initial+9], math.Float32bits(x.Price))
			order.PutUint32(bs[initial+9:initial+13], math.Float32bits(x.Cost))
			order.PutUint64(bs[initial+13:initial+21], x.ID)
		}
		previousLength += len(database)
	}

	_, err = save.Write(bs)
	return err
}

func LoadBackUp() error {
	order := binary.BigEndian

	uri, err := storage.ParseURI("file://com.android.externalstorage.documents/documents/primary%3AAndroid%2ABH_Saves%2FBackUp.red")
	if err != nil {
		return err
	}

	file, err := storage.Reader(uri)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(file)
	file.Close()

	// buf, err := ioutil.ReadFile("Saves/BackUp.red")
	if err != nil {
		return err
	}

	black := make([]Sale, len(buf)/21)

	for i, v := range black {
		c := 21 * i

		black[i].Year = uint8(buf[c])
		black[i].Month = uint8(buf[c+1])
		black[i].Day = uint8(buf[c+2])

		black[i].Quantity = order.Uint16(buf[c+3 : c+5])
		black[i].Price = math.Float32frombits(order.Uint32(buf[c+5 : c+9]))
		black[i].Cost = math.Float32frombits(order.Uint32(buf[c+9 : c+13]))
		black[i].ID = order.Uint64(buf[c+13 : c+21])

		if v.Year == 0 {
			Databases[0] = append(Databases[0], v)
		} else if v.Quantity == 0 {
			Databases[2] = append(Databases[2], v)
		} else {
			Databases[1] = append(Databases[1], v)
		}
	}

	err = func() error {
		uri, err := storage.ParseURI("file://com.android.externalstorage.documents/documents/primary%3AAndroid%2ABH_Saves%2FBackup_Keys.json")
		names, err := storage.Reader(uri)
		// names, err := os.OpenFile("Saves/Backup_Keys.json", os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
		defer names.Close()

		encoder := json.NewDecoder(names)
		err = encoder.Decode(&NameKeys)
		if err != nil {
			return err
		}
		return nil
	}()

	return err
}
