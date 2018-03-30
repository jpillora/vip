package vip

import (
	"encoding/json"
	"strconv"
	"strings"
	"testing"
)

func TestNW1(t *testing.T) {
	m := Mask(24)
	s := strconv.FormatUint(uint64(m.BitMask()), 2)
	if s != "11111111111111111111111100000000" {
		t.Fatalf("bad str: %s", s)
	}
}

func TestNW2(t *testing.T) {
	nw := QuadMask(10, 0, 0, 0, 24)
	s := nw.String()
	if s != "10.0.0.0/24" {
		t.Fatalf("bad str: %s", s)
	}
}

func TestNW3(t *testing.T) {
	nw := QuadMask(10, 0, 0, 7, 24)
	if nw.IP != Quad(10, 0, 0, 7) {
		t.Fatalf("bad nw ip")
	}
	if nw.NetworkIP() != Quad(10, 0, 0, 0) {
		t.Fatalf("bad nw addr")
	}
	if ip := nw.BroadcastIP(); ip != Quad(10, 0, 0, 255) {
		t.Fatalf("bad bc addr: %s", ip)
	}
}
func TestNWParse1(t *testing.T) {
	nw, err := ParseCIDR("1.2.3.4/24")
	if err != nil {
		t.Fatalf("parse err: %s", err)
	} else if nw != QuadMask(1, 2, 3, 4, 24) {
		t.Fatalf("parse mismatch")
	}
}
func TestNWParse2(t *testing.T) {
	_, err := ParseCIDR("1.2.3.4/33")
	if err == nil {
		t.Fatalf("expected err")
	} else if !strings.Contains(err.Error(), "value out of range") {
		t.Fatalf("expected range err")
	}
}
func TestNWParse3(t *testing.T) {
	nw, err := ParseCIDR("5.196.192.216/32")
	if err != nil {
		t.Fatalf("parse err: %s", err)
	} else if nw != QuadMask(5, 196, 192, 216, 32) {
		t.Fatalf("parse mismatch")
	}
}
func TestNWContains1(t *testing.T) {
	nw := MustParseCIDR("10.1.1.0/24")
	if !nw.Contains(MustParse("10.1.1.13")) {
		t.Fatalf("expected to be in network")
	} else if nw.Contains(MustParse("10.1.0.13")) {
		t.Fatalf("expected not to be in network 1")
	} else if nw.Contains(MustParse("10.0.1.13")) {
		t.Fatalf("expected not to be in network 2")
	}
}

func TestNWContains2(t *testing.T) {
	nw := MustParseCIDR("10.1.1.17/28")
	if nw.NetworkIP() != Quad(10, 1, 1, 16) {
		t.Fatalf("unexpected network addr")
	}
	if nw.BroadcastIP() != Quad(10, 1, 1, 31) {
		t.Fatalf("unexpected broadcast addr")
	}
	if nw.Contains(MustParse("10.1.1.13")) {
		t.Fatalf("expected not to be in network 1")
	}
	if !nw.Contains(MustParse("10.1.1.19")) {
		t.Fatalf("expected to be in network 2")
	}
	if nw.Contains(MustParse("10.1.1.32")) {
		t.Fatalf("expected not to be in network 3")
	}
}

func TestNWSize(t *testing.T) {
	if MustParseCIDR("10.1.1.0/24").Size() != 256 {
		t.Fatalf("expected size (1)")
	} else if MustParseCIDR("10.1.1.0/28").Size() != 16 {
		t.Fatalf("expected size (2)")
	}
}

func TestNWJSON(t *testing.T) {
	data := struct {
		Net IPNet `json:"net"`
	}{}
	buff := []byte(`{"net":"10.0.0.7/24"}`)
	if err := json.Unmarshal(buff, &data); err != nil {
		t.Fatalf("unmarshal json: %s", err)
	} else if data.Net != QuadMask(10, 0, 0, 7, 24) {
		t.Fatalf("data mismatch")
	}
	//change address
	data.Net.IP = Quad(10, 3, 3, 6)
	//reencode
	if buff2, err := json.Marshal(&data); err != nil {
		t.Fatalf("marshal json: %s", err)
	} else if string(buff2) != `{"net":"10.3.3.6/24"}` {
		t.Fatalf("buff mismatch")
	}
}
