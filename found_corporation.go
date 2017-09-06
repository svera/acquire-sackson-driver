package main

import (
	"errors"

	"github.com/svera/acquire-sackson-driver/corporation"
)

func (b *AcquireDriver) foundCorporation(clientName string, params newCorpMessageParams) error {
	if params.CorporationIndex < 0 || params.CorporationIndex > 6 {
		return errors.New(CorporationNotFound)
	}
	corp := b.corporations[params.CorporationIndex]
	if err := b.game.FoundCorporation(corp); err != nil {
		return err
	}
	b.history = append(b.history, i18n{
		Key: "game.history.founded_corporation",
		Arguments: map[string]string{
			"player":      clientName,
			"corporation": corp.(*corporation.Corporation).Name(),
		},
	})

	return nil
}
