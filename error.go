// Copyright (c) 2018 Vicente Romero. All rights reserved.
// Licensed under the MIT License.
// See LICENSE file in the project root for full license information.

package bitstream

import "errors"

var (
	EOF                 = errors.New("bitstream: EOF")
	ErrSizeOutOfBound   = errors.New("bitstream: size out of bound")
	ErrOffsetOutOfBound = errors.New("bitstream: offset out of bound")
	ErrUnexpected       = errors.New("bitstream: unexpected error")
)
