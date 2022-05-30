package test

import (
	"BronzeHermes/Database"
	"testing"
)

func TestToUint40(t *testing.T) {
	value := 674398202423
	buf := make([]byte, 5)
	Database.ToUint40(buf, uint64(value))

	newVal := Database.FromUint40(buf)

	if value != int(newVal) {
		t.Errorf("Values Don't match | Value: %v, New Value: %v", value, newVal)
	}
	t.Log(value)
	t.Log(newVal)
}
