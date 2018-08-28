package bitstream

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaskTable(t *testing.T) {
	for i := 0; i < 65; i++ {
		expected := ^uint64(0) >> uint(64-i)
		assert.Equal(t, expected, maskTable[i])
	}
}
