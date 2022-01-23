package oware

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type GameStatus int64

const (
	InProgress GameStatus = iota
	Player1Won
	Player2Won
	Tie
)

type Board struct {
	Status     GameStatus
	player     int
	scores     []int
	pits       []int
	validMoves []int
}

func Initialize() *Board {
	b, _ := New(
		0,
		[]int{0, 0},
		[]int{4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4},
		[]int{0, 1, 2, 3, 4, 5},
		InProgress,
	)

	return b
}

func New(playerToMove int, scores []int, pits []int, validMoves []int, status GameStatus) (*Board, error) {
	b := &Board{
		player:     playerToMove,
		scores:     scores,
		pits:       pits,
		Status:     status,
		validMoves: validMoves,
	}

	if err := b.validateInputs(); err != nil {
		return nil, err
	}

	return b, nil
}

// "Status/Player/pit0,pit1,...pit11/score1,score2/validmove1,..."
func (b *Board) ToString() string {
	s := fmt.Sprintf("%v/%v/%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v/%v,%v/",
		b.Status, b.player,
		b.pits[0], b.pits[1], b.pits[2], b.pits[3], b.pits[4], b.pits[5],
		b.pits[6], b.pits[7], b.pits[8], b.pits[9], b.pits[10], b.pits[11],
		b.scores[0], b.scores[1])
	for _, v := range b.validMoves {
		s += strconv.Itoa(v) + ","
	}
	s = strings.TrimSuffix(s, ",")
	return s
}

func NewS(s string) (*Board, error) {
	v := strings.Split(s, "/")
	if len(v) != 5 {
		return nil, errors.New("invalid number of variables")
	}

	status, err := strconv.Atoi(v[0])
	if err != nil {
		return nil, errors.New("invalid status")
	}

	player, err := strconv.Atoi(v[1])
	if err != nil {
		return nil, errors.New("invalid player")
	}

	// Pits
	p := strings.Split(v[2], ",")
	pits := []int{}
	for _, sp := range p {
		pit, err := strconv.Atoi(sp)
		if err != nil {
			return nil, errors.Wrap(err, "invalid pit")
		} else if pit < 0 {
			return nil, errors.New("invalid pit")
		}
		pits = append(pits, pit)
	}

	// Scores
	sc := strings.Split(v[3], ",")
	scores := []int{}
	for _, scr := range sc {
		score, err := strconv.Atoi(scr)
		if err != nil {
			return nil, errors.Wrap(err, "invalid score")
		} else if score < 0 {
			return nil, errors.New("invalid score")
		}
		scores = append(scores, score)
	}

	// Valid moves
	moves := []int{}
	if v[4] != "" {
		vm := strings.Split(v[4], ",")
		for _, vmo := range vm {
			move, err := strconv.Atoi(vmo)
			if err != nil {
				return nil, errors.Wrap(err, "invalid move")
			} else if move < 0 {
				return nil, errors.New("invalid move")
			}
			moves = append(moves, move)
		}
	}

	return New(player, scores, pits, moves, GameStatus(status))
}

func (b *Board) Move(pit int) (*Board, error) {
	if b.pits[pit] == 0 {
		return b, errors.New("cannot make a move on an empty pit")
	}

	nb := b.clone()
	nb.tryMove(pit)
	nb.tryEndGame()

	return nb, nil
}

func (b *Board) GetValidMoves() []int {
	return b.validMoves
}

func (b *Board) Player() int {
	return b.player
}

func (b *Board) CurrentPlayerWon() bool {
	// Invert current player since the current state updates after the current player finished the move
	return (b.Status == Player1Won && b.player == 1) || (b.Status == Player2Won && b.player == 0)
}

func (b *Board) tryEndGame() {
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

func (b *Board) computeStatus() GameStatus {
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

func (b *Board) computeValidMoves() []int {
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

func (b *Board) tryMove(pit int) {
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

func (b *Board) validateInputs() error {
	if len(b.pits) != 12 {
		return errors.New("invalid number of pits")
	}

	if len(b.scores) != 2 {
		return errors.New("invalid number of scores")
	}

	if b.player > 1 || b.player < 0 {
		return errors.New("invalid player")
	}

	if len(b.validMoves) > 6 {
		return errors.New("too many valid moves")
	}

	return nil
}

func (b *Board) applyCaptures(pit int) {
	cp := pit
	scores := append([]int{}, b.scores...)
	pits := append([]int{}, b.pits...)
	for cp > 0 && b.isOpponentPit(cp) && (b.pits[cp] == 2 || b.pits[cp] == 3) {
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

func (b *Board) opponentCanMakeMove() bool {
	return (b.player == 0 && sum(b.pits[6:11]) > 0) || (b.player == 1 && sum(b.pits[0:5]) > 0)
}

func (b *Board) clone() *Board {
	return &Board{
		Status:     b.Status,
		player:     b.player,
		scores:     append([]int{}, b.scores...),
		pits:       append([]int{}, b.pits...),
		validMoves: append([]int{}, b.validMoves...),
	}
}

func (b *Board) isOpponentPit(pit int) bool {
	return pit > 5 && b.player == 0 || pit <= 5 && b.player == 1
}

func sum(s []int) int {
	total := 0
	for _, i := range s {
		total += i
	}

	return total
}
