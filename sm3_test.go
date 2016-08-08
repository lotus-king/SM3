package sm3

import "testing"

func TestSM3_1(t *testing.T) {
	msg := "abc"
	var buffer [3]byte
	trueVal := "66c7f0f462eeedd9d1f2d46bdc10e4e24167c4875cf2f7a2297da02b8f4ba8e0"

	copy(buffer[:], msg)
	sm3 := NewSM3(buffer[:])
	calcVal := sm3.Hash()
	if calcVal != trueVal {
		t.Errorf("Expected: %v,\nbut : %v\n", trueVal, calcVal)
	}
}

func TestSM3_2(t *testing.T) {
	str := "abcd"
	msg := ""
	for i := 0; i < 16; i++ {
		msg += str
	}

	trueVal := "debe9ff92275b8a138604889c18e5a4d6fdb70e5387e5765293dcba39c0c5732"

	var buffer [64]byte
	copy(buffer[:], msg)
	sm3 := NewSM3(buffer[:])
	calcVal := sm3.Hash()
	if calcVal != trueVal {
		t.Errorf("Expected: %v,\nbut : %v\n", trueVal, calcVal)
	}
}

// func TestSM3(t *testing.T) {
// 	str := "abcd"
// 	msg := ""
// 	for i := 0; i < 16; i++ {
// 		msg += str
// 	}

// 	var buffer [64]byte
// 	copy(buffer[:], msg)
// 	sm3 := NewSM3(buffer[:])
// 	sm3.Hash()
// }
