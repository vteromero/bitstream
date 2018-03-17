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

	rem := n
	for rem > 0 {
		i := n - rem
		shift := uint(i * 8)

		switch {
		case rem == 8:
			buff.bits |= binary.LittleEndian.Uint64(b[i:]) << shift
			rem -= 8
		case rem >= 4:
			buff.bits |= uint64(binary.LittleEndian.Uint32(b[i:])) << shift
			rem -= 4
		case rem >= 2:
			buff.bits |= uint64(binary.LittleEndian.Uint16(b[i:])) << shift
			rem -= 2
		case rem == 1:
			buff.bits |= uint64(b[i]) << shift
			rem -= 1
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

	rem := n
	for rem > 0 {
		i := n - rem
		v := buff.bits >> uint(i*8)

		switch {
		case rem == 8:
			binary.LittleEndian.PutUint64(b[i:], v)
			rem -= 8
		case rem >= 4:
			binary.LittleEndian.PutUint32(b[i:], uint32(v))
			rem -= 4
		case rem >= 2:
			binary.LittleEndian.PutUint16(b[i:], uint16(v))
			rem -= 2
		case rem == 1:
			b[i] = byte(v)
			rem -= 1
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
	var mask uint64 = ^uint64(0) >> uint(64-m)
	var x uint64 = (buff.bits >> uint(buff.off)) & mask
	buff.off += m
	return x, m
}

func (buff *bitBuffer) set(v uint64, n int) int {
	m := buff.len - buff.off
	if n < m {
		m = n
	}
	var mask uint64 = ^uint64(0) >> uint(64-m)
	buff.bits |= (v & mask) << uint(buff.off)
	buff.off += m
	return m
}
