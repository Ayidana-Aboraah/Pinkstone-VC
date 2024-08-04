package Database

import (
	"encoding/binary"
	"math"
	"strings"
)

func save_customers() (result []byte) {
	if len(Customers) > 1 {
		for _, v := range Customers[1:] {
			result = append(result, v+"\n"...)
		}
	}
	return
}

func load_customers(buf []byte) {
	Customers = []string{""}
	Customers = append(Customers, strings.Split(string(buf), "\n")...)
	Customers = Customers[:len(Customers)-1]
}

func SaveNLoadCustomers() {
	load_customers(save_customers())
}

func save_sales() (result []byte) {
	result = make([]byte, binary.Size(Sale{})*len(Sales))
	for i, x := range Sales {
		c := (binary.Size(Sale{}) * i)

		order.PutUint64(result[c:c+8], uint64(x.Timestamp))
		order.PutUint16(result[c+8:c+10], x.Customer)

		order.PutUint16(result[c+10:c+12], x.ID)
		order.PutUint32(result[c+12:c+16], math.Float32bits(x.Price))
		order.PutUint32(result[c+16:c+20], math.Float32bits(x.Cost))
		order.PutUint32(result[c+20:c+24], math.Float32bits(x.Quantity))
	}

	return
}

func load_sales(buf []byte) {
	sales := make([]Sale, len(buf)/binary.Size(Sale{}))
	for i := range sales {
		c := binary.Size(Sale{}) * i

		sales[i].Timestamp = int64(order.Uint64(buf[c : c+8]))
		sales[i].Customer = uint16(order.Uint16(buf[c+8 : c+10]))

		sales[i].ID = order.Uint16(buf[c+10 : c+12])
		sales[i].Price = math.Float32frombits(order.Uint32(buf[c+12 : c+16]))
		sales[i].Cost = math.Float32frombits(order.Uint32(buf[c+16 : c+20]))
		sales[i].Quantity = math.Float32frombits(order.Uint32(buf[c+20 : c+24]))
	}
	Sales = sales
}

func SaveNLoadSales() {
	load_sales(save_sales())
}

const kvSize = 30

var order = binary.BigEndian

func save_kv() (result []byte) {
	sect := make([]byte, len(Items)*kvSize)
	names := ""
	var i int32

	for k, v := range Items {
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

		Items[order.Uint16(buf[c:c+2])] = &Entry{
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

func saveBackupMap() (buf []byte) {
	buf = append(buf, save_customers()...)
	return
}

func loadBackUpMap(buf []byte) {
	load_customers(buf)
}

func SaveNLoadBackUp() {
	load_kv(save_kv())
	load_sales(save_sales())
	loadBackUpMap(saveBackupMap())
}
