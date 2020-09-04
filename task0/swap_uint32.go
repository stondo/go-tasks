package main

import "fmt"

func main() {
	var m uint32 = 9
	var n uint32 = 19

	fmt.Println("m and n are:               ", m, n)

	n = m + n
	m = n - m
	n = n - m

	fmt.Println("m and n have been swapped: ", m, n)
}

// SwapUint32InPlace swaps 2 uint32 in place
func swapTwoUint32InPlace(m uint32, n uint32) (uint32, uint32) {
	fmt.Println("m and n are:               ", m, n)
	m, n = n, m
	fmt.Println("m and n have been swapped: ", m, n)
	return m, n
}
