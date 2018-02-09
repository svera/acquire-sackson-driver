package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/svera/acquire-sackson-driver/internal/corporation"
	"github.com/svera/acquire-sackson-driver/internal/messages"
	"github.com/svera/acquire-sackson-driver/internal/player"
	acquireInterfaces "github.com/svera/acquire/interfaces"
)

// Status return a status message with the current status of the game
// as well as player specific information
func (b *AcquireDriver) Status(playerNumber int) (interface{}, error) {
	var msg interface{}

	if !b.GameStarted() {
		return nil, errors.New(GameNotStarted)
	}

	playerInfo, rivalsInfo, err := b.playersInfo(playerNumber)
	if err != nil {
		return messages.Status{}, err
	}
	msg = messages.Status{
		Board:       b.boardOwnership(),
		State:       b.game.GameStateName(),
		Corps:       b.corpsData(),
		Hand:        b.tilesData(b.players[playerNumber]),
		PlayerInfo:  playerInfo,
		RivalsInfo:  rivalsInfo,
		RoundNumber: b.game.Round(),
		IsLastRound: b.game.IsLastRound(),
		History:     b.history,
	}
	return msg, err
}

func (b *AcquireDriver) boardOwnership() map[string]string {
	cells := make(map[string]string)
	var letters = [9]string{"A", "B", "C", "D", "E", "F", "G", "H", "I"}
	for number := 1; number < 13; number++ {
		for _, letter := range letters {
			cell := b.game.Board().Cell(number, letter)
			if cell.Type() == "corporation" {
				cells[strconv.Itoa(number)+letter] = fmt.Sprintf("%d", cell.(*corporation.Corporation).Index())
			} else {
				cells[strconv.Itoa(number)+letter] = cell.Type()
			}
		}
	}

	return cells
}

func (b *AcquireDriver) corpsData() [7]messages.CorpData {
	var data [7]messages.CorpData
	for i, corp := range b.corporations {
		data[i] = messages.CorpData{
			Name:            corp.(*corporation.Corporation).Name(),
			Price:           corp.StockPrice(),
			MajorityBonus:   corp.MajorityBonus(),
			MinorityBonus:   corp.MinorityBonus(),
			RemainingShares: corp.Stock(),
			Size:            corp.Size(),
			Defunct:         b.game.IsCorporationDefunct(corp),
			Tied:            false,
		}
	}

	for _, corp := range b.game.TiedCorps() {
		data[corp.(*corporation.Corporation).Index()].Tied = true
	}
	return data
}

func (b *AcquireDriver) tilesData(pl acquireInterfaces.Player) map[string]bool {
	hnd := map[string]bool{}
	var coords string

	for _, tl := range pl.Tiles() {
		coords = strconv.Itoa(tl.Number()) + tl.Letter()
		hnd[coords] = b.game.IsTilePlayable(tl)
	}
	return hnd
}

func (b *AcquireDriver) playersInfo(n int) (messages.PlayerData, []messages.PlayerData, error) {
	rivals := []messages.PlayerData{}
	var ply messages.PlayerData
	var err error

	if _, exist := b.players[n]; !exist {
		err = errors.New(NonexistentPlayer)
	}

	for i, p := range b.players {
		if n != i {
			rivals = append(rivals, messages.PlayerData{
				Name:        p.(*player.Player).Name(),
				Cash:        p.Cash(),
				OwnedShares: b.playersShares(i),
				InTurn:      b.isCurrentPlayer(i),
			})
		} else {
			ply = messages.PlayerData{
				Name:        p.(*player.Player).Name(),
				Cash:        p.Cash(),
				OwnedShares: b.playersShares(n),
				InTurn:      b.isCurrentPlayer(n),
			}
		}
	}
	return ply, rivals, err
}
