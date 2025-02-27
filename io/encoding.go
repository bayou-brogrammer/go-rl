package io

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
)

func DecodeSave[T any](data []byte) (*T, error) {
	buf := bytes.NewReader(data)
	r, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}

	lg := new(T)
	dec := gob.NewDecoder(r)
	err = dec.Decode(lg)
	if err != nil {
		return nil, err
	}

	r.Close()
	return lg, nil
}
