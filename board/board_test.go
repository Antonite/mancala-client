package board

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
				Player: 0,
				Scores: []int{0, 0},
				Pits:   []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{0, 0},
				Pits:   []int{0, 5, 5, 5, 5, 4, 4, 4, 4, 4, 4, 4},
			},
			move:    0,
			wantErr: false,
		},
		{
			name: "bad input pits",
			givenBoard: &Board{
				Player: 0,
				Scores: []int{0, 0},
				Pits:   []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{0, 0},
				Pits:   []int{0, 5, 5, 5, 5, 4, 4, 4, 4, 4, 4, 4},
			},
			move:    0,
			wantErr: true,
		},
		{
			name: "bad input scores",
			givenBoard: &Board{
				Player: 0,
				Scores: []int{7},
				Pits:   []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{7},
				Pits:   []int{0, 5, 5, 5, 5, 4, 4, 4, 4, 4, 4, 4},
			},
			move:    0,
			wantErr: true,
		},
		{
			name: "bad input players",
			givenBoard: &Board{
				Player: 2,
				Scores: []int{0, 0},
				Pits:   []int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{0, 0},
				Pits:   []int{0, 5, 5, 5, 5, 4, 4, 4, 4, 4, 4, 4},
			},
			move:    0,
			wantErr: true,
		},
		{
			name: "simple scoring",
			givenBoard: &Board{
				Player: 0,
				Scores: []int{0, 0},
				Pits:   []int{4, 4, 4, 4, 4, 4, 0, 1, 2, 1, 2, 1},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{7, 0},
				Pits:   []int{4, 4, 4, 4, 4, 0, 1, 0, 0, 0, 2, 1},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "simple scoring with skips",
			givenBoard: &Board{
				Player: 0,
				Scores: []int{0, 0},
				Pits:   []int{4, 4, 4, 4, 4, 4, 0, 1, 12, 1, 2, 1},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{5, 0},
				Pits:   []int{4, 4, 4, 4, 4, 0, 1, 2, 12, 0, 0, 1},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "simple scoring with two skips",
			givenBoard: &Board{
				Player: 0,
				Scores: []int{0, 0},
				Pits:   []int{4, 4, 4, 4, 4, 4, 0, 12, 12, 1, 2, 1},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{7, 0},
				Pits:   []int{4, 4, 4, 4, 4, 0, 1, 12, 12, 0, 0, 0},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "no scoring due to skips",
			givenBoard: &Board{
				Player: 0,
				Scores: []int{0, 0},
				Pits:   []int{4, 4, 4, 4, 4, 4, 0, 12, 12, 1, 2, 12},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{0, 0},
				Pits:   []int{5, 4, 4, 4, 4, 0, 1, 12, 12, 2, 3, 12},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "no scoring due to opponents landing",
			givenBoard: &Board{
				Player: 0,
				Scores: []int{0, 0},
				Pits:   []int{1, 4, 4, 4, 4, 4, 0, 12, 12, 1, 2, 12},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{0, 0},
				Pits:   []int{2, 4, 4, 4, 4, 0, 1, 12, 12, 2, 3, 12},
			},
			move:    5,
			wantErr: false,
		},
		{
			name: "no moves possible for opponent",
			givenBoard: &Board{
				Player: 0,
				Scores: []int{0, 0},
				Pits:   []int{1, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{0, 0},
				Pits:   []int{2, 4, 4, 4, 4, 0, 1, 12, 12, 2, 3, 12},
			},
			move:    0,
			wantErr: true,
		},
		{
			name: "no moves possible for opponent due to captures",
			givenBoard: &Board{
				Player: 0,
				Scores: []int{0, 0},
				Pits:   []int{1, 4, 4, 4, 4, 1, 2, 0, 0, 0, 0, 0},
			},
			wantBoard: &Board{
				Player: 1,
				Scores: []int{0, 0},
				Pits:   []int{2, 4, 4, 4, 4, 0, 1, 12, 12, 2, 3, 12},
			},
			move:    5,
			wantErr: true,
		},
		{
			name: "scoring while skipping a lot",
			givenBoard: &Board{
				Player: 1,
				Scores: []int{0, 0},
				Pits:   []int{12, 0, 0, 0, 12, 0, 0, 12, 12, 0, 0, 0},
			},
			wantBoard: &Board{
				Player: 0,
				Scores: []int{0, 2},
				Pits:   []int{12, 0, 1, 1, 12, 1, 1, 12, 0, 2, 2, 2},
			},
			move:    8,
			wantErr: false,
		},
		{
			name: "complex scoring while skipping a lot",
			givenBoard: &Board{
				Player: 1,
				Scores: []int{0, 0},
				Pits:   []int{0, 12, 0, 0, 12, 0, 0, 12, 0, 0, 12, 0},
			},
			wantBoard: &Board{
				Player: 0,
				Scores: []int{0, 4},
				Pits:   []int{2, 12, 0, 0, 12, 1, 1, 12, 1, 1, 0, 2},
			},
			move:    10,
			wantErr: false,
		},
	}

	for _, t := range tests {
		tt.Run(t.name, func(tt *testing.T) {
			b := t.givenBoard
			err := b.Move(t.move)
			if t.wantErr {
				require.Error(tt, err)
			} else {
				require.Nil(tt, err)
				require.Equal(tt, t.wantBoard, b)
			}
		})
	}
}
