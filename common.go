package vip

//commmon ranges and such

var (
	Multicast = QuadMask(224, 0, 0, 0, 4)
	SSDP      = Quad(239, 255, 255, 250)
)

func quadSpread(n uint32) (uint8, uint8, uint8, uint8) {
	a := uint8(n >> 24)
	b := uint8(n >> 16)
	c := uint8(n >> 8)
	d := uint8(n >> 0)
	return a, b, c, d
}

func quadJoin(a uint8, b uint8, c uint8, d uint8) uint32 {
	u := uint32(0)
	u = u | (uint32(a) << 24)
	u = u | (uint32(b) << 16)
	u = u | (uint32(c) << 8)
	u = u | (uint32(d) << 0)
	return u
}
