package sm3

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestSM3_2(t *testing.T) {
	msg := "abcd"
	var buffer [4]byte
	trueVal := "debe9ff92275b8a138604889c18e5a4d6fdb70e5387e5765293dcba39c0c5732"

	copy(buffer[:], msg)
	hw := New()
	for i := 0; i < 16; i++ {
		hw.Write(buffer[:])
	}

	uhash := make([]uint8, 32)
	hw.Sum(uhash[:0])
	calcVal := Byte2String(uhash)

	if calcVal != trueVal {
		t.Errorf("expected: %s,\nbut got: %s\n", trueVal, calcVal)
	}
}

func TestSM3_3(t *testing.T) {
	msg := "abcd"
	//var buffer [4]byte
	trueVal := "debe9ff92275b8a138604889c18e5a4d6fdb70e5387e5765293dcba39c0c5732"

	//copy(buffer[:], msg)
	buffer := []byte(msg)
	hw := New()
	for i := 0; i < 15; i++ {
		hw.Write(buffer[:])
	}

	uhash := make([]uint8, 32)
	hw.Sum(uhash[:0])

	// Continue write, the result still the same,
	// for hw.Sum() not change the hash state
	hw.Write(buffer[:])
	uhash = make([]uint8, 32)
	hw.Sum(uhash[:0])
	calcVal := Byte2String(uhash)

	if calcVal != trueVal {
		t.Errorf("Sum() of hash.Hash interface is not implemented.\nexpected: %s,\nbut got: %s\n", trueVal, calcVal)
	}
}

type sm3Test struct {
	out string
	in  string
}

var golden = []sm3Test{
	{"136CE3C86E4ED909B76082055A61586AF20B4DAB674732EBD4B599EEF080C9BE", "aaaaa"},
	{"623476AC18F65A2909E43C7FEC61B49C7E764A91A18CCB82F1917A29C86C5E88", "a"},
	{"E07D8EE6E54586A459E30EB8D809E02194558E2B0B235A31F3226A3687FAAB88", "ab"},
	{"66C7F0F462EEEDD9D1F2D46BDC10E4E24167C4875CF2F7A2297DA02B8F4BA8E0", "abc"},
	{"AFE4CCAC5AB7D52BCAE36373676215368BAF52D3905E1FECBE369CC120E97628", "abcde"},
}

func TestGolden(t *testing.T) {
	for i := 0; i < len(golden); i++ {
		g := golden[i]
		s := fmt.Sprintf("%x", Sum([]byte(g.in)))
		if s != strings.ToLower(g.out) {
			t.Fatalf("Sum function: sm3(%s) = %s want %s", g.in, s, g.out)
		}
		c := New()
		for j := 0; j < 3; j++ {
			if j < 2 {
				io.WriteString(c, g.in)
			} else {
				io.WriteString(c, g.in[0:len(g.in)/2])
				c.Sum(nil)
				io.WriteString(c, g.in[len(g.in)/2:])
			}
			s := fmt.Sprintf("%x", c.Sum(nil))
			if s != strings.ToLower(g.out) {
				t.Fatalf("sm3[%d](%s) = %s want %s", j, g.in, s, g.out)
			}
			c.Reset()
		}
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
