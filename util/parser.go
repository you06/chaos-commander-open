package util

import "encoding/json"

func ParseStringSlice(input string) []string {
	var s []string
	err := json.Unmarshal([]byte(input), &s)
	if err != nil {
		return []string{}
	}
	return s
}

func ParseIntSlice(input string) []int {
	var s []int
	err := json.Unmarshal([]byte(input), &s)
	if err != nil {
		return []int{}
	}
	return s
}
