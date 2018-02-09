package main

import (
	"errors"

	"github.com/svera/acquire-sackson-driver/internal/messages"
)

func (b *AcquireDriver) claimEndGame(clientName string) error {
	if !b.game.ClaimEndGame().IsLastRound() {
		return errors.New(NotEndGame)
	}
	b.history = append(b.history, messages.I18n{
		Key: "game.history.claimed_end",
		Arguments: map[string]string{
			"player": clientName,
		},
	})

	return nil
}
