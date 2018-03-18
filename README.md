# bitstream

[![Build Status](https://travis-ci.org/vteromero/bitstream.svg?branch=master)](https://travis-ci.org/vteromero/bitstream)
[![Go Report Card](https://goreportcard.com/badge/github.com/vteromero/bitstream)](https://goreportcard.com/report/github.com/vteromero/bitstream)
[![GoDoc](https://godoc.org/github.com/vteromero/bitstream?status.svg)](https://godoc.org/github.com/vteromero/bitstream)

`bitstream` is a Go library to read and write bit-length values on a stream of bytes. It has been designed and optimized to be fast and this is the main goal of this library.

## Installation

To install `bitstream`, simply run:

```
$ go get -v -t github.com/vteromero/bitstream
```

Once the `get` completes, you should find the `bitstream` executable inside `$GOPATH/bin`.

## Examples

Here is an example of `bitstream.Reader`:

```go
package main

import (
	"fmt"

	"github.com/vteromero/bitstream"
)

func showValues(v uint64, m int, err error) {
	fmt.Printf("value: %b, size: %d, error: %v\n", v, m, err)
}

func main() {
	data := []byte{0x55, 0x55} // in binary: 01010101 01010101

	// create a new Reader
	r := bitstream.NewReader(data)

	// read 16 bits
	v, m, err := r.Read(16)
	showValues(v, m, err)

	// reset the Reader to start over
	r.Reset()

	// read bit to bit until EOF
	for {
		v, m, err := r.Read(1)
		if err != nil {
			fmt.Println(err)
			break
		}
		showValues(v, m, err)
	}
}
```

And the following shows how to use `bitstream.Writer`:

```go
package main

import (
	"fmt"

	"github.com/vteromero/bitstream"
)

func main() {
	// create an empty byte slice
	data := make([]byte, 10)

	// create a new Writer
	w := bitstream.NewWriter(data)

	// some writings
	w.Write(0x8877665544332211, 64)
	w.Write(0x2, 2)
	w.Write(0x2, 2)
	w.Write(0x2, 2)
	w.Write(0x2, 2)
	w.Write(0xb, 4)
	w.Write(0xb, 4)

	// always close the writer
	w.Close()

	fmt.Printf("% x\n", data)
}
```
