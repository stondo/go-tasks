package main

import (
	"testing"
)

func TestSwapUint32InPlace(t *testing.T) {
	m, n := swapTwoUint32InPlace(9, 19)

	if m != 19 && n != 9 {
		t.Errorf("swapTwoUint32InPlace(9, 19) = %d, %d; want 19, 9", m, n)
	}
}
