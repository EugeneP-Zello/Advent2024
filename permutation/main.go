package main

import (
	"fmt"
	"sort"
)

func comparePerm(p1 []int32, p2 []int32) bool {
	for idx := 0; idx < len(p1); idx++ {
		if p1[idx] == p2[idx] {
			continue
		}
		return p1[idx] < p2[idx]
	}
	return true
}

func CheckPermutation(values []int32, best []int32) {
	v0 := 0
	for idx := 0; idx < len(values); idx++ {
		v0 += int(values[idx] * best[idx])
	}
	// check all permutations
	bestAll := findBestPemutationO2(values)

	for idx := 0; idx < len(values); idx++ {
		if best[idx] != bestAll[idx] {
			fmt.Printf("Error: %v != %v\n", best[idx], bestAll[idx])
			return
		}
	}
	fmt.Printf("OK: %d %v\n", v0, best)
}

func findBestPemutation(values []int32) []int32 {
	maxIdx := int32(len(values))
	best := make([]int32, maxIdx)
	for i := int32(0); i < maxIdx; i++ {
		best[i] = i
	}
	sort.Slice(best, func(i, j int) bool {
		return values[best[i]] < values[best[j]]
	})

	sorted := make([]int32, len(values))
	for i, idx := range best {
		sorted[i] = values[idx]
	}

	indexes := make([]int32, maxIdx)
	for i, v := range best {
		indexes[int(v)] = int32(i)
	}

	total := 0
	for idx := 0; idx < len(values); idx++ {
		indexes[idx] = indexes[idx] + 1
		total += int(values[idx] * indexes[idx])
	}
	fmt.Printf("BestPemutation(%d): %v\n", total, indexes)

	return indexes
}

func findBestPemutationO2(values []int32) []int32 {
	permutations := generatePermutations(int32(len(values)))
	best := make([]int32, len(values))
	bestValue := 0
	for _, perm := range permutations {
		v := 0
		for idx := 0; idx < len(values); idx++ {
			v += int(values[idx] * perm[idx])
		}
		if v > bestValue {
			bestValue = v
			copy(best, perm)
		} else if v == bestValue {
			if comparePerm(perm, best) {
				copy(best, perm)
			}
		}
	}
	return best
}

func permute(arr []int32, l int, r int, result *[][]int32) {
	if l == r {
		perm := make([]int32, len(arr))
		copy(perm, arr)
		*result = append(*result, perm)
	} else {
		for i := l; i <= r; i++ {
			arr[l], arr[i] = arr[i], arr[l]
			permute(arr, l+1, r, result)
			arr[l], arr[i] = arr[i], arr[l] // backtrack
		}
	}
}

func generatePermutations(n int32) [][]int32 {
	arr := make([]int32, n)
	for i := int32(0); i < n; i++ {
		arr[i] = i + 1
	}
	var result [][]int32
	permute(arr, 0, int(n-1), &result)
	return result
}

func main() {
	initValues := [6]int32{2, 4, 1, 2, 3, 2}

	//initValues := [6]int32{2, 4, 1, 8, 3, 5}
	//initValues := [4]int32{20, 40, 10, 5}

	res := findBestPemutation(initValues[:])

	CheckPermutation(initValues[:], res)

}