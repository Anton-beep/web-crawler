package utils

import (
	"encoding/json"
	"io"
	"strings"
)

// GetReaderFromStruct returns a reader from a struct
func GetReaderFromStruct(s interface{}) io.Reader {
	res, err := json.Marshal(s)

	if err != nil {
		panic(err)
	}

	return strings.NewReader(string(res))
}
