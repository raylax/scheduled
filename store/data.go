package store

import (
	"github.com/raylax/scheduled/codec"
	"io"
)

type Data map[string]string

func (data Data) Set(key string, val string) {
	data[key] = val
}

func NewData() *Data {
	return &Data{}
}

func (data Data) Persist(w io.Writer) error {
	encoder := codec.NewEncoder(w)
	return encoder.Encode(data)
}

func (data Data) Restore(r io.ReadCloser) error {
	decoder := codec.NewDecoder(r)
	return decoder.Decode(&data)
}
