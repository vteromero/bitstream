// Copyright (c) 2018 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

// Writer is a structure to write bits on a stream of bytes.
type Writer struct {
	b   []byte // byte stream
	off int    // writing position
}

// NewWriter returns a Writer that can write on the stream of bytes provided.
func NewWriter(b []byte) *Writer {
	return &Writer{
		b:   b,
		off: 0,
	}
}

// Write appends to the byte stream the least-significant n bits of the value v.
// It returns an error value indicating if something went wrong.
// When the end of the stream is reached, it returns an error EOF.
func (w *Writer) Write(v uint64, n int) error {
	if n < 0 || n > 64-7 {
		return ErrSizeOutOfBound
	}
	if w.off+n > len(w.b)<<3 {
		return EOF
	}
	i := w.off >> 3
	bits := uint64(w.b[i]) | ((v & maskTable[n]) << uint(w.off&7))
	write64LE(bits, w.b[i:])
	w.off += n
	return nil
}

// Offset returns the current writing position.
// It also indicates the number of bits already written by using Write function
// exclusively.
func (w *Writer) Offset() int {
	return w.off
}

// Reset resets the Writer so that it can be written from the beginning.
func (w *Writer) Reset() {
	w.off = 0
}
