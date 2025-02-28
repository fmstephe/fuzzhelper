// Copyright 2025 Francis Michael Stephens. All rights reserved.  Use of this
// source code is governed by an MIT license that can be found in the LICENSE
// file.

package fuzzhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// It turns out that pushing and then reading the same negative integer for sizes less than 64 bits is subtle.
// So we write a unit test for it here.
func TestByteConsumer_PushAndGetInt64(t *testing.T) {
	consumer := NewByteConsumer([]byte{})

	consumer.pushInt64(-1, NativeBytes)
	assert.Equal(t, int64(-1), consumer.Int64(NativeBytes))
	//
	consumer.pushInt64(-1, BytesFor64)
	assert.Equal(t, int64(-1), consumer.Int64(BytesFor64))
	//
	consumer.pushInt64(-1, 4)
	assert.Equal(t, int64(-1), consumer.Int64(4))
	//
	consumer.pushInt64(-1, BytesFor16)
	assert.Equal(t, int64(-1), consumer.Int64(BytesFor16))
	//
	consumer.pushInt64(-1, BytesFor8)
	assert.Equal(t, int64(-1), consumer.Int64(BytesFor8))
}

func TestByteConsumer_Bytes(t *testing.T) {
	consumer := NewByteConsumer([]byte{})
	consumer.pushBytes([]byte{1, 2, 3, 4, 5, 6})
	consumer.pushByte(7)
	assert.Equal(t, 7, consumer.Len())

	// Consume the available bytes
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6}, consumer.Bytes(6))
	assert.Equal(t, 1, consumer.Len())

	// Consume bytes, but not enoough available - get remaining bytes and zeroes
	assert.Equal(t, []byte{7, 0, 0, 0, 0, 0}, consumer.Bytes(6))
	assert.Equal(t, 0, consumer.Len())

	// Consume bytes, but none available - get zeroes
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0}, consumer.Bytes(6))
	assert.Equal(t, 0, consumer.Len())
}

func TestByteConsumer_Byte(t *testing.T) {
	consumer := NewByteConsumer([]byte{})
	consumer.pushByte(12)
	assert.Equal(t, 1, consumer.Len())

	// Consume the available bytes
	assert.Equal(t, byte(12), consumer.Byte())
	assert.Equal(t, 0, consumer.Len())

	// Consume bytes, but none available - get zeroes
	assert.Equal(t, byte(0), consumer.Byte())
	assert.Equal(t, 0, consumer.Len())
}

func TestByteConsumer_Uint16(t *testing.T) {
	consumer := NewByteConsumer([]byte{})
	consumer.pushUint64(10_000, 2)
	consumer.pushByte(7)
	assert.Equal(t, 3, consumer.Len())

	// Consume the available bytes
	assert.Equal(t, uint16(10_000), uint16(consumer.Uint64(2)))
	assert.Equal(t, 1, consumer.Len())

	// Consume bytes, but not enoough available - get remaining bytes and zeroes
	assert.Equal(t, uint16(7), uint16(consumer.Uint64(2)))
	assert.Equal(t, 0, consumer.Len())

	// Consume bytes, but none available - get zeroes
	assert.Equal(t, uint16(0), uint16(consumer.Uint64(2)))
	assert.Equal(t, 0, consumer.Len())
}

func TestByteConsumer_Uint32(t *testing.T) {
	consumer := NewByteConsumer([]byte{})
	consumer.pushUint64(100_000, 4)
	consumer.pushByte(7)
	assert.Equal(t, 5, consumer.Len())

	// Consume the available bytes
	assert.Equal(t, uint32(100_000), uint32(consumer.Uint64(4)))
	assert.Equal(t, 1, consumer.Len())

	// Consume bytes, but not enoough available - get remaining bytes and zeroes
	assert.Equal(t, uint32(7), uint32(consumer.Uint64(4)))
	assert.Equal(t, 0, consumer.Len())

	// Consume bytes, but none available - get zeroes
	assert.Equal(t, uint32(0), uint32(consumer.Uint64(4)))
	assert.Equal(t, 0, consumer.Len())
}

func TestByteConsumer_Combined(t *testing.T) {
	consumer := NewByteConsumer([]byte{})
	consumer.pushBytes([]byte{1, 2, 3, 4, 5, 6})
	consumer.pushByte(12)
	consumer.pushUint64(10_000, 2)
	consumer.pushUint64(100_000, 4)
	assert.Equal(t, 13, consumer.Len())

	// Consume the available bytes
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6}, consumer.Bytes(6))
	assert.Equal(t, 7, consumer.Len())

	assert.Equal(t, byte(12), consumer.Byte())
	assert.Equal(t, 6, consumer.Len())

	assert.Equal(t, uint16(10_000), uint16(consumer.Uint64(2)))
	assert.Equal(t, 4, consumer.Len())

	assert.Equal(t, uint32(100_000), uint32(consumer.Uint64(4)))
	assert.Equal(t, 0, consumer.Len())
}
