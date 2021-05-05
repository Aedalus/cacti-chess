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
	lineNumber int
	fen        string
	depths     [7]int
}

func (t *testCasePerft) String() string {
	return fmt.Sprintf("%v | %v\n", t.fen, t.depths)
}

func getPerftTestCases2(t *testing.T) []*testCasePerft {
	file, err := os.OpenFile("perft_test_2.txt", os.O_RDONLY, 0644)
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

func Test_Perft_All(t *testing.T) {
	tsc := getPerftTestCases2(t)

	for depth := 0; depth < 6; depth++ {
		for _, tc := range tsc {
			p, err := ParseFen(tc.fen)
			require.Nil(t, err)

			got := p.Perft(depth)
			want := tc.depths[depth]
			if want != got {
				t.Fatalf("perft error (line %d depth %d): %v | got %v, want %v\n", tc.lineNumber, depth, tc.fen, got, want)
			} else {
				fmt.Printf("%d %d | %v: got %v want %v\n", depth, tc.lineNumber, tc.fen, got, want)
			}
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
	assert.Equal(t, 4865609, p.Perft(5))
	//assert.Equal(t, 119060324, p.Perft(6))

}

func Test_Perft_Sample_A(t *testing.T) {
	fen := "r3kbnr/2qn2p1/8/pppBpp1P/3P1Pb1/P1P1P3/1P2Q2P/RNB1K1NR w KQkq - 0 1"
	p, err := ParseFen(fen)
	require.Nil(t, err)

	pft := p.Perft(2)
	want := map[string]int{
		"b2b3": 40,
		"h2h3": 40,
		"a3a4": 40,
		"c3c4": 40,
		"e3e4": 40,
		"h5h6": 40,
		"b2b4": 41,
		"h2h4": 40,
		"d4e5": 38,
		"d4c5": 40,
		"f4e5": 40,
		"b1d2": 40,
		"g1f3": 39,
		"g1h3": 40,
		"c1d2": 40,
		"d5a2": 40,
		"d5g2": 41,
		"d5b3": 40,
		"d5f3": 40,
		"d5c4": 40,
		"d5e4": 41,
		"d5c6": 38,
		"d5e6": 40,
		"d5b7": 39,
		"d5f7": 3,
		"d5a8": 35,
		"d5g8": 38,
		"a1a2": 40,
		"e2d1": 41,
		"e2f1": 41,
		"e2c2": 41,
		"e2d2": 41,
		"e2f2": 41,
		"e2g2": 41,
		"e2d3": 41,
		"e2f3": 39,
		"e2c4": 41,
		"e2g4": 37,
		"e2b5": 37,
		"e1d1": 40,
		"e1f1": 40,
		"e1d2": 40,
		"e1f2": 40,
	}
	assert.Equal(t, 1674, pft)
	assert.Equal(t, want, perftCounts)
	fmt.Println(perftCounts)
	fmt.Println(perftSubmoves)
}

func Test_Perft_Sample_B(t *testing.T) {
	// all based off of "rn1q1bnr/3kppp1/p1pp3p/1p3b2/1P6/2P2N1P/P1QPPPB1/RNB1K2R b KQ - 0 1"
	tcs := []struct {
		move     string
		fen      string
		expected int
	}{
		{
			"d4e3",
			"r2n3r/1bNk2pp/6P1/pP3p2/5qnP/1P1Ppp1R/2P3B1/Q1B1bKN1 w - - 0 2",
			36,
		},
	}

	for _, tc := range tcs {
		p, err := ParseFen(tc.fen)
		if err != nil {
			panic(err)
		}
		pft := p.Perft(1)
		if tc.expected != pft {
			panic(fmt.Sprintf("%v - got %v want %v", tc.move, pft, tc.expected))
		}
	}
}
