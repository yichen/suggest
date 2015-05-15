package suggest

import (
	"sort"
)

type RankedResult struct {
	query  string
	result []string
}

func (r RankedResult) Len() int {
	return len(r.result)
}

func (r RankedResult) Swap(i, j int) {
	r.result[i], r.result[j] = r.result[j], r.result[i]
}

func (r RankedResult) Less(i, j int) bool {
	s1 := Score(r.query, r.result[i])
	s2 := Score(r.query, r.result[j])

	return s1 > s2
}

func Score(query, candidate string) float64 {
	lev := LevenshteinDistance(query, candidate)
	length := Max(len(candidate), len(query))
	return float64(length-lev) / float64(length+lev) //Jacard score
}

func SortByRank(query string, result []string) {
	// create a map from the score to result
	// then sort the score, and then map to the result

	rr := RankedResult{query, result}
	sort.Sort(rr)

	result = rr.result
}

//Levenshtein distance is the number of inserts, deletions,
//and substitutions that differentiate one word from another.
//This algorithm is dynamic programming found at
//http://en.wikipedia.org/wiki/Levenshtein_distance
func LevenshteinDistance(s, t string) int {
	m := len(s)
	n := len(t)
	width := n - 1
	d := make([]int, m*n)
	//y * w + h for position in array
	for i := 1; i < m; i++ {
		d[i*width+0] = i
	}

	for j := 1; j < n; j++ {
		d[0*width+j] = j
	}

	for j := 1; j < n; j++ {
		for i := 1; i < m; i++ {
			if s[i] == t[j] {
				d[i*width+j] = d[(i-1)*width+(j-1)]
			} else {
				d[i*width+j] = Min(d[(i-1)*width+j]+1, //deletion
					d[i*width+(j-1)]+1,     //insertion
					d[(i-1)*width+(j-1)]+1) //substitution
			}
		}
	}
	return d[m*(width)+0]
}
