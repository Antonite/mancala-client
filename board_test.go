package oware

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMove(tt *testing.T) {
	type test struct {
		name       string
		givenBoard *Board
		wantBoard  *Board
		move       int
		wantErr    bool
	}

	tests := []test{
		{
			name: "default",
			givenBoard: &Board{
				player: 0,
				scores: []int{0, 0},
				pits:   []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
			},
			wantBoard: &Board{
				player:     1,
				scores:     []int{0, 0},
				pits:       []int{0, 5, 5, 5, 5, 4, 4, 4, 4, 4, 4, 4},
				validMoves: []int{6, 7, 8, 9, 10, 11},
			},
			move:    0,
			wantErr: false,
		},
		{
			name: "simple scoring",
			givenBoard: &Board{
				player: 0,
				scores: []int{0, 0},
				pits:   []int{4, 4, 4, 4, 4, 4, 0, 1, 2, 1, 2, 1},
			},
			wantBoard: &Board{
				player:     1,
				scores:     []int{7, 0},
				pits:       []int{4, 4, 4, 4, 4, 0, 1, 0, 0, 0, 2, 1},
				validMoves: []int{6, 10, 11},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "simple scoring with skips",
			givenBoard: &Board{
				player: 0,
				scores: []int{0, 0},
				pits:   []int{4, 4, 4, 4, 4, 4, 0, 1, 12, 1, 2, 1},
			},
			wantBoard: &Board{
				player:     1,
				scores:     []int{5, 0},
				pits:       []int{4, 4, 4, 4, 4, 0, 1, 2, 12, 0, 0, 1},
				validMoves: []int{6, 7, 8, 11},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "simple scoring with two skips",
			givenBoard: &Board{
				player: 0,
				scores: []int{0, 0},
				pits:   []int{4, 4, 4, 4, 4, 4, 0, 12, 12, 1, 2, 1},
			},
			wantBoard: &Board{
				player:     1,
				scores:     []int{7, 0},
				pits:       []int{4, 4, 4, 4, 4, 0, 1, 12, 12, 0, 0, 0},
				validMoves: []int{6, 7, 8},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "no scoring due to skips",
			givenBoard: &Board{
				player: 0,
				scores: []int{0, 0},
				pits:   []int{4, 4, 4, 4, 4, 4, 0, 12, 12, 1, 2, 12},
			},
			wantBoard: &Board{
				player:     1,
				scores:     []int{0, 0},
				pits:       []int{5, 4, 4, 4, 4, 0, 1, 12, 12, 2, 3, 12},
				validMoves: []int{6, 7, 8, 9, 10, 11},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "no scoring due to opponents landing",
			givenBoard: &Board{
				player: 0,
				scores: []int{0, 0},
				pits:   []int{1, 4, 4, 4, 4, 4, 0, 12, 12, 1, 2, 12},
			},
			wantBoard: &Board{
				player:     1,
				scores:     []int{0, 0},
				pits:       []int{2, 4, 4, 4, 4, 0, 1, 12, 12, 2, 3, 12},
				validMoves: []int{6, 7, 8, 9, 10, 11},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "no moves possible for opponent",
			givenBoard: &Board{
				player: 0,
				scores: []int{10, 0},
				pits:   []int{1, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0},
			},
			wantBoard: &Board{
				player:     1,
				scores:     []int{27, 0},
				pits:       []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				validMoves: []int{},
				Status:     Player1Won,
			},
			move:    0,
			wantErr: false,
		},
		{
			name: "capture not possible since no moves left for opponent after capture",
			givenBoard: &Board{
				player: 0,
				scores: []int{0, 0},
				pits:   []int{1, 4, 4, 4, 4, 1, 2, 0, 0, 0, 0, 0},
			},
			wantBoard: &Board{
				player:     1,
				scores:     []int{0, 0},
				pits:       []int{1, 4, 4, 4, 4, 0, 3, 0, 0, 0, 0, 0},
				validMoves: []int{6},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "scoring while skipping a lot",
			givenBoard: &Board{
				player: 1,
				scores: []int{0, 0},
				pits:   []int{12, 0, 0, 0, 12, 0, 0, 12, 12, 0, 0, 0},
			},
			wantBoard: &Board{
				player:     0,
				scores:     []int{0, 2},
				pits:       []int{12, 0, 1, 1, 12, 1, 1, 12, 0, 2, 2, 2},
				validMoves: []int{0, 2, 3, 4, 5},
			},
			move:    8,
			wantErr: false,
		},
		{
			name: "complex scoring while skipping a lot",
			givenBoard: &Board{
				player: 1,
				scores: []int{0, 0},
				pits:   []int{0, 12, 0, 0, 12, 0, 0, 12, 0, 0, 12, 0},
			},
			wantBoard: &Board{
				player:     0,
				scores:     []int{0, 4},
				pits:       []int{2, 12, 0, 0, 12, 1, 1, 12, 1, 1, 0, 2},
				validMoves: []int{0, 1, 4, 5},
			},
			move:    10,
			wantErr: false,
		},
	}

	for _, t := range tests {
		tt.Run(t.name, func(tt *testing.T) {
			b, err := t.givenBoard.Move(t.move)
			if t.wantErr {
				require.Error(tt, err)
			} else {
				require.Nil(tt, err)
				require.Equal(tt, t.wantBoard, b)
			}
		})
	}
}

func TestToString(tt *testing.T) {
	type test struct {
		name       string
		givenBoard *Board
		wantString string
	}

	tests := []test{
		{
			name: "default",
			givenBoard: &Board{
				player:     0,
				scores:     []int{0, 0},
				pits:       []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
				validMoves: []int{0, 1, 2, 3, 4, 5},
				Status:     InProgress,
			},
			wantString: "0/0/4,4,4,4,4,4,4,4,4,4,4,4/0,0/0,1,2,3,4,5",
		},
		{
			name: "complex valid",
			givenBoard: &Board{
				player:     1,
				scores:     []int{25, 5},
				pits:       []int{12, 0, 0, 1, 3, 4, 0, 0, 4, 0, 4, 2},
				validMoves: []int{8, 10, 11},
				Status:     Player1Won,
			},
			wantString: "1/1/12,0,0,1,3,4,0,0,4,0,4,2/25,5/8,10,11",
		},
	}

	for _, t := range tests {
		tt.Run(t.name, func(tt *testing.T) {
			s := t.givenBoard.ToString()
			require.Equal(tt, t.wantString, s)
		})
	}
}

func TestNewS(tt *testing.T) {
	type test struct {
		name        string
		givenString string
		wantBoard   *Board
		wantErr     bool
	}

	tests := []test{
		{
			name: "default",
			wantBoard: &Board{
				player:     0,
				scores:     []int{0, 0},
				pits:       []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
				validMoves: []int{0, 1, 2, 3, 4, 5},
				Status:     InProgress,
			},
			givenString: "0/0/4,4,4,4,4,4,4,4,4,4,4,4/0,0/0,1,2,3,4,5",
			wantErr:     false,
		},
		{
			name: "complex valid",
			wantBoard: &Board{
				player:     1,
				scores:     []int{25, 5},
				pits:       []int{12, 0, 0, 1, 3, 4, 0, 0, 4, 0, 4, 2},
				validMoves: []int{8, 10, 11},
				Status:     Player1Won,
			},
			givenString: "1/1/12,0,0,1,3,4,0,0,4,0,4,2/25,5/8,10,11",
			wantErr:     false,
		},
		{
			name: "bad pits",
			wantBoard: &Board{
				player:     1,
				scores:     []int{0, 0},
				pits:       []int{12, 0, 0, 1, 3, 4, 0, 0, 4, 0, 4, 2, 3},
				validMoves: []int{8, 10, 11},
				Status:     Player1Won,
			},
			givenString: "1/1/12,0,0,3,4,0,0,4,0,4,2/25,5/8,10,11",
			wantErr:     true,
		},
		{
			name: "negative pits",
			wantBoard: &Board{
				player:     1,
				scores:     []int{0, 0},
				pits:       []int{12, 0, -1, 3, 4, 0, 0, 4, 0, 4, 2, 3},
				validMoves: []int{8, 10, 11},
				Status:     Player1Won,
			},
			givenString: "1/1/12,0,0,1,-3,4,0,0,4,0,4,2/25,5/8,10,11",
			wantErr:     true,
		},
		{
			name: "bad scores",
			wantBoard: &Board{
				player:     1,
				scores:     []int{},
				pits:       []int{12, 0, 0, 1, 3, 4, 0, 0, 4, 0, 4, 2},
				validMoves: []int{8, 10, 11},
				Status:     Player1Won,
			},
			givenString: "1/1/12,0,0,1,3,4,0,0,4,0,4,2//8,10,11",
			wantErr:     true,
		},
		{
			name: "negative scores",
			wantBoard: &Board{
				player:     1,
				scores:     []int{2, -2},
				pits:       []int{12, 0, 1, 3, 4, 0, 0, 4, 0, 4, 2, 3},
				validMoves: []int{8, 10, 11},
				Status:     Player1Won,
			},
			givenString: "1/1/12,0,0,1,3,4,0,0,4,0,4,2/25,-5/8,10,11",
			wantErr:     true,
		},
		{
			name: "bad moves",
			wantBoard: &Board{
				player:     1,
				scores:     []int{},
				pits:       []int{12, 0, 0, 1, 3, 4, 0, 0, 4, 0, 4, 2},
				validMoves: []int{8, 10, 11, 0, 2, 3, 5},
				Status:     Player1Won,
			},
			givenString: "1/1/12,0,0,1,3,4,0,0,4,0,4,2//8,10,11",
			wantErr:     true,
		},
		{
			name: "negative moves",
			wantBoard: &Board{
				player:     1,
				scores:     []int{2, -2},
				pits:       []int{12, 0, 1, 3, 4, 0, 0, 4, 0, 4, 2, 3},
				validMoves: []int{8, -10, 11},
				Status:     Player1Won,
			},
			givenString: "1/1/12,0,0,1,3,4,0,0,4,0,4,2/25,-5/8,10,11",
			wantErr:     true,
		},
		{
			name: "no valid moves",
			wantBoard: &Board{
				player:     1,
				scores:     []int{35, 12},
				pits:       []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				validMoves: []int{},
				Status:     Player1Won,
			},
			givenString: "1/1/0,0,0,0,0,0,0,0,0,0,0,0/35,12/",
			wantErr:     false,
		},
	}

	for _, t := range tests {
		tt.Run(t.name, func(tt *testing.T) {
			s, err := NewS(t.givenString)
			if t.wantErr {
				require.Error(tt, err)
			} else {
				require.Nil(tt, err)
				require.Equal(tt, t.wantBoard, s)
			}
		})
	}
}
