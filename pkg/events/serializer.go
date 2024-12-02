package events

import (
	"encoding/json"
)

func UnsafeDeserializeAny[T any](b []byte) *T {
	out := new(T)
	err := json.Unmarshal(b, &out)
	if err != nil {
		panic(err)
	}
	return out
}

func UnsafeSerializeAny[T any](t T) []byte {

	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return b
}
