// Copyright (c) 2019 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

import "encoding/binary"

func read64LE(b []byte) (bits uint64) {
	n := len(b)

	switch {
	case n >= 8:
		bits = binary.LittleEndian.Uint64(b)
	case n == 7:
		bits = uint64(binary.LittleEndian.Uint32(b)) | (uint64(binary.LittleEndian.Uint16(b[4:])) << 32) | (uint64(b[6]) << 48)
	case n == 6:
		bits = uint64(binary.LittleEndian.Uint32(b)) | (uint64(binary.LittleEndian.Uint16(b[4:])) << 32)
	case n == 5:
		bits = uint64(binary.LittleEndian.Uint32(b)) | (uint64(b[4]) << 32)
	case n == 4:
		bits = uint64(binary.LittleEndian.Uint32(b))
	case n == 3:
		bits = uint64(binary.LittleEndian.Uint16(b)) | (uint64(b[2]) << 16)
	case n == 2:
		bits = uint64(binary.LittleEndian.Uint16(b))
	case n == 1:
		bits = uint64(b[0])
	case n == 0:
		break
	default:
		panic(ErrUnexpected)
	}

	return bits
}

func write64LE(bits uint64, b []byte) {
	n := len(b)

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
