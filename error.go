package bitstream

import "errors"

var (
	EOF                 = errors.New("bitstream: EOF")
	ErrSizeOutOfBound   = errors.New("bitstream: size out of bound")
	ErrOffsetOutOfBound = errors.New("bitstream: offset out of bound")
	ErrUnexpected       = errors.New("bitstream: unexpected error")
)
