package main

import "fmt"

// sumArray рекурсивно обчислює суму всіх чисел у вкладених масивах
func sumArray(arr []interface{}) int {
	sum := 0

	for _, item := range arr {
		switch v := item.(type) {
		case int:
			// Якщо елемент - число, додаємо до суми
			sum += v
		case []interface{}:
			// Якщо елемент - масив, рекурсивно обчислюємо його суму
			sum += sumArray(v)
		default:
			// Ігноруємо інші типи (strings, nil, etc.)
			continue
		}
	}

	return sum
}

func main() {
	// Test 1: Simple array
	arr1 := []interface{}{1, 2, 3, 4, 5}
	fmt.Printf("sumArray(%v) = %d\n", arr1, sumArray(arr1))
	// Output: 15

	// Test 2: Nested array
	arr2 := []interface{}{1, []interface{}{2, 3}, 4, []interface{}{5, 6}}
	fmt.Printf("sumArray(%v) = %d\n", arr2, sumArray(arr2))
	// Output: 21

	// Test 3: Empty array
	arr3 := []interface{}{}
	fmt.Printf("sumArray(%v) = %d\n", arr3, sumArray(arr3))
	// Output: 0

	// Test 4: Single element
	arr4 := []interface{}{10}
	fmt.Printf("sumArray(%v) = %d\n", arr4, sumArray(arr4))
	// Output: 10

	// Test 5: Deep nesting
	arr5 := []interface{}{
		1,
		[]interface{}{
			2,
			[]interface{}{
				3,
				[]interface{}{
					4,
					[]interface{}{5},
				},
			},
		},
	}
	fmt.Printf("sumArray(deeply nested) = %d\n", sumArray(arr5))
	// Output: 15

	// Test 6: With non-numbers (ignored)
	arr6 := []interface{}{1, "a", 2, nil, 3}
	fmt.Printf("sumArray(%v) = %d\n", arr6, sumArray(arr6))
	// Output: 6
}
