package main

import (
	"errors"

	"github.com/svera/acquire-sackson-driver/internal/corporation"
	"github.com/svera/acquire-sackson-driver/internal/messages"
)

func (b *AcquireDriver) foundCorporation(clientName string, params messages.NewCorp) error {
	if params.CorporationIndex < 0 || params.CorporationIndex > 6 {
		return errors.New(CorporationNotFound)
	}
	corp := b.corporations[params.CorporationIndex]
	if err := b.game.FoundCorporation(corp); err != nil {
		return err
	}
	b.history = append(b.history, messages.I18n{
		Key: "game.history.founded_corporation",
		Arguments: map[string]string{
			"player":      clientName,
			"corporation": corp.(*corporation.Corporation).Name(),
		},
	})

	return nil
}
