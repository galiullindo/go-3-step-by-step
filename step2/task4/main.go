package main

import (
	"encoding/json"
	"io"
)

type Student struct {
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

func EncodeStudentsToWriter(w io.Writer, students []Student) (err error) {
	err = json.NewEncoder(w).Encode(&students)
	return err
}
