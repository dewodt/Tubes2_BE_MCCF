package models

// Path Data Structure (simple array of numbers)
// Each number represents the index from the articles []Article array.
// example: const path: Path = [0, 5, 4, 3] represents articles[0] (start)->articles[5]->articles[4]->articles[3] (target)
type Path = []int
