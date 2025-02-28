// Copyright 2025 Francis Michael Stephens. All rights reserved.  Use of this
// source code is governed by an MIT license that can be found in the LICENSE
// file.

package fuzzhelper

import (
	"encoding/binary"
	"math"
	"unicode/utf8"
	"unsafe"
)

var nativeInt = 0

const (
	BytesForNative = uintptr(unsafe.Sizeof(nativeInt))
	BytesFor64     = 8
	BytesFor32     = 4
	BytesFor16     = 2
	BytesFor8      = 1
)

type ByteConsumer struct {
	bytes []byte
}

func NewByteConsumer(bytes []byte) *ByteConsumer {
	return &ByteConsumer{
		bytes: bytes,
	}
}

func (c *ByteConsumer) Len() int {
	return len(c.bytes)
}

func (c *ByteConsumer) Bytes(size int) []byte {
	consumed := make([]byte, size)
	copy(consumed, c.bytes)

	if len(c.bytes) <= size {
		c.bytes = c.bytes[:0]
	} else {
		c.bytes = c.bytes[size:]
	}
	return consumed
}

// Test only
func (c *ByteConsumer) pushBytes(bytes []byte) {
	c.bytes = append(c.bytes, bytes...)
}

func (c *ByteConsumer) Byte() byte {
	dest := c.Bytes(1)
	return dest[0]
}

// Test only
func (c *ByteConsumer) pushByte(b byte) {
	c.pushBytes([]byte{b})
}

func (c *ByteConsumer) Uint64(bytes uintptr) uint64 {
	switch bytes {
	case 8:
		dest := c.Bytes(8)
		return binary.LittleEndian.Uint64(dest)
	case 4:
		dest := c.Bytes(4)
		return uint64(binary.LittleEndian.Uint32(dest))
	case 2:
		dest := c.Bytes(2)
		return uint64(binary.LittleEndian.Uint16(dest))
	case 1:
		dest := c.Bytes(1)
		return uint64(dest[0])
	default:
		panic("Must provided either 8, 4, 2, or 1 as bytes argument")
	}
}

func (c *ByteConsumer) Int64(bytes uintptr) int64 {
	switch bytes {
	case 8:
		dest := c.Bytes(8)
		return int64(binary.LittleEndian.Uint64(dest))
	case 4:
		dest := c.Bytes(4)
		return int64(int32((binary.LittleEndian.Uint32(dest))))
	case 2:
		dest := c.Bytes(2)
		return int64(int16((binary.LittleEndian.Uint16(dest))))
	case 1:
		dest := c.Bytes(1)
		return int64(int8((dest[0])))
	default:
		panic("Must provided either 8, 4, 2, or 1 as bytes argument")
	}
}

func (c *ByteConsumer) Float64(bytes uintptr) float64 {
	switch bytes {
	case 8:
		dest := c.Bytes(8)
		return math.Float64frombits(binary.LittleEndian.Uint64(dest))
	case 4:
		dest := c.Bytes(4)
		return float64(math.Float32frombits(binary.LittleEndian.Uint32(dest)))
	default:
		panic("Must provided either 8 or 4 as bytes argument")
	}
}

// Returns a valid UTF8 string, which is at _most_ as long as length in runes
func (c *ByteConsumer) String(length int) string {
	bytes := c.Bytes(length)

	// Extract the valid runes from the bytes
	//
	// The way this is implemented right now will create strings which are
	// randomly shorter than asked for. We may want to correct this.
	validRunes := []rune{}
	for len(bytes) > 0 {
		r, l := utf8.DecodeRune(bytes)
		bytes = bytes[l:]
		if r != utf8.RuneError {
			validRunes = append(validRunes, r)
		}
	}

	return string(validRunes)
}

func (c *ByteConsumer) Bool() bool {
	bytes := c.Bytes(1)
	return bytes[0]%2 == 1
}

// test only
func (c *ByteConsumer) pushUint64(value uint64, bytes uintptr) {
	switch bytes {
	case 8:
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, value)
		c.pushBytes(bytes)
	case 4:
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, uint32(value))
		c.pushBytes(bytes)
	case 2:
		bytes := make([]byte, 2)
		binary.LittleEndian.PutUint16(bytes, uint16(value))
		c.pushBytes(bytes)
	case 1:
		bytes := make([]byte, 1)
		bytes[0] = byte(value)
		c.pushBytes(bytes)
	default:
		panic("Must provided either 8, 4, 2, or 1 as bytes argument")
	}
}

// test only
func (c *ByteConsumer) pushInt64(value int64, bytes uintptr) {
	switch bytes {
	case 8:
		c.pushUint64(uint64(value), bytes)
	case 4:
		c.pushUint64(uint64(int32(value)), bytes)
	case 2:
		c.pushUint64(uint64(int16(value)), bytes)
	case 1:
		c.pushUint64(uint64(int8(value)), bytes)
	default:
		panic("Must provided either 8, 4, 2, or 1 as bytes argument")
	}
}

func (c *ByteConsumer) pushFloat64(value float64, bytes uintptr) {
	switch bytes {
	case 8:
		floatBits := math.Float64bits(value)
		c.pushUint64(floatBits, 8)
	case 4:
		floatBits := uint64(math.Float32bits(float32(value)))
		c.pushUint64(floatBits, 4)
	default:
		panic("Must provided either 8 or 4 as bytes argument")
	}
}

// test only
func (c *ByteConsumer) pushString(str string) {
	c.pushBytes([]byte(str))
}

func (c *ByteConsumer) pushBool(value bool) {
	if value {
		c.pushBytes([]byte{1})
	} else {
		c.pushBytes([]byte{0})
	}
}
