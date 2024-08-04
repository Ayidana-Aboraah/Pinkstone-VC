package main

import (
	"fmt"
	"time"
)

type Sale struct {
	Customer              uint8
	ID                    uint16
	Price, Cost, Quantity float32
	timestamp             int64
}

func main() {
	fmt.Println(time.Now().Local().String())
}
