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
	depths     [6]int
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
		for i := 0; i < 6; i++ {
			sg := segments[i+1]
			d, err := strconv.Atoi(sg)
			if err != nil {
				panic(err)
			}
			tc.depths[i] = d
		}

		testCases = append(testCases, tc)
	}

	return testCases
}

func Test_Perft_All(t *testing.T) {
	depth := 2
	tsc := getPerftTestCases2(t)

	for _, tc := range tsc {
		p, err := ParseFen(tc.fen)
		require.Nil(t, err)

		got := p.Perft(depth)
		want := tc.depths[depth-1]
		if want != got {
			t.Fatalf("perft error (line %d): %v | got %v, want %v\n", tc.lineNumber, tc.fen, got, want)
		} else {
			fmt.Printf("%d | %v: got %v want %v\n", tc.lineNumber, tc.fen, got, want)
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
}

func Test_Perft_Sample_A(t *testing.T) {
	fen := "2b1kbnB/rppqp3/3p3p/3P1pp1/pnP3P1/PP2P2P/4QP2/RN2KBNR b KQ - 0 1"
	p, err := ParseFen(fen)
	require.Nil(t, err)

	// movegen isn't picking up enPas after c7 c5
	want := map[string]int{
		"f5f4": 28,
		"h6h5": 30,
		"b7b6": 29,
		"c7c6": 30,
		"e7e6": 30,
		"b7b5": 30,
		"c7c5": 29,
		"e7e5": 27,
		"a4b3": 29,
		"f5g4": 30,
		"b4a2": 29,
		"b4c2": 3,
		"b4d3": 3,
		"b4d5": 30,
		"b4a6": 29,
		"b4c6": 30,
		"g8f6": 25,
		"f8g7": 24,
		"a7a5": 29,
		"a7a6": 29,
		"a7a8": 29,
		"d7b5": 30,
		"d7c6": 30,
		"d7e6": 30,
		"d7d8": 29,
		"e8f7": 29,
		"e8d8": 29,
	}

	p.Perft(2)
	assert.Equal(t, want, perftCounts)
	fmt.Println(perftSubmoves)
	//assert.Equal(t, 729, p.Perft(2))
}

func Test_Perft_Sample_B(t *testing.T) {
	// all based off of "rn1q1bnr/3kppp1/p1pp3p/1p3b2/1P6/2P2N1P/P1QPPPB1/RNB1K2R b KQ - 0 1"
	tcs := []struct {
		move     string
		fen      string
		expected int
	}{
		{
			"c7c5",
			"2b1kbnB/rp1qp3/3p3p/2pP1pp1/pnP3P1/PP2P2P/4QP2/RN2KBNR w KQ c6 0 2",
			29,
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
