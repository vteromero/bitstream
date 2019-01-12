# bitstream

[![Build Status](https://travis-ci.org/vteromero/bitstream.svg?branch=master)](https://travis-ci.org/vteromero/bitstream)
[![GoDoc](https://godoc.org/github.com/vteromero/bitstream?status.svg)](https://godoc.org/github.com/vteromero/bitstream)
[![Go Report Card](https://goreportcard.com/badge/github.com/vteromero/bitstream)](https://goreportcard.com/report/github.com/vteromero/bitstream)
[![Coverage Status](https://coveralls.io/repos/github/vteromero/bitstream/badge.svg?branch=master)](https://coveralls.io/github/vteromero/bitstream?branch=master)

`bitstream` is a Go library to read and write bit-length values on a stream of bytes.

### Installation

```
$ go get -v -t github.com/vteromero/bitstream
```

### Usage

The [API documentation](https://godoc.org/github.com/vteromero/bitstream) is available on godoc.org.

### Examples

Here is an example of the usage of `bitstream.Reader`:

```go
package main

import (
	"fmt"

	"github.com/vteromero/bitstream"
)

func showValues(bits uint64, err error) {
	fmt.Printf("value: %b, error: %v\n", bits, err)
}

func main() {
	data := []byte{0x55, 0x55} // in binary: 01010101 01010101

	// create a new Reader
	r := bitstream.NewReader(data)

	// read 16 bits
	bits, err := r.Read(16)
	showValues(bits, err)

	// reset the Reader to start over
	r.Reset()

	// read bit to bit until EOF
	for {
		bits, err := r.Read(1)
		if err != nil {
			fmt.Println(err)
			break
		}
		showValues(bits, err)
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
	w.Write(0x77665544332211, 56)
	w.Write(0x2, 2)
	w.Write(0x2, 2)
	w.Write(0x2, 2)
	w.Write(0x2, 2)
	w.Write(0xb, 4)
	w.Write(0xb, 4)
	w.Write(0xcc, 8)

	fmt.Printf("% x\n", data)
}
```

### Benchmarks

You can test the performance of the library in this way:

```
$ cd $GOPATH/src/github.com/vteromero/bitstream
$ go test -run=^$ -bench=.
```

As a reference, here is the outcome on a laptop Ubuntu Desktop 18.11 with a Core i7-6700HQ CPU @ 2.60GHz x 8

```
BenchmarkRead/SmallSizes-8         	100000000	        16.3 ns/op
BenchmarkRead/MediumSizes-8        	100000000	        16.3 ns/op
BenchmarkRead/LargeSizes-8         	100000000	        16.3 ns/op
BenchmarkRead/ExtraLargeSizes-8    	100000000	        16.3 ns/op
BenchmarkRead/AllSizes-8           	100000000	        16.3 ns/op
BenchmarkWrite/SmallSizes-8        	100000000	        17.0 ns/op
BenchmarkWrite/MediumSizes-8       	100000000	        17.0 ns/op
BenchmarkWrite/LargeSizes-8        	100000000	        17.0 ns/op
BenchmarkWrite/ExtraLargeSizes-8   	100000000	        17.0 ns/op
BenchmarkWrite/AllSizes-8          	100000000	        17.0 ns/op
```
