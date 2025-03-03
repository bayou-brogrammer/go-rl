package io

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/gob"
	"log"
	"os"
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

func DecodeBinary(filename string) (*gzip.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	r, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer r.Close()

	return r, nil
}
