package suggest

import (
	"strings"
)

const NoSuggestion = ""

func FindClosestMatch(input string, pool []string) string {
	minDistance := len(input)
	closestMatch := NoSuggestion

	for _, p := range pool {
		distance := levenshteinDistance(input, p)
		if distance < minDistance {
			minDistance = distance
			closestMatch = p
		}
	}
	return closestMatch
}

func levenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	m := len(s1)
	n := len(s2)
	d := make([][]int, m+1)
	for i := range d {
		d[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		d[i][0] = i
	}
	for j := 0; j <= n; j++ {
		d[0][j] = j
	}

	for j := 1; j <= n; j++ {
		for i := 1; i <= m; i++ {
			if s1[i-1] == s2[j-1] {
				d[i][j] = d[i-1][j-1]
			} else {
				min := d[i-1][j]
				if d[i][j-1] < min {
					min = d[i][j-1]
				}
				if d[i-1][j-1] < min {
					min = d[i-1][j-1]
				}
				d[i][j] = min + 1
			}
		}
	}

	return d[m][n]
}
