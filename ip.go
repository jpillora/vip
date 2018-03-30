package vip

import (
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

//IP represents a single IPv4 address
type IP uint32

//EmptyIP is 0.0.0.0
var EmptyIP = IP(0)

//MaxIP is 255.255.255.255
var MaxIP = Quad(255, 255, 255, 255)

//Parse an IPv4 "dotted-quad"
func Parse(ip string) (IP, error) {
	q := strings.SplitN(ip, ".", 4)
	if len(q) != 4 {
		return IP(0), errors.New("expected dotted-quad")
	}
	u := uint32(0)
	for i := 0; i < 4; i++ {
		n, err := strconv.ParseUint(q[i], 10, 8)
		if err != nil {
			return IP(0), err
		}
		u = u | uint32(n)<<uint8((3-i)*8)
	}
	return IP(u), nil
}

//MustParse panics on Parse failure
func MustParse(ip string) IP {
	x, err := Parse(ip)
	if err != nil {
		panic(err)
	}
	return x
}

func (ip IP) String() string {
	a, b, c, d := ip.Quad()
	return fmt.Sprintf("%d.%d.%d.%d", a, b, c, d)
}

//IsZero returns if the IP is 0.0.0.0
func (ip IP) IsZero() bool {
	return ip == EmptyIP
}

//Delta many spots from the IP
func (ip IP) Delta(d int64) IP {
	n := int64(ip) + d
	if n < 0 {
		return EmptyIP
	} else if n > math.MaxUint32 {
		return MaxIP
	}
	return IP(n)
}

//Next IP after this one
func (ip IP) Next() IP {
	return ip.Delta(1)
}

//Prev IP before this one
func (ip IP) Prev() IP {
	return ip.Delta(-1)
}

//Quad returns 4 bytes representing each quad
func (ip IP) Quad() (uint8, uint8, uint8, uint8) {
	return quadSpread(uint32(ip))
}

//Quad returns 4 bytes representing each quad
func Quad(a uint8, b uint8, c uint8, d uint8) IP {
	return IP(quadJoin(a, b, c, d))
}

//MarshalJSON allows IPNets to be json.Unmarshalled
func (ip IP) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ip.String() + `"`), nil
}

//UnmarshalJSON allows IPNets to be json.Unmarshalled
func (ip *IP) UnmarshalJSON(b []byte) error {
	tmp, err := Parse(strings.Trim(string(b), `"`))
	if err != nil {
		return err
	}
	*ip = tmp
	return nil
}

//StdIP converts a standard-library net.IP to a vip.IP
func StdIP(ip net.IP) IP {
	v4 := ip.To4()
	return Quad(v4[0], v4[1], v4[2], v4[3])
}

func (ip IP) ToStd() net.IP {
	a, b, c, d := ip.Quad()
	return net.IP([]byte{a, b, c, d})
}

func (ip IP) IsMulticast() bool {
	return Multicast.Contains(ip)
}

func (ip IP) IsSSDP() bool {
	return ip == SSDP
}

func (ip IP) Mask(bits uint8) IPNet {
	return IPNet{
		IP:   ip,
		Mask: Mask(bits),
	}
}

//IPBytes 4 bytes to an vip.IP
func IPBytes(b [4]byte) IP {
	return Quad(b[0], b[1], b[2], b[3])
}
