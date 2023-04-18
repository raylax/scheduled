package store

import (
	"github.com/raylax/scheduled/codec"
	"io"
)

type Data struct {
	m map[string]string
}

func NewData() *Data {
	return &Data{
		m: map[string]string{},
	}
}

func (d *Data) Persist(w io.Writer) error {
	encoder := codec.NewEncoder(w)
	return encoder.Encode(d)
}

func RestoreData(r io.ReadCloser) (*Data, error) {
	data := &Data{}
	decoder := codec.NewDecoder(r)
	err := decoder.Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
