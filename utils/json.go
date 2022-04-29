package utils

import (
	"encoding/json"
	"log"
)

func JsonMarshal(v any) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		log.Panicln("json marshal error")
	}
	return result
}

func JsonUnmarshal[T any](data []byte) T {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Panicln("json unmarshal error: ", err)
	}
	return result
}
