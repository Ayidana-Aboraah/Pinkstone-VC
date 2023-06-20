package Database

import (
	"encoding/binary"
	"math"
	"strings"
)

const DATA_SIZE = 15
const ITEM_DS = 2

func save_users() (result []byte) {
	for _, v := range Users {
		result = append(result, []byte(v+"\n")...)
	}
	return
}

func load_users(buf []byte) {
	Users = strings.Split(string(buf), "\n")
	Users = Users[:len(Users)-1]
}

func save_report(data []Sale, order binary.ByteOrder) (result []byte) {
	result = make([]byte, DATA_SIZE*len(data))
	for i, x := range data {
		c := (DATA_SIZE * i)

		result[c] = x.Year
		result[c+1] = x.Month
		result[c+2] = x.Day
		result[c+3] = x.Usr

		order.PutUint16(result[c+4:c+6], x.Quantity)
		order.PutUint32(result[c+6:c+10], math.Float32bits(x.Price))
		PutUint40(result[c+10:c+DATA_SIZE], x.ID)
	}

	return
}

func load_report(buf []byte, order binary.ByteOrder) (report []Sale) {
	report = make([]Sale, len(buf)/DATA_SIZE)
	for i := range report {
		c := DATA_SIZE * i

		report[i].Year = uint8(buf[c])
		report[i].Month = uint8(buf[c+1])
		report[i].Day = uint8(buf[c+2])
		report[i].Usr = uint8(buf[c+3])

		report[i].Quantity = order.Uint16(buf[c+4 : c+6])
		report[i].Price = math.Float32frombits(order.Uint32(buf[c+6 : c+10]))
		report[i].ID = FromUint40(buf[c+10 : c+DATA_SIZE])
	}
	return
}

func save_expense(order binary.ByteOrder) (result []byte) {
	for _, v := range Expenses {
		buff := make([]byte, 8)
		buff[0] = v.Frequency
		buff[1] = v.Date[0]
		buff[2] = v.Date[1]
		buff[3] = v.Date[2]
		order.PutUint32(buff[4:], math.Float32bits(v.Amount))
		if v.Name == "" {
			v.Name = "_"
		}
		buff = append(buff, []byte(v.Name+"\n\n")...)
		result = append(result, buff...)
	}
	return
}

func load_expense(buf []byte, order binary.ByteOrder) {
	raw := strings.Split(string(buf), "\n\n")
	for _, v := range raw[:len(raw)-1] {
		var expense Expense
		process := []byte(v)
		expense.Frequency = process[0]
		expense.Date[0] = process[1]
		expense.Date[1] = process[2]
		expense.Date[2] = process[3]
		expense.Amount = math.Float32frombits(order.Uint32(process[4:8])) //TODO: Look into using a float16 to save space
		expense.Name = v[8:]
		if expense.Name == "_" {
			expense.Name = ""
		}
		Expenses = append(Expenses, expense)
	}
}

const kvSize = 11

func save_kv(order binary.ByteOrder) (result []byte) {
	sect := make([]byte, len(ItemKeys)*kvSize)
	names := ""
	var i int32

	for k, v := range ItemKeys {
		// ID[5], Quantity[2], Price[4]
		order.PutUint16(sect[i:i+2], v.Quantity)
		order.PutUint32(sect[i+2:i+6], math.Float32bits(v.Price))
		PutUint40(sect[i+6:i+kvSize], k)
		i += kvSize
		names += v.Name + "\n"
	}
	names = "\n\n" + names
	result = append(sect, names...)
	return
}

func load_kv(buf []byte, order binary.ByteOrder) {
	sides := strings.SplitN(string(buf), "\n\n", 2)

	names := strings.Split(sides[1], "\n")
	names = names[:len(names)-1]

	var c int
	for i := 0; i < len(sides[0])/kvSize; i++ {
		c = i * kvSize
		ItemKeys[FromUint40(buf[c+6:c+kvSize])] = &ItemEV{Quantity: order.Uint16(buf[c : c+2]), Price: math.Float32frombits(order.Uint32(buf[c+2 : c+6])), Name: names[i]}
	}
}
