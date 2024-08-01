package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(time.Unix(time.Now().Local().Unix(), 0))
}
