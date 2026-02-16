package main

import "encoding/json"

func DeserializeStringMap(data string) (map[string]string, error) {
	var m map[string]string
	err := json.Unmarshal([]byte(data), &m)
	return m, err
}
