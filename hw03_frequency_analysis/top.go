package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type Pair struct {
	Key   string
	Value int
}

var re = regexp.MustCompile(`[.,()!?^&*'";:/\\|-]+`)

func Top10(s string) []string {
	// get top 10 words from string
	data := strings.Fields(s)

	if len(data) == 0 {
		return []string{}
	}

	wordFrequency := make(map[string]int)
	for _, word := range data {
		// Remove unwanted characters like .,()- from the word
		word = re.ReplaceAllString(word, "")
		if len(word) == 0 {
			continue
		}
		wordFrequency[strings.ToLower(word)]++
	}

	pairs := make([]Pair, 0, len(wordFrequency))
	for key, value := range wordFrequency {
		pairs = append(pairs, Pair{key, value})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].Value == pairs[j].Value {
			return pairs[i].Key < pairs[j].Key
		}
		return pairs[i].Value > pairs[j].Value
	})

	var result []string

	outputSize := 10
	if len(pairs) < 10 {
		outputSize = len(pairs)
	}

	for i := 0; i < outputSize; i++ {
		result = append(result, pairs[i].Key)
	}

	return result
}
