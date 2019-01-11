// Copyright (c) 2018 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

// Reader is a structure to read bits from a stream of bytes.
type Reader struct {
	b   []byte // byte stream
	off int    // reading position
}

// NewReader returns a Reader that can read bits from the byte slice provided.
func NewReader(b []byte) *Reader {
	return &Reader{
		b:   b,
		off: 0,
	}
}

// Read reads the next n bits from the stream.
// It returns a 64-bit integer which holds the bits, and an error value.
// The returning error is not nil when something went wrong or when the end of
// the stream has been reached (EOF error).
func (r *Reader) Read(n int) (bits uint64, err error) {
	if n < 0 || n > 64-7 {
		return bits, ErrSizeOutOfBound
	}
	if r.off+n > len(r.b)<<3 {
		return bits, EOF
	}
	bits = (read64LE(r.b[r.off>>3:]) >> uint(r.off&7)) & maskTable[n]
	r.off += n
	return bits, err
}

// ReadAt reads n bits starting at position off.
// The returning values are the same as Read function.
func (r *Reader) ReadAt(n, off int) (uint64, error) {
	if off < 0 || off > len(r.b)<<3 {
		return 0, ErrOffsetOutOfBound
	}
	r.off = off
	return r.Read(n)
}

// Reset resets the Reader so that it can be read from the beginning.
func (r *Reader) Reset() {
	r.off = 0
}

// Offset returns the current reading position.
// It also indicates the number of bits already read by using Read function
// exclusively.
func (r *Reader) Offset() int {
	return r.off
}
