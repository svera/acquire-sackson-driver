package main

import (
	"errors"
	"strconv"

	"github.com/svera/acquire-sackson-driver/corporation"
	acquireInterfaces "github.com/svera/acquire/interfaces"
)

func (b *AcquireDriver) buyStock(clientName string, params buyMessageParams) error {
	buy := map[acquireInterfaces.Corporation]int{}

	for corpIndex, amount := range params.CorporationsIndexes {
		index, _ := strconv.Atoi(corpIndex)
		if index < 0 || index > 6 {
			return errors.New(CorporationNotFound)
		}

		buy[b.corporations[index]] = amount
	}

	if err := b.game.BuyStock(buy); err != nil {
		return err
	}
	for corp, amount := range buy {
		if amount > 0 {
			b.history = append(b.history, i18n{
				Key: "game.history.bought_stock",
				Arguments: map[string]string{
					"player":      clientName,
					"amount":      strconv.Itoa(amount),
					"corporation": corp.(*corporation.Corporation).Name(),
				},
			})
		}
	}
	return nil
}
