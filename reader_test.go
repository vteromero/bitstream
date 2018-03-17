package bitstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader_NewReader(t *testing.T) {
	b := []byte{0x11, 0x22}
	r := NewReader(b)
	assert.Equal(t, b, r.b)
	assert.Equal(t, 0, r.i)
	assert.NotNil(t, r.buff)
}

func TestReader_Read(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
	}
	params := []struct {
		n   int
		v   uint64
		m   int
		err error
	}{
		{64, 0x0706050403020100, 64, nil},
		{0, 0x0, 0, nil},
		{-100, 0x0, 0, ErrSizeOutOfBound},
		{100, 0x0, 0, ErrSizeOutOfBound},
		{64, 0x0f0e0d0c0b0a0908, 64, nil},
		{7, 0x10, 7, nil},
		{3, 0x02, 3, nil},
		{6, 0x04, 6, nil},
		{0, 0x0, 0, nil},
		{29, 0x15141312, 29, nil},
		{18, 0xb8b0, 18, nil},
		{8, 0, 1, EOF},
		{8, 0, 0, EOF},
		{8, 0, 0, EOF},
	}

	r := NewReader(data)

	for i := 0; i < len(params); i++ {
		x, y, z := r.Read(params[i].n)

		assert.Equal(t, params[i].v, x)
		assert.Equal(t, params[i].m, y)
		assert.Equal(t, params[i].err, z)
	}
}

func TestReader_ReadAt(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
	}
	params := []struct {
		n   int
		off int
		v   uint64
		m   int
		err error
	}{
		{0, 20, 0x0, 0, nil},
		{8, 8, 0x01, 8, nil},
		{8, 40, 0x05, 8, nil},
		{8, 96, 0x0c, 8, nil},
		{8, 160, 0x14, 8, nil},
		{32, 56, 0x0a090807, 32, nil},
		{32, 80, 0x0d0c0b0a, 32, nil},
		{64, 64, 0x0f0e0d0c0b0a0908, 64, nil},
		{64, 160, 0x17161514, 32, EOF},
		{16, 184, 0x17, 8, EOF},
		{-250, 0, 0x0, 0, ErrSizeOutOfBound},
		{250, 0, 0x0, 0, ErrSizeOutOfBound},
		{8, -10, 0x0, 0, ErrOffsetOutOfBound},
		{8, len(data)*8 + 10, 0x0, 0, ErrOffsetOutOfBound},
	}

	r := NewReader(data)

	for i := 0; i < len(params); i++ {
		x, y, z := r.ReadAt(params[i].n, params[i].off)

		assert.Equal(t, params[i].v, x)
		assert.Equal(t, params[i].m, y)
		assert.Equal(t, params[i].err, z)
	}
}

func TestReader_Reset(t *testing.T) {
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	}

	r := NewReader(data)

	x, y, z := r.Read(64)
	assert.Equal(t, uint64(0x0706050403020100), x)
	assert.Equal(t, 64, y)
	assert.Nil(t, z)

	r.Reset()
	assert.Equal(t, 64, y)

	x, y, z = r.Read(64)
	assert.Equal(t, uint64(0x0706050403020100), x)
	assert.Nil(t, z)
}
