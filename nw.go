package vip

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

//IPNet reprents an IPv4 address and network mask
type IPNet struct {
	IP
	Mask
}

//CIDR is an alias for IPNet
type CIDR = IPNet

//IPNets is a slice of IPNet
type IPNets []IPNet

//EmptyIPNet is 0.0.0.0/0
var EmptyIPNet = IPNet{}

//ParseHostCIDR converts a CIDR string into an IPNet,
//and ensures the resulting IP is a host.
func ParseHostCIDR(cidr string) (IPNet, error) {
	nw, err := ParseCIDR(cidr)
	if err != nil {
		return IPNet{}, err
	}
	if !nw.IsHost(nw.IP) {
		return IPNet{}, errors.New("expecting host ip")
	}
	return nw, nil
}

//ParseCIDR converts a CIDR string into an IPNet
func ParseCIDR(cidr string) (IPNet, error) {
	c := strings.SplitN(cidr, "/", 2)
	if len(c) != 2 {
		return IPNet{}, errors.New("expected <ip>/<mask>")
	}
	ip, err := Parse(c[0])
	if err != nil {
		return IPNet{}, err
	}
	mask, err := strconv.ParseUint(c[1], 10, 8)
	if err != nil {
		return IPNet{}, err
	} else if mask > 32 {
		return IPNet{}, errors.New("value out of range")
	}
	return IPNet{IP: ip, Mask: Mask(mask)}, nil
}

//MustParseCIDR panics on failure
func MustParseCIDR(cidr string) IPNet {
	x, err := ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return x
}

//QuadMask creates a IPNet from a quad and a mask
func QuadMask(a uint8, b uint8, c uint8, d uint8, bits uint8) IPNet {
	return IPNet{IP: Quad(a, b, c, d), Mask: Mask(bits)}
}

func (nw IPNet) String() string {
	return nw.IP.String() + "/" + nw.Mask.String()
}

//NetworkIP returns the network address
//For example 10.0.0.7/24 returns 10.0.0.0
func (nw IPNet) NetworkIP() IP {
	return IP(uint32(nw.IP) & nw.BitMask())
}

//IsNetworkIP determines if the IP is network's
//network address
func (nw IPNet) IsNetworkIP(ip IP) bool {
	return nw.Mask != 32 && nw.NetworkIP() == ip
}

//BroadcastIP returns the broadcast address
//For example 10.0.0.7/24 returns 10.0.0.255
func (nw IPNet) BroadcastIP() IP {
	return IP(uint32(nw.IP) | ^nw.BitMask())
}

//IsBroadcastIP determines if the IP is network's
//broadcast address
func (nw IPNet) IsBroadcastIP(ip IP) bool {
	return nw.Mask != 32 && nw.BroadcastIP() == ip
}

//IsHost determines if the IP is a host in the
//network (not the network ad not broadcast address)
func (nw IPNet) IsHost(ip IP) bool {
	return nw.Mask == 32 || (!nw.IsNetworkIP(ip) && !nw.IsBroadcastIP(ip))
}

//Contains determines if IP is contained in the network
func (nw IPNet) Contains(ip IP) bool {
	return nw.NetworkIP() == IP(uint32(ip)&nw.BitMask())
}

//Size represents number of addresses in the network
func (nw IPNet) Size() uint32 {
	return 1 << (32 - nw.Mask)
}

//MarshalJSON allows IPNets to be json.Unmarshalled
func (nw IPNet) MarshalJSON() ([]byte, error) {
	return []byte(`"` + nw.String() + `"`), nil
}

//UnmarshalJSON allows IPNets to be json.Unmarshalled
func (nw *IPNet) UnmarshalJSON(b []byte) error {
	tmp, err := ParseCIDR(strings.Trim(string(b), `"`))
	if err != nil {
		return err
	}
	*nw = tmp
	return nil
}

//MarshalText ...
func (nw IPNet) MarshalText() (text []byte, err error) {
	return []byte(nw.String()), nil
}

//UnmarshalText ...
func (nw *IPNet) UnmarshalText(text []byte) error {
	return nw.UnmarshalJSON(text)
}

//StdNet converts a standard-library *net.IPNet to a vip.IPNet
func StdNet(nw *net.IPNet) IPNet {
	m, err := StdMask(nw.Mask)
	if err != nil {
		m = Mask(32)
	}
	return IPNet{
		IP:   StdIP(nw.IP),
		Mask: m,
	}
}

//ToStd ...
func (nw IPNet) ToStd() net.IPNet {
	return net.IPNet{
		IP:   nw.NetworkIP().ToStd(),
		Mask: nw.Mask.ToStd(),
	}
}
