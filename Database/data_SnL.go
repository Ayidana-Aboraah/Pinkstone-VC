package Database

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

const DATA_SIZE = 18
const ITEM_DS = 6

func save_itemDB(order binary.ByteOrder) (result []byte) {
	result = make([]byte, len(Items)*ITEM_DS)
	for i := range Items {
		c := i * ITEM_DS
		order.PutUint16(result[c:c+2], Items[i].Quantity)
		order.PutUint32(result[c+2:c+ITEM_DS], math.Float32bits(Items[i].Cost))
	}
	return
}

func load_itemDB(buf []byte, order binary.ByteOrder) {
	Items = make([]Item, len(buf)/ITEM_DS)

	for i := range Items {
		c := ITEM_DS * i
		Items[i].Quantity = order.Uint16(buf[c : c+2])
		Items[i].Cost = math.Float32frombits(order.Uint32(buf[c+2 : c+ITEM_DS]))
	}
}

func save_report(data []Sale, order binary.ByteOrder) (result []byte) {
	result = make([]byte, DATA_SIZE*len(data))
	for i, x := range data {
		c := (DATA_SIZE * i)

		result[c] = x.Year
		result[c+1] = x.Month
		result[c+2] = x.Day

		order.PutUint16(result[c+3:c+5], x.Quantity)
		order.PutUint32(result[c+5:c+9], math.Float32bits(x.Price))
		order.PutUint32(result[c+9:c+13], math.Float32bits(x.Cost))
		PutUint40(result[c+13:c+DATA_SIZE], x.ID)
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

		report[i].Quantity = order.Uint16(buf[c+3 : c+5])
		report[i].Price = math.Float32frombits(order.Uint32(buf[c+5 : c+9]))
		report[i].Cost = math.Float32frombits(order.Uint32(buf[c+9 : c+13]))
		report[i].ID = FromUint40(buf[c+13 : c+DATA_SIZE])
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

func save_kv(order binary.ByteOrder) (result []byte) {
	for k, v := range ItemKeys {
		mine := make([]byte, 9+(len(v.Idxes)*2))
		PutUint40(mine[:5], k)
		order.PutUint32(mine[5:9], math.Float32bits(v.Price))

		for i := range v.Idxes {
			c := 9 + (i * 2)
			order.PutUint16(mine[c:c+2], uint16(v.Idxes[i]))
		}
		mine = append(mine, []byte{255, 255}...)
		mine = append(mine, []byte(v.Name+"\n\n")...)
		result = append(result, mine...)
	}
	return
}

func load_kv(buf []byte, order binary.ByteOrder) {
	entries := strings.Split(string(buf), "\n\n")

	for _, v := range entries[:len(entries)-1] {
		data, name, found := strings.Cut(v, string([]byte{255, 255}))
		if !found {
			fmt.Println("Seems as if the character does not exist") // NOTE: return a error
		}
		ItemKeys[FromUint40([]byte(data))] = &ItemEV{
			Price: math.Float32frombits(order.Uint32([]byte(data)[5:9])),
			Name:  name,
		}

		for i := 0; i < len([]byte(data)[9:])/2; i++ {
			c := 9 + (i * 2)
			ItemKeys[FromUint40([]byte(data))].Idxes = append(ItemKeys[FromUint40([]byte(data))].Idxes, int(order.Uint16([]byte(data)[c:c+2])))
		}
	}
}
