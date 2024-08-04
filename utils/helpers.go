package utils

import (
	"encoding/json"
)

func JsonSerilizeDeserialize(src any, dest any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return err
	}

	return nil
}

func GetBytes(d any) ([]byte, error) {
	return json.Marshal(d)
	// var buf bytes.Buffer
	// enc := gob.NewEncoder(&buf)
	// err := enc.Encode(d)
	// if err != nil {
	// 	return nil, err
	// }
	// return buf.Bytes(), nil
}
