// Copyright 2018 Noskov Artem, Alex Browne. All rights reserved.
// Use of this source code is governed by the MIT
// license, which can be found in the LICENSE file.

package store

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

var (
	// BinaryEncoding is a ready-to-use implementation of EncoderDecoder which
	// encodes data structures in a binary format using the gob package.
	BinaryEncoding = &binaryEncoderDecoder{}
	// JSONEncoding is a ready-to-use implementation of EncoderDecoder which
	// encodes data structures as json.
	JSONEncoding = &jsonEncoderDecoder{}
)

// Encoder is an interface implemented by objects which can encode an arbitrary
// go object into a slice of bytes.
type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

// Decoder is an interface implemented by objects which can decode a slice
// of bytes into an arbitrary go object.
type Decoder interface {
	Decode([]byte, interface{}) error
}

// EncoderDecoder is an interface implemented by objects which can both encode
// an arbitrary go object into a slice of bytes and decode that slice of bytes
// into an arbitrary go object. EncoderDecoders should have the property that
// Encode(Decode(x)) == x for all objects x which are encodable.
type EncoderDecoder interface {
	Encoder
	Decoder
}

// jsonEncoderDecoder is an implementation of EncoderDecoder which uses json
// encoding.
type jsonEncoderDecoder struct{}

// Encode implements the Encode method of Encoder
func (jsonEncoderDecoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Decode implements the Decode method of Decoder
func (jsonEncoderDecoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// binaryEncoderDecoder is an implementation of EncoderDecoder which uses binary
// encoding via the gob package in the standard library.
type binaryEncoderDecoder struct{}

// Encode implements the Encode method of Encoder
func (b binaryEncoderDecoder) Encode(v interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode implements the Decode method of Decoder
func (b binaryEncoderDecoder) Decode(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(v)
}
