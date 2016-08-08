package sm3

import (
	// "bytes"
	"encoding/binary"
	"fmt"
	// "hash"
)

type SM3 struct {
	digest   [8]uint32  // digest represents the partial evaluation of V
	T        [64]uint32 // constant
	message  []byte     // uint8  //
	hashcode string     // hash result
}

func NewSM3(message []byte) *SM3 {
	sm3 := &SM3{message: message, hashcode: ""}
	sm3.digest[0] = 0x7380166f
	sm3.digest[1] = 0x4914b2b9
	sm3.digest[2] = 0x172442d7
	sm3.digest[3] = 0xda8a0600
	sm3.digest[4] = 0xa96f30bc
	sm3.digest[5] = 0x163138aa
	sm3.digest[6] = 0xe38dee4d
	sm3.digest[7] = 0xb0fb0e4e

	// Set T[i]
	for i := 0; i < 16; i++ {
		sm3.T[i] = 0x79cc4519
	}
	for i := 16; i < 64; i++ {
		sm3.T[i] = 0x7a879d8a
	}

	// Padding
	sm3.messagePad()

	// reset hashcode
	sm3.hashcode = ""

	return sm3
}

func (sm3 *SM3) printMsg() {
	for i := 0; i < len(sm3.message); i++ {
		fmt.Printf("%02x", sm3.message[i])
		if i%4 == 3 {
			fmt.Printf(" ")
		}
		if i%(4*8) == 31 {
			fmt.Println("")
		}
	}
	fmt.Println("")
}

func (sm3 *SM3) messagePad() {
	fmt.Println("messagePad------->")
	fmt.Println("Before padding:")
	sm3.printMsg()

	msgLen := uint64(len(sm3.message) * 8)

	// append '1'
	sm3.message = append(sm3.message, 0x80)

	// append until the resulting message length (in bits) is congruent to 448 (mod 512)
	blockSize := 64
	for len(sm3.message)%blockSize != 56 {
		sm3.message = append(sm3.message, 0x00)
	}

	// append message length
	sm3.message = append(sm3.message, uint8(msgLen>>56&0xff))
	sm3.message = append(sm3.message, uint8(msgLen>>48&0xff))
	sm3.message = append(sm3.message, uint8(msgLen>>40&0xff))
	sm3.message = append(sm3.message, uint8(msgLen>>32&0xff))
	sm3.message = append(sm3.message, uint8(msgLen>>24&0xff))
	sm3.message = append(sm3.message, uint8(msgLen>>16&0xff))
	sm3.message = append(sm3.message, uint8(msgLen>>8&0xff))
	sm3.message = append(sm3.message, uint8(msgLen>>0&0xff))

	if len(sm3.message)%64 != 0 {
		fmt.Println("------messagePad: Error, length =", len(sm3.message))
	}

	fmt.Println("After padding:")
	sm3.printMsg()
}

