package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileRankLookups(t *testing.T) {
	// just spot check a couple rank/files

	assert.Equal(t, FILE_A, fileLookups[A1])
	assert.Equal(t, FILE_B, fileLookups[B1])
	assert.Equal(t, FILE_H, fileLookups[H4])
	assert.Equal(t, FILE_C, fileLookups[C8])

	assert.Equal(t, RANK_1, rankLookups[A1])
	assert.Equal(t, RANK_3, rankLookups[A3])
	assert.Equal(t, RANK_8, rankLookups[F8])
}
