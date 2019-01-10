// Copyright (c) 2018 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

import "encoding/binary"

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

func (w *Writer) write64LE(bits uint64) {
	i := w.off >> 3
	n := len(w.b) - i
	b := w.b[i:]

	switch {
	case n >= 8:
		binary.LittleEndian.PutUint64(b, bits)
	case n == 7:
		binary.LittleEndian.PutUint32(b, uint32(bits))
		binary.LittleEndian.PutUint16(b[4:], uint16(bits>>32))
		b[6] = byte(bits >> 48)
	case n == 6:
		binary.LittleEndian.PutUint32(b, uint32(bits))
		binary.LittleEndian.PutUint16(b[4:], uint16(bits>>32))
	case n == 5:
		binary.LittleEndian.PutUint32(b, uint32(bits))
		b[4] = byte(bits >> 32)
	case n == 4:
		binary.LittleEndian.PutUint32(b, uint32(bits))
	case n == 3:
		binary.LittleEndian.PutUint16(b, uint16(bits))
		b[2] = byte(bits >> 16)
	case n == 2:
		binary.LittleEndian.PutUint16(b, uint16(bits))
	case n == 1:
		b[0] = byte(bits)
	case n == 0:
		break
	default:
		panic(ErrUnexpected)
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
	bits := uint64(w.b[w.off>>3]) | ((v & maskTable[n]) << uint(w.off&7))
	w.write64LE(bits)
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
