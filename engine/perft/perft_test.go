package perft

import (
	"bufio"
	"cacti-chess/engine/position"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"strings"
	"testing"
)

type testCasePerft struct {
	lineNumber int
	fen        string
	depths     [7]int
}

func (t *testCasePerft) String() string {
	return fmt.Sprintf("%v | %v\n", t.fen, t.depths)
}

func getPerftTestCases(t *testing.T) []*testCasePerft {
	file, err := os.OpenFile("perft_test_cases.txt", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	testCases := []*testCasePerft{}

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber += 1
		ln := scanner.Text()

		segments := strings.Split(ln, ",")
		require.Equal(t, 7, len(segments))

		tc := &testCasePerft{}
		tc.fen = strings.TrimSpace(segments[0]) + " 0 1"
		tc.lineNumber = lineNumber

		// todo - parse depths
		for i := 0; i < 7; i++ {
			if i == 0 {
				tc.depths[i] = 1
			} else {
				sg := segments[i]
				d, err := strconv.Atoi(sg)
				if err != nil {
					panic(err)
				}
				tc.depths[i] = d
			}
		}

		testCases = append(testCases, tc)
	}

	return testCases
}

// Test_Perft_All reads in all cases from the text file, and asserts all at a depth of 2.
// It's possible to increment the depth, but starts to have noticeable runtime.
func Test_Perft_All(t *testing.T) {
	tsc := getPerftTestCases(t)

	maxDepth := 2
	for depth := 0; depth <= maxDepth; depth++ {
		for _, tc := range tsc {
			p, err := position.FromFen(tc.fen)
			require.Nil(t, err)

			got := Perft(p, depth)
			want := tc.depths[depth]
			if want != got {
				t.Fatalf("perft error (line %d depth %d): %v | got %v, want %v\n", tc.lineNumber, depth, tc.fen, got, want)
			} else {
				fmt.Printf("%d %d | %v: got %v want %v\n", depth, tc.lineNumber, tc.fen, got, want)
			}
		}
	}
}

// Test_Perft_StartingPos tests the starting position up to depth 5
func Test_Perft_StartingPos(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, err := position.FromFen(fen)
	require.Nil(t, err)

	assert.Equal(t, 1, Perft(p, 0))
	assert.Equal(t, 20, Perft(p, 1))
	assert.Equal(t, 400, Perft(p, 2))
	assert.Equal(t, 8902, Perft(p, 3))
	assert.Equal(t, 197281, Perft(p, 4))
	assert.Equal(t, 4865609, Perft(p, 5))
}
