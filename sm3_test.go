package sm3

import (
	"fmt"
	"testing"
)

func TestSM3_1(t *testing.T) {
	msg := "abc"
	var buffer [3]byte
	trueVal := "66c7f0f462eeedd9d1f2d46bdc10e4e24167c4875cf2f7a2297da02b8f4ba8e0"

	copy(buffer[:], msg)
	hw := NewSM3()
	hw.Write(buffer[:])

	uhash := make([]uint8, 32)
	hw.Sum(uhash[:0])
	calcVal := Byte2String(uhash)

	if calcVal != trueVal {
		t.Errorf("expected: %x,\nbut got: %v\n", trueVal, calcVal)
	}
}

func TestSM3_2(t *testing.T) {
	msg := "abcd"
	var buffer [4]byte
	trueVal := "debe9ff92275b8a138604889c18e5a4d6fdb70e5387e5765293dcba39c0c5732"

	copy(buffer[:], msg)
	hw := NewSM3()
	for i := 0; i < 16; i++ {
		hw.Write(buffer[:])
	}

	uhash := make([]uint8, 32)
	hw.Sum(uhash[:0])
	calcVal := Byte2String(uhash)

	if calcVal != trueVal {
		t.Errorf("expected: %x,\nbut got: %v\n", trueVal, calcVal)
	}
}

func Byte2String(b []byte) string {
	ret := ""
	for i := 0; i < len(b); i++ {
		ret += fmt.Sprintf("%02x", b[i])
	}
	return ret
}

func printByteSlice(list []byte) {
	for i := 0; i < len(list); i++ {
		fmt.Printf("%02x", list[i])
		if i%4 == 3 {
			fmt.Printf(" ")
		}
		if i%(4*8) == 31 {
			fmt.Println("")
		}
	}
	fmt.Println("")
}
