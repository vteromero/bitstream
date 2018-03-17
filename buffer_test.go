package bitstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBitBuffer(t *testing.T) {
	buff := newBitBuffer(20)
	assert.Equal(t, uint64(0), buff.bits)
	assert.Equal(t, 20, buff.len)
	assert.Equal(t, 0, buff.off)
}

func TestBitBuffer_LoadFrom(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	expected := []struct {
		bits uint64
		len  int
		off  int
		n    int
	}{
		{0x0706050403020100, 64, 0, 8},
		{0x07060504030201, 56, 0, 7},
		{0x070605040302, 48, 0, 6},
		{0x0706050403, 40, 0, 5},
		{0x07060504, 32, 0, 4},
		{0x070605, 24, 0, 3},
		{0x0706, 16, 0, 2},
		{0x07, 8, 0, 1},
		{0x0, 0, 0, 0},
	}

	for i := 0; i < len(expected); i++ {
		buff := newBitBuffer(0)
		n := buff.loadFrom(data[i:])

		assert.Equal(t, expected[i].bits, buff.bits)
		assert.Equal(t, expected[i].len, buff.len)
		assert.Equal(t, expected[i].off, buff.off)
		assert.Equal(t, expected[i].n, n)
	}
}

func TestBitBuffer_WriteTo(t *testing.T) {
	expected := []struct {
		bits      uint64
		len       int
		sliceSize int
		result    []byte
	}{
		{0x0706050403020100, 64, 16, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}},
		{0x0706050403020100, 64, 8, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}},
		{0x0706050403020100, 64, 7, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}},
		{0x0706050403020100, 64, 6, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}},
		{0x0706050403020100, 64, 5, []byte{0x00, 0x01, 0x02, 0x03, 0x04}},
		{0x0706050403020100, 64, 4, []byte{0x00, 0x01, 0x02, 0x03}},
		{0x0706050403020100, 64, 3, []byte{0x00, 0x01, 0x02}},
		{0x0706050403020100, 64, 2, []byte{0x00, 0x01}},
		{0x0706050403020100, 64, 1, []byte{0x00}},
		{0x0706050403020100, 64, 0, []byte{}},
		{0x0706050403020100, 0, 8, []byte{}},
		{0x0706050403020100, 8, 8, []byte{0x00}},
		{0x0706050403020100, 16, 8, []byte{0x00, 0x01}},
		{0x0706050403020100, 24, 8, []byte{0x00, 0x01, 0x02}},
		{0x0706050403020100, 32, 8, []byte{0x00, 0x01, 0x02, 0x03}},
		{0x0706050403020100, 40, 8, []byte{0x00, 0x01, 0x02, 0x03, 0x04}},
		{0x0706050403020100, 48, 8, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}},
		{0x0706050403020100, 56, 8, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}},
		{0x0706050403020100, 64, 8, []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}},
	}

	for i := 0; i < len(expected); i++ {
		data := make([]byte, expected[i].sliceSize)
		buff := newBitBuffer(64)
		buff.bits = expected[i].bits
		buff.len = expected[i].len
		n := buff.writeTo(data)

		assert.Equal(t, len(expected[i].result), n)
		assert.Equal(t, expected[i].result, data[:len(expected[i].result)])
	}
}

func TestBitBuffer_Reset(t *testing.T) {
	buff := newBitBuffer(0)
	buff.loadFrom([]byte{0xff, 0xff})
	buff.reset()

	assert.Equal(t, uint64(0), buff.bits)
	assert.Equal(t, 0, buff.len)
	assert.Equal(t, 0, buff.off)
}

func TestBitBuffer_Get(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	params := []struct {
		n int
		v uint64
		m int
	}{
		{4, 0x0, 4},
		{4, 0x0, 4},
		{4, 0x1, 4},
		{4, 0x0, 4},
		{8, 0x02, 8},
		{8, 0x03, 8},
		{3, 0x4, 3},
		{13, 0xa0, 13},
		{30, 0x0706, 16},
		{10, 0x0, 0},
	}

	buff := newBitBuffer(0)
	buff.loadFrom(data)

	for i := 0; i < len(params); i++ {
		x, y := buff.get(params[i].n)

		assert.Equal(t, params[i].v, x)
		assert.Equal(t, params[i].m, y)
	}
}

func TestBitBuffer_Set(t *testing.T) {
	params := []struct {
		v uint64
		n int
		m int
	}{
		{0xff, 8, 8},
		{0xaa, 8, 8},
		{0x1, 1, 1},
		{0x0, 1, 1},
		{0x1, 1, 1},
		{0x0, 1, 1},
		{0xf, 4, 4},
		{0x7762, 16, 16},
		{0x90807060, 32, 24},
		{0xff, 8, 0},
	}
	expected := &bitBuffer{
		bits: 0x8070607762f5aaff,
		len:  64,
		off:  64,
	}

	buff := newBitBuffer(64)

	for i := 0; i < len(params); i++ {
		m := buff.set(params[i].v, params[i].n)

		assert.Equal(t, params[i].m, m)
	}

	assert.Equal(t, expected, buff)
}
