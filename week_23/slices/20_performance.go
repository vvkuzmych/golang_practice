package main

import "fmt"

// 20. Performance Tips - Preallocate, reuse, avoid unnecessary allocations

func main() {
	// 1. Preallocate capacity when size is known
	knownSize := 100
	slice := make([]int, 0, knownSize)
	for i := 0; i < knownSize; i++ {
		slice = append(slice, i)
	}
	fmt.Printf("Preallocated: len=%d cap=%d\n", len(slice), cap(slice))

	// 2. Prefer copy over append for existing slices
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, len(src))
	copy(dst, src)
	fmt.Println("Copied:", dst)

	// 3. Reuse slice buffer - reset with slice = slice[:0]
	buf := make([]int, 0, 100)
	buf = append(buf, 1, 2, 3)
	fmt.Println("Before reset:", buf)
	buf = buf[:0] // reuse same underlying array
	buf = append(buf, 4, 5, 6)
	fmt.Println("After reset:", buf)

	// 4. Avoid allocation in loop - use index
	for i := range slice {
		slice[i] = slice[i] * 2 // no allocation
	}
	fmt.Println("In-place double:", slice[:5])

	// 5. Use []byte for string building when possible
	bytes := make([]byte, 0, 32)
	bytes = append(bytes, "hello"...)
	bytes = append(bytes, " world"...)
	fmt.Println("String from bytes:", string(bytes))
}