func (sm3 *SM3) ff0(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func (sm3 *SM3) ff1(x, y, z uint32) uint32 {
	return (x & y) | (x & z) | (y & z)
}

func (sm3 *SM3) gg0(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func (sm3 *SM3) gg1(x, y, z uint32) uint32 {
	return (x & y) | (^x & z)
}

func (sm3 *SM3) p0(x uint32) uint32 {
	return x ^ sm3.leftRotate(x, 9) ^ sm3.leftRotate(x, 17)
}

func (sm3 *SM3) p1(x uint32) uint32 {
	return x ^ sm3.leftRotate(x, 15) ^ sm3.leftRotate(x, 23)
}

func (sm3 *SM3) messageExtend(data []byte) (W [68]uint32, W1 [64]uint32) {
	fmt.Println("messageExtend--------->")

	// big endian
	for i := 0; i < 16; i++ {
		W[i] = binary.BigEndian.Uint32(data[4*i : 4*(i+1)])
	}
	for i := 16; i < 68; i++ {
		W[i] = sm3.p1(W[i-16]^W[i-9]^sm3.leftRotate(W[i-3], 15)) ^ sm3.leftRotate(W[i-13], 7) ^ W[i-6]
	}
	for i := 0; i < 64; i++ {
		W1[i] = W[i] ^ W[i+4]
	}
	return W, W1
}

func (sm3 *SM3) leftRotate(x uint32, i uint32) uint32 {
	i %= 32
	return (x<<i | x>>(32-i))
}

// cf is compress function
func (sm3 *SM3) cf(W [68]uint32, W1 [64]uint32) {
	fmt.Println("cf------->")

	A := sm3.digest[0]
	B := sm3.digest[1]
	C := sm3.digest[2]
	D := sm3.digest[3]
	E := sm3.digest[4]
	F := sm3.digest[5]
	G := sm3.digest[6]
	H := sm3.digest[7]

	for i := 0; i < 16; i++ {
		SS1 := sm3.leftRotate(sm3.leftRotate(A, 12)+E+sm3.leftRotate(sm3.T[i], uint32(i)), 7)
		SS2 := SS1 ^ sm3.leftRotate(A, 12)
		TT1 := sm3.ff0(A, B, C) + D + SS2 + W1[i]
		TT2 := sm3.gg0(E, F, G) + H + SS1 + W[i]
		D = C
		C = sm3.leftRotate(B, 9)
		B = A
		A = TT1
		H = G
		G = sm3.leftRotate(F, 19)
		F = E
		E = sm3.p0(TT2)

		// debug
		fmt.Printf("%02d: ", i)
		fmt.Printf("%08x ", A)
		fmt.Printf("%08x ", B)
		fmt.Printf("%08x ", C)
		fmt.Printf("%08x ", D)
		fmt.Printf("%08x ", E)
		fmt.Printf("%08x ", F)
		fmt.Printf("%08x ", G)
		fmt.Printf("%08x\n", H)
	}

	for i := 16; i < 64; i++ {
		SS1 := sm3.leftRotate(sm3.leftRotate(A, 12)+E+sm3.leftRotate(sm3.T[i], uint32(i)), 7)
		SS2 := SS1 ^ sm3.leftRotate(A, 12)
		TT1 := sm3.ff1(A, B, C) + D + SS2 + W1[i]
		TT2 := sm3.gg1(E, F, G) + H + SS1 + W[i]
		D = C
		C = sm3.leftRotate(B, 9)
		B = A
		A = TT1
		H = G
		G = sm3.leftRotate(F, 19)
		F = E
		E = sm3.p0(TT2)

		// debug
		fmt.Printf("%02d: ", i)
		fmt.Printf("%08x ", A)
		fmt.Printf("%08x ", B)
		fmt.Printf("%08x ", C)
		fmt.Printf("%08x ", D)
		fmt.Printf("%08x ", E)
		fmt.Printf("%08x ", F)
		fmt.Printf("%08x ", G)
		fmt.Printf("%08x\n", H)
	}

	sm3.digest[0] ^= A
	sm3.digest[1] ^= B
	sm3.digest[2] ^= C
	sm3.digest[3] ^= D
	sm3.digest[4] ^= E
	sm3.digest[5] ^= F
	sm3.digest[6] ^= G
	sm3.digest[7] ^= H
}

func (sm3 *SM3) Hash() string {
	blockSize := 64
	nblocks := len(sm3.message) / blockSize

	for i := 0; i < nblocks; i++ {
		startPos := i * blockSize
		W, W1 := sm3.messageExtend(sm3.message[startPos : startPos+blockSize])

		// debug
		printUint32Slice(W[:])
		printUint32Slice(W1[:])

		sm3.cf(W, W1)
	}

	// Transform sm3.digest[] to hashcode string
	sm3.constructHashCode()

	sm3.printValues()

	return sm3.hashcode
}

func (sm3 *SM3) printValues() {
	for i := 0; i < 8; i++ {
		fmt.Printf("%x ", sm3.digest[i])
	}
	fmt.Printf("\n")
	// fmt.Println(sm3.hashcode)
}

func (sm3 *SM3) constructHashCode() {
	fmt.Println("constructHashCode------------------------->")
	for i := 0; i < 8; i++ {
		sm3.hashcode += fmt.Sprintf("%x", sm3.digest[i])
	}
	fmt.Println(sm3.hashcode)
}

func (sm3 *SM3) GetHashCode() string {
	return sm3.hashcode
}

func printUint32Slice(list []uint32) {
	for i := 0; i < len(list); i++ {
		fmt.Printf("%08x ", list[i])
		if i%8 == 7 {
			fmt.Println("")
		}
	}
	fmt.Println("")
}
