package bitstream

import (
	"math/rand"
	"testing"
)

func randomByteSlice(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte(rand.Intn(256))
	}
	return b
}

func randomIntSlice(n, min, max int) []int {
	s := make([]int, n)
	for i := 0; i < n; i++ {
		s[i] = min + rand.Intn(max-min+1)
	}
	return s
}

func randomUint64Slice(n int) []uint64 {
	s := make([]uint64, n)
	for i := 0; i < n; i++ {
		s[i] = rand.Uint64()
	}
	return s
}

func BenchmarkRead(b *testing.B) {
	benchmarks := []struct {
		name        string
		minReadSize int
		maxReadSize int
	}{
		{"SmallSizes", 1, 16},
		{"MediumSizes", 17, 32},
		{"LargeSizes", 33, 48},
		{"ExtraLargeSizes", 49, 64},
		{"AllSize", 1, 64},
	}
	n := 10000000

	for _, bm := range benchmarks {
		data := randomByteSlice(n * 8)
		sizes := randomIntSlice(n, bm.minReadSize, bm.maxReadSize)
		r := NewReader(data)

		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if i == n {
					r.Reset()
				}
				r.Read(sizes[i%n])
			}
		})
	}
}

func BenchmarkWrite(b *testing.B) {
	benchmarks := []struct {
		name        string
		minReadSize int
		maxReadSize int
	}{
		{"SmallSizes", 1, 16},
		{"MediumSizes", 17, 32},
		{"LargeSizes", 33, 48},
		{"ExtraLargeSizes", 49, 64},
		{"AllSize", 1, 64},
	}
	n := 10000000

	for _, bm := range benchmarks {
		data := make([]byte, n*8)
		sizes := randomIntSlice(n, bm.minReadSize, bm.maxReadSize)
		values := randomUint64Slice(n)
		w := NewWriter(data)

		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if i == n {
					w = NewWriter(data)
				}
				w.Write(values[i%n], sizes[i%n])
			}
		})
	}
}
