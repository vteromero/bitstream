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
	assert.Equal(t, 0, w.off)
}

func TestWriter_Write(t *testing.T) {
	params := []struct {
		v           uint64
		n           int
		expectedErr error
	}{
		{0x0, -100, ErrSizeOutOfBound},
		{0x0, 100, ErrSizeOutOfBound},
		{0x44, 8, nil},
		{0x2, 4, nil},
		{0xc, 4, nil},
		{0x5, 3, nil},
		{0x1b, 5, nil},
		{0x49, 7, nil},
		{0x1, 1, nil},
		{0x0, 0, nil},
		{0x0, 0, nil},
		{0xaabbccdd, 32, nil},
		{0x0, 0, EOF},
		{0x0, 0, EOF},
		{0xffff, 16, EOF},
		{0x1111, 16, EOF},
	}
	expected := []byte{0x44, 0xc2, 0xdd, 0xc9, 0xdd, 0xcc, 0xbb, 0xaa}

	data := make([]byte, len(expected))
	w := NewWriter(data)

	for i := 0; i < len(params); i++ {
		err := w.Write(params[i].v, params[i].n)
		assert.Equal(t, params[i].expectedErr, err)
	}

	assert.Equal(t, expected, data)
}

func TestWriter_Offset(t *testing.T) {
	data := make([]byte, 20)
	w := NewWriter(data)
	assert.Equal(t, 0, w.Offset())

	w.Write(0x11223344556677, 50)
	w.Write(0xffffff, 20)
	w.Write(0xffff, 10)
	assert.Equal(t, 80, w.Offset())
}

func TestWriter_Reset(t *testing.T) {
	data := make([]byte, 2)
	w := NewWriter(data)

	w.Write(0xffff, 16)
	assert.Equal(t, 16, w.Offset())

	w.Reset()
	assert.Equal(t, 0, w.Offset())
}
