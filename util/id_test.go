package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashIdEncodeAndDecode(t *testing.T) {
	idStr := HashIDEncode(1)
	num, err := HashIDDecode(idStr)
	t.Logf("id: %d=>%s", num, idStr)
	assert.NotEqual(t, idStr, "")
	assert.NoError(t, err)
	assert.Equal(t, num, 1)
}

func BenchmarkHashIdEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashIDEncode(i)
	}
}

func BenchmarkHashIdDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HashIDEncode(i)
	}
}
