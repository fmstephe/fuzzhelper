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
	consumer := newByteConsumer([]byte{})

	consumer.pushInt64(-1, bytesForNative)
	assert.Equal(t, int64(-1), consumer.consumeInt64(bytesForNative))
	//
	consumer.pushInt64(-1, bytesFor64)
	assert.Equal(t, int64(-1), consumer.consumeInt64(bytesFor64))
	//
	consumer.pushInt64(-1, 4)
	assert.Equal(t, int64(-1), consumer.consumeInt64(4))
	//
	consumer.pushInt64(-1, bytesFor16)
	assert.Equal(t, int64(-1), consumer.consumeInt64(bytesFor16))
	//
	consumer.pushInt64(-1, bytesFor8)
	assert.Equal(t, int64(-1), consumer.consumeInt64(bytesFor8))
}

func TestByteConsumer_Bytes(t *testing.T) {
	consumer := newByteConsumer([]byte{})
	consumer.pushBytes([]byte{1, 2, 3, 4, 5, 6})
	consumer.pushByte(7)
	assert.Equal(t, 7, consumer.len())

	// Consume the available bytes
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6}, consumer.consume(6))
	assert.Equal(t, 1, consumer.len())

	// Consume bytes, but not enoough available - get remaining bytes and zeroes
	assert.Equal(t, []byte{7, 0, 0, 0, 0, 0}, consumer.consume(6))
	assert.Equal(t, 0, consumer.len())

	// Consume bytes, but none available - get zeroes
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0}, consumer.consume(6))
	assert.Equal(t, 0, consumer.len())
}

func TestByteConsumer_Byte(t *testing.T) {
	consumer := newByteConsumer([]byte{})
	consumer.pushByte(12)
	assert.Equal(t, 1, consumer.len())

	// Consume the available bytes
	assert.Equal(t, byte(12), consumer.singleByte())
	assert.Equal(t, 0, consumer.len())

	// Consume bytes, but none available - get zeroes
	assert.Equal(t, byte(0), consumer.singleByte())
	assert.Equal(t, 0, consumer.len())
}

func TestByteConsumer_Uint16(t *testing.T) {
	consumer := newByteConsumer([]byte{})
	consumer.pushUint64(10_000, 2)
	consumer.pushByte(7)
	assert.Equal(t, 3, consumer.len())

	// Consume the available bytes
	assert.Equal(t, uint16(10_000), uint16(consumer.consumeUint64(2)))
	assert.Equal(t, 1, consumer.len())

	// Consume bytes, but not enoough available - get remaining bytes and zeroes
	assert.Equal(t, uint16(7), uint16(consumer.consumeUint64(2)))
	assert.Equal(t, 0, consumer.len())

	// Consume bytes, but none available - get zeroes
	assert.Equal(t, uint16(0), uint16(consumer.consumeUint64(2)))
	assert.Equal(t, 0, consumer.len())
}

func TestByteConsumer_Uint32(t *testing.T) {
	consumer := newByteConsumer([]byte{})
	consumer.pushUint64(100_000, 4)
	consumer.pushByte(7)
	assert.Equal(t, 5, consumer.len())

	// Consume the available bytes
	assert.Equal(t, uint32(100_000), uint32(consumer.consumeUint64(4)))
	assert.Equal(t, 1, consumer.len())

	// Consume bytes, but not enoough available - get remaining bytes and zeroes
	assert.Equal(t, uint32(7), uint32(consumer.consumeUint64(4)))
	assert.Equal(t, 0, consumer.len())

	// Consume bytes, but none available - get zeroes
	assert.Equal(t, uint32(0), uint32(consumer.consumeUint64(4)))
	assert.Equal(t, 0, consumer.len())
}

func TestByteConsumer_Combined(t *testing.T) {
	consumer := newByteConsumer([]byte{})
	consumer.pushBytes([]byte{1, 2, 3, 4, 5, 6})
	consumer.pushByte(12)
	consumer.pushUint64(10_000, 2)
	consumer.pushUint64(100_000, 4)
	assert.Equal(t, 13, consumer.len())

	// Consume the available bytes
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6}, consumer.consume(6))
	assert.Equal(t, 7, consumer.len())

	assert.Equal(t, byte(12), consumer.singleByte())
	assert.Equal(t, 6, consumer.len())

	assert.Equal(t, uint16(10_000), uint16(consumer.consumeUint64(2)))
	assert.Equal(t, 4, consumer.len())

	assert.Equal(t, uint32(100_000), uint32(consumer.consumeUint64(4)))
	assert.Equal(t, 0, consumer.len())
}
