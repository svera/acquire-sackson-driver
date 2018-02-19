// Package bots implements different types of AI for playing Acquire games
package bots

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/svera/acquire-sackson-driver/internal/messages"
	"github.com/svera/acquire/interfaces"
	"github.com/svera/sackson-server/api"
)

const (
	endGameCorporationSize = 41
	safeCorporationSize    = 11
)

var rn *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rn = rand.New(source)
}

// Chaotic is a struct which implements a very stupid AI, which basically
// chooses all its decisions randomly (So not that much an AI but an AS)
type Chaotic struct {
	*base
}

// NewChaotic returns a new instance of the chaotic AI bot
func NewChaotic() *Chaotic {
	return &Chaotic{
		&base{},
	}
}

// Play analyses the current game status and returns a message with the
// next play movement by the bot AI
func (r *Chaotic) Play() api.Action {
	var msg api.Action

	if !r.status.IsLastRound && r.claimEndGame() {
		msg = api.Action{
			Type: messages.TypeEndGame,
		}
	} else {
		switch r.status.State {
		case interfaces.PlayTileStateName:
			ser, _ := json.Marshal(r.playTile())
			msg = api.Action{
				Type:   messages.TypePlayTile,
				Params: ser,
			}
		case interfaces.FoundCorpStateName:
			ser, _ := json.Marshal(r.foundCorporation())
			msg = api.Action{
				Type:   messages.TypeFoundCorporation,
				Params: ser,
			}
		case interfaces.BuyStockStateName:
			ser, _ := json.Marshal(r.buyStock())
			msg = api.Action{
				Type:   messages.TypeBuyStock,
				Params: ser,
			}
		case interfaces.SellTradeStateName:
			ser, _ := json.Marshal(r.sellTrade())
			msg = api.Action{
				Type:   messages.TypeSellTrade,
				Params: ser,
			}
		case interfaces.UntieMergeStateName:
			ser, _ := json.Marshal(r.untieMerge())
			msg = api.Action{
				Type:   messages.TypeUntieMerge,
				Params: ser,
			}
		}
	}

	return msg
}

func (r *Chaotic) playTile() messages.PlayTile {
	tileCoords := r.tileCoords()
	tileNumber := rn.Intn(len(tileCoords))

	return messages.PlayTile{
		Tile: tileCoords[tileNumber],
	}
}

// As the tiles in hand come as a map, we need to store its coordinates in an array
// before selecting a random one (only the playable ones)
func (r *Chaotic) tileCoords() []string {
	coords := make([]string, 0, len(r.status.Hand))
	for k, playable := range r.status.Hand {
		if playable {
			coords = append(coords, k)
		}
	}
	return coords
}

func (r *Chaotic) foundCorporation() messages.NewCorp {
	var corpNumber int
	response := messages.NewCorp{}
	for {
		corpNumber = rn.Intn(len(r.status.Corps))
		if r.status.Corps[corpNumber].Size == 0 {
			response.CorporationIndex = corpNumber
			break
		}
	}
	return response
}

// buyStock buys stock from a random active corporation
func (r *Chaotic) buyStock() messages.Buy {
	buy := 0
	var corpIndex int
	var corp messages.CorpData

	for {
		corpIndex = rn.Intn(len(r.status.Corps))
		corp = r.status.Corps[corpIndex]
		if corp.Size > 0 {
			break
		}
	}
	if corp.RemainingShares > 3 && corp.Size > 0 && r.hasEnoughCash(3, corp.Price) {
		buy = 3
	} else if corp.Size > 0 && r.hasEnoughCash(corp.RemainingShares, corp.Price) {
		buy = corp.RemainingShares
	}
	index := strconv.Itoa(corpIndex)
	return messages.Buy{
		CorporationsIndexes: map[string]int{
			index: buy,
		},
	}
}

func (r *Chaotic) hasEnoughCash(amount int, price int) bool {
	return amount*price < r.status.PlayerInfo.Cash
}

func (r *Chaotic) sellTrade() messages.SellTrade {
	var sellTrade messages.SellTrade
	sellTradeCorporations := map[string]messages.SellTradeAmounts{}

	for i, corp := range r.status.Corps {
		if corp.Defunct && r.status.PlayerInfo.OwnedShares[i] > 0 {
			index := strconv.Itoa(i)
			sellTradeCorporations[index] = messages.SellTradeAmounts{
				Sell: r.status.PlayerInfo.OwnedShares[i],
			}
		}
	}
	sellTrade.CorporationsIndexes = sellTradeCorporations
	return sellTrade
}

func (r *Chaotic) untieMerge() messages.UntieMerge {
	var untieMerge messages.UntieMerge
	for i, corp := range r.status.Corps {
		if corp.Tied {
			untieMerge = messages.UntieMerge{
				CorporationIndex: i,
			}
			break
		}
	}
	return untieMerge
}

func (r *Chaotic) claimEndGame() bool {
	var active, safe int
	for _, corp := range r.status.Corps {
		if corp.Size >= endGameCorporationSize {
			return true
		}
		if corp.Size > 0 {
			active++
		}
		if corp.Size >= safeCorporationSize {
			safe++
		}
	}
	if active > 0 && active == safe {
		return true
	}
	return false
}
