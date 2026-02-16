package main

import (
	"encoding/json"
	"io"
)

type Student struct {
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

func DecodeStudentFromReader(r io.Reader) (s Student, err error) {
	err = json.NewDecoder(r).Decode(&s)
	return s, err
}
