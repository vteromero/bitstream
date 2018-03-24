// Copyright (c) 2018 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriter_NewWriter(t *testing.T) {
	b := []byte{0x0, 0x0}
	w := NewWriter(b)
	assert.Equal(t, b, w.b)
	assert.Equal(t, 0, w.i)
	assert.NotNil(t, w.buff)
}

func TestWriter_Write(t *testing.T) {
	params := []struct {
		v   uint64
		n   int
		m   int
		err error
	}{
		{0x0, -100, 0, ErrSizeOutOfBound},
		{0x0, 100, 0, ErrSizeOutOfBound},
		{0x44, 8, 8, nil},
		{0x2, 4, 4, nil},
		{0xc, 4, 4, nil},
		{0x5, 3, 3, nil},
		{0x1b, 5, 5, nil},
		{0x49, 7, 7, nil},
		{0x1, 1, 1, nil},
		{0x0, 0, 0, nil},
		{0x0, 0, 0, nil},
		{0xaabbccdd, 32, 32, nil},
		{0xff, 8, 0, EOF},
	}
	expected := []byte{0x44, 0xc2, 0xdd, 0xc9, 0xdd, 0xcc, 0xbb, 0xaa}

	data := make([]byte, len(expected))
	w := NewWriter(data)

	for i := 0; i < len(params); i++ {
		m, err := w.Write(params[i].v, params[i].n)

		assert.Equal(t, params[i].m, m)
		assert.Equal(t, params[i].err, err)
	}

	w.Close()

	assert.Equal(t, expected, data)
}

func TestWriter_Close(t *testing.T) {
	data := make([]byte, 2)
	w := NewWriter(data)
	w.Write(0xff, 8)
	assert.Equal(t, []byte{0x0, 0x0}, data)

	w.Close()
	assert.Equal(t, []byte{0xff, 0x0}, data)

	w.Close()
	assert.Equal(t, []byte{0xff, 0x0}, data)
}

func TestWriter_Offset(t *testing.T) {
	data := make([]byte, 20)
	w := NewWriter(data)
	assert.Equal(t, 0, w.Offset())

	w.Write(0x11223344556677, 60)
	w.Write(0xffffff, 20)
	w.Write(0xffff, 10)
	w.Close()
	assert.Equal(t, 90, w.Offset())
}

func TestWriter_Reset(t *testing.T) {
	data := make([]byte, 2)
	expected := []byte{0x11, 0x11}

	w := NewWriter(data)

	w.Write(0xffff, 16)
	w.Close()

	w.Reset()

	w.Write(0x1111, 16)
	w.Close()

	assert.Equal(t, expected, data)
}
