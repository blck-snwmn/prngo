package pcg

import (
	"testing"
)

func TestPCG(t *testing.T) {
	seed := uint64(1111111111111)
	p := New(seed)
	wants := []uint32{
		1207103580,
		865683596,
		2250586291,
		4181475314,
		3072241397,
	}
	for _, want := range wants {
		got := p.next()
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	}
}
