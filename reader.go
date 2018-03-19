// Copyright (c) 2018 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

// A Reader is a structure to read bits from a stream of bytes.
type Reader struct {
	b    []byte
	i    int
	buff *bitBuffer
}

// NewReader returns a Reader that can read bits from the byte slice provided.
func NewReader(b []byte) *Reader {
	return &Reader{
		b:    b,
		i:    0,
		buff: newBitBuffer(0),
	}
}

// Read reads the next n bits from the stream.
// It returns 3 values: a 64-bit integer that holds the bits, the actual number
// of bits read and an error value.
// The returning number of bits is equal to n if there have been no errors and
// there are enough bits to read. Otherwise, the value will be less than n.
// The returning error is not nil when something went wrong or when the end of
// the stream has been reached (EOF error).
func (r *Reader) Read(n int) (uint64, int, error) {
	if n < 0 || n > 64 {
		return 0, 0, ErrSizeOutOfBound
	}
	v, m := r.buff.get(n)
	if m < n {
		r.i += r.buff.loadFrom(r.b[r.i:])
		vv, mm := r.buff.get(n - m)
		v |= vv << uint(m)
		m += mm
	}
	if m < n {
		return v, m, EOF
	}
	return v, m, nil
}

// ReadAt reads n bits starting at offset off.
// The returning values are the same as Read function.
func (r *Reader) ReadAt(n, off int) (uint64, int, error) {
	if off < 0 || off > (len(r.b)*8-1) {
		return 0, 0, ErrOffsetOutOfBound
	}
	r.i = off / 8
	r.i += r.buff.loadFrom(r.b[r.i:])
	r.buff.moveTo(off % 8)
	return r.Read(n)
}

// Reset resets the Reader so that it can be read from the beginning.
func (r *Reader) Reset() {
	r.i = 0
	r.buff.reset()
}
