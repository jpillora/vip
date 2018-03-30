package vip

import (
	"strings"
	"testing"
)

func TestIP1(t *testing.T) {
	ip := Quad(1, 1, 1, 1)
	s := ip.String()
	if s != "1.1.1.1" {
		t.Fatalf("bad str: %s", s)
	}
}

func TestIP2(t *testing.T) {
	ip := Quad(1, 2, 3, 4)
	s := ip.String()
	if s != "1.2.3.4" {
		t.Fatalf("bad str: %s", s)
	}
}

func TestIPParse1(t *testing.T) {
	ip, err := Parse("1.2.3.4")
	if err != nil {
		t.Fatalf("parse err: %s", err)
	} else if ip != Quad(1, 2, 3, 4) {
		t.Fatalf("parse mismatch")
	}
}

func TestIPParse2(t *testing.T) {
	_, err := Parse("1.2.333.4")
	if err == nil {
		t.Fatalf("expected err")
	} else if !strings.Contains(err.Error(), "value out of range") {
		t.Fatalf("expected range err")
	}
}
