// Copyright (c) 2018 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

import (
	"encoding/binary"
)

type bitBuffer struct {
	bits uint64
	len  int
	off  int
}

func newBitBuffer(size int) *bitBuffer {
	return &bitBuffer{len: size}
}

func (buff *bitBuffer) reset() {
	buff.bits = 0
	buff.len = 0
	buff.off = 0
}

func (buff *bitBuffer) moveTo(i int) {
	buff.off = i
}

func (buff *bitBuffer) loadFrom(b []byte) int {
	buff.reset()

	n := len(b)
	if n > 8 {
		n = 8
	}

	left := n
	for left > 0 {
		i := n - left
		shift := uint(i * 8)

		switch {
		case left == 8:
			buff.bits |= binary.LittleEndian.Uint64(b[i:]) << shift
			left -= 8
		case left >= 4:
			buff.bits |= uint64(binary.LittleEndian.Uint32(b[i:])) << shift
			left -= 4
		case left >= 2:
			buff.bits |= uint64(binary.LittleEndian.Uint16(b[i:])) << shift
			left -= 2
		case left == 1:
			buff.bits |= uint64(b[i]) << shift
			left--
		default:
			panic(ErrUnexpected)
		}
	}

	buff.len = n * 8

	return n
}

func (buff *bitBuffer) writeTo(b []byte) int {
	n := 0
	if buff.len > 0 {
		n = ((buff.len - 1) / 8) + 1
	}
	if len(b) < n {
		n = len(b)
	}

	left := n
	for left > 0 {
		i := n - left
		v := buff.bits >> uint(i*8)

		switch {
		case left == 8:
			binary.LittleEndian.PutUint64(b[i:], v)
			left -= 8
		case left >= 4:
			binary.LittleEndian.PutUint32(b[i:], uint32(v))
			left -= 4
		case left >= 2:
			binary.LittleEndian.PutUint16(b[i:], uint16(v))
			left -= 2
		case left == 1:
			b[i] = byte(v)
			left--
		default:
			panic(ErrUnexpected)
		}
	}

	return n
}

func (buff *bitBuffer) get(n int) (uint64, int) {
	m := buff.len - buff.off
	if n < m {
		m = n
	}
	x := (buff.bits >> uint(buff.off)) & maskTable[m]
	buff.off += m
	return x, m
}

func (buff *bitBuffer) set(v uint64, n int) int {
	m := buff.len - buff.off
	if n < m {
		m = n
	}
	buff.bits |= (v & maskTable[m]) << uint(buff.off)
	buff.off += m
	return m
}
