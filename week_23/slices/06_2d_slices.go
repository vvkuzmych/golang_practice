package main

import "fmt"

// 06. 2D Slices - Multi-dimensional slices (slice of slices)

func main() {
	// Rectangular 2D slice
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	fmt.Println("Matrix:", matrix)
	fmt.Println("matrix[1][2]:", matrix[1][2])

	// Jagged 2D slice - rows can have different lengths
	jagged := [][]int{
		{1},
		{2, 3},
		{4, 5, 6},
	}
	fmt.Println("Jagged:", jagged)

	// Create 3x4 matrix with make
	rows, cols := 3, 4
	grid := make([][]int, rows)
	for i := range grid {
		grid[i] = make([]int, cols)
	}
	grid[0][0] = 1
	fmt.Println("3x4 grid:", grid)
}
