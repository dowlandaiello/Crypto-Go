package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Hash - hash specified interface, return string
func Hash(obj interface{}) (string, error) {
	b, err := GetBytes(obj)

	if err != nil {
		return "", err
	}

	h := sha256.Sum256(b)

	return string(h[:]), nil
}

// GetBytes - get bytes of specified interface, return byte array
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
