package engine

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"strings"
	"testing"
)

type testCasePerft struct {
	fen    string
	depths [6]int
}

func (t *testCasePerft) String() string {
	return fmt.Sprintf("%v | %v\n", t.fen, t.depths)
}

func getPerftTestCases(t *testing.T) []*testCasePerft {
	file, err := os.OpenFile("perft_test.txt", os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	testCases := []*testCasePerft{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ln := scanner.Text()

		segments := strings.Split(ln, ";")
		require.Equal(t, 7, len(segments))

		tc := &testCasePerft{}
		tc.fen = strings.TrimSpace(segments[0])

		// todo - parse depths
		for i := 0; i < 6; i++ {
			sg := segments[i+1]
			str := strings.ReplaceAll(sg, fmt.Sprintf("D%d", i+1), "")
			str = strings.TrimSpace(str)
			d, err := strconv.Atoi(str)
			if err != nil {
				panic(err)
			}
			tc.depths[i] = d
		}

		testCases = append(testCases, tc)
	}

	return testCases
}

func Test_Perft_Deserialization(t *testing.T) {
	tcs := getPerftTestCases(t)

	assert.Equal(t, 126, len(tcs))

	for _, tc := range tcs {
		for i := 0; i < 6; i++ {
			assert.NotEqual(t, 0, tc.depths[i])
		}
	}
}

func Test_Perft_All(t *testing.T) {
	depth := 1
	tsc := getPerftTestCases(t)

	for _, tc := range tsc {
		p, err := ParseFen(tc.fen)
		require.Nil(t, err)

		got := p.Perft(depth)
		want := tc.depths[depth-1]
		if want != got {
			t.Fatalf("perft error: %v | got %v, want %v\n", tc.fen, got, want)
		} else {
			fmt.Printf("%v: got %v want %v\n", tc.fen, got, want)
		}
	}
}
func Test_Perft_StartingPos(t *testing.T) {
	fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p, err := ParseFen(fen)
	require.Nil(t, err)

	assert.Equal(t, 1, p.Perft(0))
	assert.Equal(t, 20, p.Perft(1))
	assert.Equal(t, 400, p.Perft(2))
	assert.Equal(t, 8902, p.Perft(3))
	assert.Equal(t, 197281, p.Perft(4))
	//assert.Equal(t, 4865609, p.Perft(5))
}
