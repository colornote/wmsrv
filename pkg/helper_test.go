package pkg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test GeneratePasswordHash
func TestGeneratePasswordHash(t *testing.T) {
	password := "test"
	hash, err := GeneratePasswordHash(password)
	assert.Nil(t, err)
	assert.NotEqual(t, password, hash)
}

// Test ComparePasswordHash
func TestComparePasswordHash(t *testing.T) {
	password := "test"
	hash, err := GeneratePasswordHash(password)
	assert.Nil(t, err)
	assert.NotEqual(t, password, hash)
	err = ComparePasswordHash(hash, password)
	assert.Nil(t, err)
}

// Test WordSegment
func TestWordSegment(t *testing.T) {
	os.Setenv("DICT_PATH", "/Users/color/Downloads/dictionary.txt")
	SegmentInit()

	text := []byte("己的手臂肌肉强壮")
	tt := WordSegment(text)
	assert.Equal(t, "中华人民共和国 中央人民政府", tt)
}
