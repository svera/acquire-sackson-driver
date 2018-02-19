package bots

import (
	"errors"

	"github.com/svera/sackson-server/api"
)

const (
	// BotNotFound is an error message returned when trying to instance un inexistent bot.
	BotNotFound = "bot_not_found"
)

// Create returns a new instance of a bot.
func Create(level string) (api.AI, error) {
	switch level {
	case "chaotic":
		return NewChaotic(), nil
	default:
		return nil, errors.New(BotNotFound)
	}
}
