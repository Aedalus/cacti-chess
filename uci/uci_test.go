package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func Test_parseGoLine(t *testing.T) {
	t.Run("It can set defaults", func(t *testing.T) {
		defaults := GoCmdArgs{
			SearchMoves: []string{},
			Ponder:      false,
			Wtime:       0,
			Btime:       0,
			Winc:        0,
			Binc:        0,
			MovesToGo:   0,
			Depth:       3,
			Nodes:       0,
			Mate:        0,
			MoveTime:    0,
			Infinite:    false,
		}

		assert.Equal(t, defaults, parseGoCmdArgs([]string{}))
	})

	t.Run("searchmoves", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.SearchMoves = []string{"a2a4", "h7h8q"}

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go infinite searchmoves a2a4 h7h8q", " ")))
	})

	t.Run("ponder", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.Ponder = true

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go ponder", " ")))
	})

	t.Run("wtime/btime", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.Wtime = time.Millisecond * 500
		want.Btime = time.Millisecond * 1000

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go infinite wtime 500 btime 1000", " ")))
	})

	t.Run("winc/binc", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.Winc = time.Millisecond * 500
		want.Binc = time.Millisecond * 1000

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go infinite winc 500 binc 1000", " ")))
	})

	t.Run("movestogo", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.MovesToGo = 23

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go infinite movestogo 23", " ")))
	})

	t.Run("depth", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.Depth = 5

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go depth 5", " ")))
	})

	t.Run("nodes", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.Nodes = 1000

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go nodes 1000", " ")))
	})

	t.Run("mate", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.Mate = 3

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go mate 3", " ")))
	})

	t.Run("movetime", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.MoveTime = time.Millisecond * 500

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go movetime 500", " ")))
	})

	t.Run("infinite", func(t *testing.T) {
		want := parseGoCmdArgs([]string{})
		want.Infinite = true

		assert.Equal(t, want, parseGoCmdArgs(strings.Split("go infinite", " ")))
	})
}
