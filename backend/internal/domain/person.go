package domain

import "errors"

var ErrPersonNotFound = errors.New("person not found")

type Person struct {
	ID      int64
	Name    string
	Age     int32
	Address string
	Work    string
}
