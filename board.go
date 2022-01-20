package oware

import (
	"sort"

	"github.com/pkg/errors"
)

type GameStatus int64

const (
	InProgress GameStatus = iota
	Player1Won
	Player2Won
	Tie
)

type board struct {
	Status     GameStatus
	player     int
	scores     []int
	pits       []int
	validMoves []int
}

func New(playerToMove int, scores []int, pits []int) (*board, error) {
	b := &board{
		player: playerToMove,
		scores: scores,
		pits:   pits,
		Status: InProgress,
	}

	if err := b.validateInputs(); err != nil {
		return nil, err
	}

	return b, nil
}

func (b *board) Move(pit int) (*board, error) {
	if b.pits[pit] == 0 {
		return b, errors.New("cannot make a move on an empty pit")
	}

	b.tryMove(pit)
	b.tryEndGame()

	return b, nil
}

func (b *board) GetValidMoves() []int {
	return b.validMoves
}

func (b *board) tryEndGame() {
	b.Status = b.computeStatus()
	if b.Status != InProgress {
		b.validMoves = []int{}
		return
	}

	b.validMoves = b.computeValidMoves()
	if len(b.validMoves) == 0 {
		b.scores[0] += sum(b.pits[0:5])
		b.scores[1] += sum(b.pits[6:11])
		b.pits = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	}

	b.Status = b.computeStatus()
}

func (b *board) computeStatus() GameStatus {
	if b.scores[0] > 24 {
		return Player1Won
	}

	if b.scores[1] > 24 {
		return Player1Won
	}

	if b.scores[0] == 24 && b.scores[1] == 24 {
		return Tie
	}

	return InProgress
}

func (b *board) computeValidMoves() []int {
	validMoves := []int{}
	moves := []int{}
	valid := make(map[int]bool)
	if b.player == 0 {
		moves = []int{0, 1, 2, 3, 4, 5}
	} else {
		moves = []int{6, 7, 8, 9, 10, 11}
	}

	for _, m := range moves {
		if b.pits[m] > 0 {
			valid[m] = true
		}
	}

	for k := range valid {
		c := b.clone()
		c.tryMove(k)
		if !c.opponentCanMakeMove() {
			valid[k] = false
		} else {
			validMoves = append(validMoves, k)
		}
	}

	sort.Ints(validMoves)
	return validMoves
}

func (b *board) tryMove(pit int) {
	seeds := b.pits[pit]
	b.pits[pit] = 0
	// Play seeds
	cp := pit
	for i := 0; i < seeds; i++ {
		cp = (cp + 1) % 12
		// Skip the starting pit if we made a circle or if current pit already contains max seeds
		for cp == pit || b.pits[cp] == 12 {
			cp = (cp + 1) % 12
		}

		b.pits[cp]++
	}

	b.applyCaptures(cp)
	b.player = (b.player + 1) % 2
}

func (b *board) validateInputs() error {
	if len(b.pits) != 12 {
		return errors.New("invalid number of pits")
	}

	if len(b.scores) != 2 {
		return errors.New("invalid number of scores")
	}

	if b.player > 1 || b.player < 0 {
		return errors.New("invalid player")
	}

	if b.Status != InProgress {
		return errors.New("cannot create finished board")
	}

	return nil
}

func (b *board) applyCaptures(pit int) {
	cp := pit
	scores := append([]int{}, b.scores...)
	pits := append([]int{}, b.pits...)
	for cp > 0 && isOpponentPit(cp, b.player) && (b.pits[cp] == 2 || b.pits[cp] == 3) {
		b.scores[b.player] += b.pits[cp]
		b.pits[cp] = 0
		cp--
	}

	if !b.opponentCanMakeMove() {
		// restore scores and pits before captures
		b.scores = scores
		b.pits = pits
	}
}

func (b *board) opponentCanMakeMove() bool {
	return (b.player == 0 && sum(b.pits[6:11]) > 0) || (b.player == 1 && sum(b.pits[0:5]) > 0)
}

func (b *board) clone() *board {
	return &board{
		Status:     b.Status,
		player:     b.player,
		scores:     append([]int{}, b.scores...),
		pits:       append([]int{}, b.pits...),
		validMoves: append([]int{}, b.validMoves...),
	}
}

func isOpponentPit(pit int, player int) bool {
	return pit > 5 && player == 0 || pit <= 5 && player == 1
}

func sum(s []int) int {
	total := 0
	for _, i := range s {
		total += i
	}

	return total
}
