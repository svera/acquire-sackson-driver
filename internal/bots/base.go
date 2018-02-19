package bots

import (
	"encoding/json"

	"github.com/svera/acquire-sackson-driver/internal/messages"
)

type base struct {
	status messages.Status
}

func (b *base) FeedGameStatus(message json.RawMessage) error {
	var content messages.Status

	if err := json.Unmarshal(message, &content); err != nil {
		return err
	}
	b.status = content
	return nil
}
