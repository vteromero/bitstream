// Copyright (c) 2018 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader_NewReader(t *testing.T) {
	b := []byte{0x11, 0x22}
	r := NewReader(b)
	assert.Equal(t, b, r.b)
	assert.Equal(t, 0, r.off)
	assert.Equal(t, 16, r.maxOffset)
}

func TestReader_Read(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
	}
	params := []struct {
		n            int
		expectedBits uint64
		expectedErr  error
	}{
		{-100, 0x0, ErrSizeOutOfBound},
		{100, 0x0, ErrSizeOutOfBound},
		{58, 0x0, ErrSizeOutOfBound},
		{32, 0x03020100, nil},
		{32, 0x07060504, nil},
		{8, 0x08, nil},
		{0, 0x0, nil},
		{0, 0x0, nil},
		{57, 0x00f0e0d0c0b0a09, nil},
		{3, 0x0, nil},
		{4, 0x1, nil},
		{6, 0x11, nil},
		{6, 0x08, nil},
		{4, 0x1, nil},
		{40, 0x1716151413, nil},
		{0, 0x0, EOF},
		{1, 0x0, EOF},
		{10, 0x0, EOF},
	}

	r := NewReader(data)

	for i := 0; i < len(params); i++ {
		x, y := r.Read(params[i].n)

		assert.Equal(t, params[i].expectedBits, x)
		assert.Equal(t, params[i].expectedErr, y)
	}
}

func TestReader_ReadAt(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
	}
	params := []struct {
		n            int
		off          int
		expectedBits uint64
		expectedErr  error
	}{
		{0, 20, 0x0, nil},
		{8, 8, 0x01, nil},
		{8, 40, 0x05, nil},
		{8, 96, 0x0c, nil},
		{8, 160, 0x14, nil},
		{20, 0, 0x20100, nil},
		{32, 56, 0x0a090807, nil},
		{32, 80, 0x0d0c0b0a, nil},
		{56, 64, 0x0e0d0c0b0a0908, nil},
		{56, 136, 0x17161514131211, nil},
		{16, 184, 0x0, EOF},
		{-250, 0, 0x0, ErrSizeOutOfBound},
		{250, 0, 0x0, ErrSizeOutOfBound},
		{58, 0, 0x0, ErrSizeOutOfBound},
		{8, -10, 0x0, ErrOffsetOutOfBound},
		{8, len(data) * 8, 0x0, ErrOffsetOutOfBound},
		{8, len(data)*8 + 10, 0x0, ErrOffsetOutOfBound},
	}

	r := NewReader(data)

	for i := 0; i < len(params); i++ {
		x, y := r.ReadAt(params[i].n, params[i].off)

		assert.Equal(t, params[i].expectedBits, x)
		assert.Equal(t, params[i].expectedErr, y)
	}
}

func TestReader_Reset(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06,
	}

	r := NewReader(data)

	x, y := r.Read(56)
	assert.Equal(t, uint64(0x06050403020100), x)
	assert.Nil(t, y)

	r.Reset()

	x, y = r.Read(56)
	assert.Equal(t, uint64(0x06050403020100), x)
	assert.Nil(t, y)
}

func TestReader_Offset(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	}

	r := NewReader(data)
	assert.Equal(t, 0, r.Offset())

	r.Read(55)
	r.Read(50)
	assert.Equal(t, 105, r.Offset())

	r.Reset()
	assert.Equal(t, 0, r.Offset())

	r.ReadAt(20, 100)
	assert.Equal(t, 120, r.Offset())

	r.ReadAt(50, 0)
	assert.Equal(t, 50, r.Offset())
}
