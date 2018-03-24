// Copyright (c) 2018 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

// A Writer is a structure to write bits on a stream of bytes.
type Writer struct {
	b    []byte     // byte stream
	i    int        // byte stream index at which buff writes
	off  int        // bit-length offset
	buff *bitBuffer // bit buffer
}

// NewWriter returns a Writer that can write on the stream of bytes provided.
func NewWriter(b []byte) *Writer {
	return &Writer{
		b:    b,
		i:    0,
		off:  0,
		buff: newBitBuffer(0),
	}
}

func (w *Writer) flush() {
	w.i += w.buff.writeTo(w.b[w.i:])
	w.off += w.buff.off
}

// Write appends n bits of the value v on the stream.
// It returns the actual number the bits written and an error value indicating
// if something went wrong.
// The returning size of bits written is always less or equal to n.
// When the end of the stream is reached, it returns an error EOF.
func (w *Writer) Write(v uint64, n int) (int, error) {
	if n < 0 || n > 64 {
		return 0, ErrSizeOutOfBound
	}
	m := w.buff.set(v, n)
	if m < n {
		w.flush()
		sz := (len(w.b) - w.i) * 8
		if sz > 64 {
			sz = 64
		}
		w.buff = newBitBuffer(sz)
		m += w.buff.set(v>>uint(m), n-m)
	}
	if m < n {
		return m, EOF
	}
	return m, nil
}

// Close ends the writing process.
// It is very important to always close the Writer once all the writes have
// been done. The Writer uses a bit buffer and some writes (the last ones)
// might not have actually been written on the stream of bytes.
func (w *Writer) Close() {
	w.flush()
	w.buff = newBitBuffer(0)
}

// Offset returns the current writing position.
// It also indicates the number of bits already written by using Write function
// exclusively.
func (w *Writer) Offset() int {
	return w.off
}
