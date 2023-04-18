package codec

import (
	"github.com/vmihailenco/msgpack/v5"
	"io"
)

type Encoder interface {
	Encode(v any) error
}

type Decoder interface {
	Decode(v any) error
}

func NewEncoder(w io.Writer) Encoder {
	return msgpack.NewEncoder(w)
}

func NewDecoder(r io.Reader) Decoder {
	return msgpack.NewDecoder(r)
}
