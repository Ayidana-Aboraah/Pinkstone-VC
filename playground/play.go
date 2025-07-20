package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Muslim struct {
	ID     uint8
	Houses []float64
	Wives  [4]string
}

type Mumin struct {
	ID     uint8
	Houses []float64
	Wives  [4]string
	Belief float64
}

func main() {
	man := Mumin{ID: 5, Houses: []float64{5554, 914.341, 5164.1534134}, Wives: [4]string{"Fatima", "Aisha", "Khadija", "Becky"}, Belief: 9999.444}
	men := []Mumin{ // TODO: Test out encoding
		man,
		Mumin{ID: 5, Houses: []float64{0, 0, 0}, Wives: [4]string{"Fatima", "Aisha", "Khadija", "Becky"}},
		Mumin{ID: 5, Houses: []float64{0, 0, 0}, Wives: [4]string{"Fatima", "Aisha", "Khadija", "Becky"}},
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(men)
	if err != nil {
		log.Fatal("Error encoding user data:", err)
	}

	fmt.Println("Data: ", men)

	var decoded_man []Muslim
	decoder := gob.NewDecoder(&buf)
	err = decoder.Decode(&decoded_man)
	if err != nil {
		log.Fatal("Error decoding user data:", err)
	}
	fmt.Println("Decoded Data: %+v\n", decoded_man)
}
