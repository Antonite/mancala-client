package board

import (
	"github.com/pkg/errors"
)

type Board struct {
	Player int
	Scores []int
	Pits   []int
}

func (b *Board) Move(pit int) error {
	if b.Pits[pit] == 0 {
		return errors.New("cannot make a move on an empty pit")
	}

	if err := b.validateInputs(); err != nil {
		return err
	}

	pits := b.Pits
	seeds := pits[pit]
	pits[pit] = 0
	// Play seeds
	cp := pit
	for i := 0; i < seeds; i++ {
		cp = (cp + 1) % 12
		// Skip the starting pit if we made a circle or if current pit already contains max seeds
		for cp == pit || pits[cp] == 12 {
			cp = (cp + 1) % 12
		}

		pits[cp]++
	}

	b.Pits = pits
	if err := b.applyCaptures(cp); err != nil {
		return err
	}

	b.Player = (b.Player + 1) % 2
	return nil
}

func (b *Board) validateInputs() error {
	if len(b.Pits) != 12 {
		return errors.New("invalid number of pits")
	}

	if len(b.Scores) != 2 {
		return errors.New("invalid number of scores")
	}

	if b.Player > 1 || b.Player < 0 {
		return errors.New("invalid player")
	}

	return nil
}

func (b *Board) applyCaptures(pit int) error {
	scores := b.Scores
	pits := b.Pits
	cp := pit
	for cp > 0 && isOpponentPit(cp, b.Player) && (pits[cp] == 2 || pits[cp] == 3) {
		scores[b.Player] += b.Pits[cp]
		pits[cp] = 0
		cp--
	}

	if err := validatePits(pits, b.Player); err != nil {
		return err
	}

	b.Scores = scores
	b.Pits = pits
	return nil
}

func validatePits(pits []int, player int) error {
	// Cannot leave no moves for your opponent
	if (player == 0 && pits[6] == 0 && pits[7] == 0 && pits[8] == 0 && pits[9] == 0 && pits[10] == 0 && pits[11] == 0) ||
		(player == 1 && pits[0] == 0 && pits[1] == 0 && pits[2] == 0 && pits[3] == 0 && pits[4] == 0 && pits[5] == 0) {
		return errors.New("move will leave no seeds for opponent")
	}

	return nil
}

func isOpponentPit(pit int, player int) bool {
	return pit > 5 && player == 0 || pit <= 5 && player == 1
}
