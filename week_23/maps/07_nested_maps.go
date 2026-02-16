package main

import "fmt"

// 07. Nested Maps - map[string]map[string]int and similar

func main() {
	// Scoreboard: player -> game -> score
	scores := map[string]map[string]int{
		"Alice": {"game1": 100, "game2": 200},
		"Bob":   {"game1": 150, "game2": 180},
	}

	fmt.Println("Scores:", scores)
	fmt.Println("Alice game1:", scores["Alice"]["game1"])

	// Add new nested map - must initialize
	if scores["Carol"] == nil {
		scores["Carol"] = make(map[string]int)
	}
	scores["Carol"]["game1"] = 120
	fmt.Println("After adding Carol:", scores)

	// Nested lookup with existence check
	if player, ok := scores["Alice"]; ok {
		if score, ok := player["game2"]; ok {
			fmt.Println("Alice game2:", score)
		}
	}
}
