package vip

import (
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
)

//Mask represents the size of a bit-mask
type Mask uint8

//BitMask returns the bit-mask
func (m Mask) BitMask() uint32 {
	return math.MaxUint32 << (32 - uint8(m))
}

//String number of bits
func (m Mask) String() string {
	return strconv.FormatUint(uint64(m), 10)
}

//Hex representation of bit mask
func (m Mask) Hex() string {
	return fmt.Sprintf("%x", m.BitMask())
}

//IP representation of bit mask
func (m Mask) IP() IP {
	return IP(uint32(MaxIP) & m.BitMask())
}

//ToStd converts back to standard-library mask
func (m Mask) ToStd() net.IPMask {
	a, b, c, d := quadSpread(m.BitMask())
	return net.IPMask([]byte{a, b, c, d})
}

//MaskBytes converts a quad into a mask
func MaskBytes(b [4]byte) (Mask, error) {
	target := quadJoin(b[0], b[1], b[2], b[3])
	for b := uint8(32); b >= 0; b-- {
		if target == math.MaxUint32<<(32-b) {
			return Mask(b), nil
		}
	}
	return Mask(0), errors.New("invalid mask")
}

//StdMask converts a standard-library mask
func StdMask(m net.IPMask) (Mask, error) {
	b := [4]byte{}
	n := copy(b[:], m)
	if n != 4 {
		return Mask(0), errors.New("invalid mask")
	}
	return MaskBytes(b)
}
