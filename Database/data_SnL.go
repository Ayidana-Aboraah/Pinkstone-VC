package Database

import (
	"encoding/binary"
	"math"
	"strings"
)

func save_users() (result []byte) {
	for _, v := range Users {
		result = append(result, v+"\n"...)
	}
	return
}

func load_users(buf []byte) {
	Users = strings.Split(string(buf), "\n")
	Users = Users[:len(Users)-1]
}

func SaveNLoadUsers() {
	load_users(save_users())
}

func save_customers() (result []byte) {
	for _, v := range Customers {
		result = append(result, v+"\n"...)
	}
	return
}

func load_customers(buf []byte) {
	Customers = strings.Split(string(buf), "\n")
	Customers = Customers[:len(Customers)-1]
}

func SaveNLoadCustomers() {
	load_customers(save_customers())
}

func save_report(data []Sale) (result []byte) {
	order := binary.BigEndian
	result = make([]byte, 19*len(data))
	for i, x := range data {
		c := (19 * i)

		result[c] = x.Year
		result[c+1] = x.Month
		result[c+2] = x.Day
		result[c+3] = x.Usr
		result[c+4] = x.Customer

		order.PutUint16(result[c+5:c+7], x.ID)
		order.PutUint32(result[c+7:c+11], math.Float32bits(x.Price))
		order.PutUint32(result[c+11:c+15], math.Float32bits(x.Cost))
		order.PutUint32(result[c+15:c+19], math.Float32bits(x.Quantity))
	}

	return
}

func load_report(buf []byte) {
	order := binary.BigEndian
	report := make([]Sale, len(buf)/19)
	for i := range report {
		c := 19 * i

		report[i].Year = uint8(buf[c])
		report[i].Month = uint8(buf[c+1])
		report[i].Day = uint8(buf[c+2])
		report[i].Usr = uint8(buf[c+3])
		report[i].Customer = uint8(buf[c+4])

		report[i].ID = order.Uint16(buf[c+5 : c+7])
		report[i].Price = math.Float32frombits(order.Uint32(buf[c+7 : c+11]))
		report[i].Cost = math.Float32frombits(order.Uint32(buf[c+11 : c+15]))
		report[i].Quantity = math.Float32frombits(order.Uint32(buf[c+15 : c+19]))
	}
	Report = report
}

const kvSize = 30

func save_kv() (result []byte) {
	order := binary.BigEndian
	sect := make([]byte, len(Item)*kvSize)
	names := ""
	var i int32

	for k, v := range Item {
		// ID[2], Price [4], Cost[4*3], Quantity[4*3]
		order.PutUint16(sect[i:i+2], k)

		order.PutUint32(sect[i+2:i+6], math.Float32bits(v.Price))

		order.PutUint32(sect[i+6:i+10], math.Float32bits(v.Cost[0]))
		order.PutUint32(sect[i+10:i+14], math.Float32bits(v.Cost[1]))
		order.PutUint32(sect[i+14:i+18], math.Float32bits(v.Cost[2]))

		order.PutUint32(sect[i+18:i+22], math.Float32bits(v.Quantity[0]))
		order.PutUint32(sect[i+22:i+26], math.Float32bits(v.Quantity[1]))
		order.PutUint32(sect[i+26:i+kvSize], math.Float32bits(v.Quantity[2]))

		i += kvSize
		names += v.Name + "\n"
	}

	names = "\n\n" + names
	result = append(sect, names...)
	return
}

func load_kv(buf []byte) {
	order := binary.BigEndian
	if len(buf) == 0 {
		return
	}
	sides := strings.SplitN(string(buf), "\n\n", 2)

	names := strings.Split(sides[1], "\n")
	names = names[:len(names)-1]

	c := 0
	for i := 0; i < len(sides[0])/kvSize; i++ {

		Item[order.Uint16(buf[c:c+2])] = &Entry{
			Name:     names[i],
			Price:    math.Float32frombits(order.Uint32(buf[c+2 : c+6])),
			Cost:     [3]float32{math.Float32frombits(order.Uint32(buf[c+6 : c+10])), math.Float32frombits(order.Uint32(buf[c+10 : c+14])), math.Float32frombits(order.Uint32(buf[c+14 : c+18]))},
			Quantity: [3]float32{math.Float32frombits(order.Uint32(buf[c+18 : c+22])), math.Float32frombits(order.Uint32(buf[c+22 : c+26])), math.Float32frombits(order.Uint32(buf[c+26 : c+kvSize]))},
		}

		c += kvSize
	}
}

func SaveNLoadKV() {
	load_kv(save_kv())
}
