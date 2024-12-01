package utils

import (
	"encoding/json"
	"io"
	"strings"
)

func GetReaderFromStruct(s interface{}) io.Reader {
	res, err := json.Marshal(s)

	if err != nil {
		panic(err)
	}

	return strings.NewReader(string(res))
}
