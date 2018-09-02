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

	switch n {
	case 8:
		buff.bits |= binary.LittleEndian.Uint64(b)
	case 7:
		buff.bits |= uint64(binary.LittleEndian.Uint32(b)) | (uint64(binary.LittleEndian.Uint16(b[4:])) << 32) | (uint64(b[6]) << 48)
	case 6:
		buff.bits |= uint64(binary.LittleEndian.Uint32(b)) | (uint64(binary.LittleEndian.Uint16(b[4:])) << 32)
	case 5:
		buff.bits |= uint64(binary.LittleEndian.Uint32(b)) | (uint64(b[4]) << 32)
	case 4:
		buff.bits |= uint64(binary.LittleEndian.Uint32(b))
	case 3:
		buff.bits |= uint64(binary.LittleEndian.Uint16(b)) | (uint64(b[2]) << 16)
	case 2:
		buff.bits |= uint64(binary.LittleEndian.Uint16(b))
	case 1:
		buff.bits |= uint64(b[0])
	case 0:
	default:
		panic(ErrUnexpected)
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

	switch n {
	case 8:
		binary.LittleEndian.PutUint64(b, buff.bits)
	case 7:
		binary.LittleEndian.PutUint32(b, uint32(buff.bits))
		binary.LittleEndian.PutUint16(b[4:], uint16(buff.bits>>32))
		b[6] = byte(buff.bits >> 48)
	case 6:
		binary.LittleEndian.PutUint32(b, uint32(buff.bits))
		binary.LittleEndian.PutUint16(b[4:], uint16(buff.bits>>32))
	case 5:
		binary.LittleEndian.PutUint32(b, uint32(buff.bits))
		b[4] = byte(buff.bits >> 32)
	case 4:
		binary.LittleEndian.PutUint32(b, uint32(buff.bits))
	case 3:
		binary.LittleEndian.PutUint16(b, uint16(buff.bits))
		b[2] = byte(buff.bits >> 16)
	case 2:
		binary.LittleEndian.PutUint16(b, uint16(buff.bits))
	case 1:
		b[0] = byte(buff.bits)
	case 0:
	default:
		panic(ErrUnexpected)
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
