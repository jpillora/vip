package vip

import "testing"

func TestMask(t *testing.T) {
	m, err := MaskBytes([4]byte{255, 255, 255, 0})
	if err != nil {
		t.Fatal(err)
	}
	if uint(m) != 24 {
		t.Fatalf("unexpected mask: %d", m)
	}
}
