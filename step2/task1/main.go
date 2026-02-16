package main

import "encoding/json"

func SerializeIntSlice(numbers []int) ([]byte, error) {
	return json.Marshal(numbers)
}
