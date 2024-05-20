package pcg

const (
	multiplier = 6364136223846793005
	increment  = 1442695040888963407
)

type PCG struct {
	state uint64
}

func New(seed uint64) *PCG {
	pcg := &PCG{state: seed + increment}
	pcg.next()
	return pcg
}

func (p *PCG) next() uint32 {
	state := p.state
	p.state = update(state)
	x := xsh(state)
	return rr(uint32(x), uint(state>>59))
}

func rr(n uint32, shift uint) uint32 {
	return n>>shift | n<<(32-shift)
}

func xsh(state uint64) uint64 {
	shiftedState := state >> 18
	return (shiftedState ^ state) >> 27
}

func update(state uint64) uint64 {
	return state*multiplier + increment
}
