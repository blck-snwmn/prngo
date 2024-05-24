package chacha8

import (
	"encoding/binary"
	"fmt"
	"math/rand/v2"
)

type keySizeError int

func (k keySizeError) Error() string {
	return fmt.Sprintf("invalid key length. got=%d, want=%d", k, 32)
}

type nonceSizeError int

func (n nonceSizeError) Error() string {
	return fmt.Sprintf("invalid nonce length. got=%d, want=%d", n, 12)
}
func newState(key, nonce []byte, counter uint32) (state, error) {
	if len(key) != 32 {
		return nil, keySizeError(len(key))
	}
	if len(nonce) != 12 {
		return nil, nonceSizeError(len(nonce))
	}
	s := make(state, 16)

	// magic
	s[0] = 0x61707865
	s[1] = 0x3320646e
	s[2] = 0x79622d32
	s[3] = 0x6b206574

	s[4] = binary.LittleEndian.Uint32(key[0:4])
	s[5] = binary.LittleEndian.Uint32(key[4:8])
	s[6] = binary.LittleEndian.Uint32(key[8:12])
	s[7] = binary.LittleEndian.Uint32(key[12:16])

	s[8] = binary.LittleEndian.Uint32(key[16:20])
	s[9] = binary.LittleEndian.Uint32(key[20:24])
	s[10] = binary.LittleEndian.Uint32(key[24:28])
	s[11] = binary.LittleEndian.Uint32(key[28:32])

	s[12] = counter
	s[13] = binary.LittleEndian.Uint32(nonce[0:4])
	s[14] = binary.LittleEndian.Uint32(nonce[4:8])
	s[15] = binary.LittleEndian.Uint32(nonce[8:12])

	return s, nil
}

type state []uint32

func (s state) quarterRound(x, y, z, w uint) {
	s[x], s[y], s[z], s[w] = quarterRound(s[x], s[y], s[z], s[w])
}

func rotationN(n uint32, shift uint) uint32 {
	return n>>(32-shift) | n<<shift
}

func quarterRound(a, b, c, d uint32) (uint32, uint32, uint32, uint32) {
	a += b
	d ^= a
	d = rotationN(d, 16)

	c += d
	b ^= c
	b = rotationN(b, 12)

	a += b
	d ^= a
	d = rotationN(d, 8)

	c += d
	b ^= c
	b = rotationN(b, 7)
	return a, b, c, d
}

func (s state) innerBlock() {
	// column rounds
	s.quarterRound(0, 4, 8, 12)
	s.quarterRound(1, 5, 9, 13)
	s.quarterRound(2, 6, 10, 14)
	s.quarterRound(3, 7, 11, 15)

	// diagonal rounds
	s.quarterRound(0, 5, 10, 15)
	s.quarterRound(1, 6, 11, 12)
	s.quarterRound(2, 7, 8, 13)
	s.quarterRound(3, 4, 9, 14)
}

func (s state) clone() state {
	newS := make(state, 16)
	copy(newS, s)
	return newS
}

type Chacha8 struct {
	states []state
	flip   bool
	i      uint32
}

var _ rand.Source = (*Chacha8)(nil)

func NewChaCha8(seed [32]byte) *Chacha8 {
	ss := make([]state, 0, 4)
	for i := 0; i < 4; i++ {
		nonce := make([]byte, 12)
		s, _ := newState(seed[:], nonce, uint32(i))
		init := s.clone()
		for i := 0; i < 4; i++ { // 4 iterations = 8 rounds
			s.innerBlock()
		}
		s[4] += init[4]
		s[5] += init[5]
		s[6] += init[6]
		s[7] += init[7]
		s[8] += init[8]
		s[9] += init[9]
		s[10] += init[10]
		s[11] += init[11]

		ss = append(ss, s)
	}
	return &Chacha8{
		states: ss,
		flip:   true,
	}
}

func (c *Chacha8) Uint64() uint64 {
	data := make([]byte, 8)

	index := c.i

	l, r := 0, 1
	if !c.flip {
		l, r = 2, 3
		// if flip false, increment i
		c.i++
	}
	c.flip = !c.flip
	// TODO bigendian?
	binary.LittleEndian.PutUint32(data[:], c.states[l][index])
	binary.LittleEndian.PutUint32(data[4:], c.states[r][index])
	return binary.LittleEndian.Uint64(data[:])
}
