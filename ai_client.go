package main

import (
	"encoding/json"

	"github.com/svera/acquire-sackson-driver/internal/messages"
	"github.com/svera/acquire/bots"
	acquireInterfaces "github.com/svera/acquire/interfaces"
	"github.com/svera/sackson-server/api"
)

// AIClient is a struct that implements Sackson's AI interface,
// storing data related to a specific user and provides
// several functions to send/receive data to/from a client using a websocket
// connection
type AIClient struct {
	bot acquireInterfaces.Bot
}

// FeedGameStatus updates the AI client with the current status of the game
func (c *AIClient) FeedGameStatus(message json.RawMessage) error {
	var content messages.Status

	if err := json.Unmarshal(message, &content); err != nil {
		return err
	}
	c.updateBot(content)
	return nil
}

// Play makes the AI execute an action, returning its type and the content
// of the action to be executed as a JSON message.
func (c *AIClient) Play() api.Action {
	m := c.bot.Play()
	bm := m.(bots.Message)
	return c.encodeResponse(bm)
}

func (c *AIClient) updateBot(parsed messages.Status) {
	hand := map[string]bool{}
	var corps [7]bots.CorpData
	var playerInfo bots.PlayerData
	var rivalsInfo []bots.PlayerData

	for coords, playable := range parsed.Hand {
		hand[coords] = playable
	}
	for i := range parsed.Corps {
		corps[i] = bots.CorpData{
			Name:            parsed.Corps[i].Name,
			Price:           parsed.Corps[i].Price,
			MajorityBonus:   parsed.Corps[i].MajorityBonus,
			MinorityBonus:   parsed.Corps[i].MinorityBonus,
			RemainingShares: parsed.Corps[i].RemainingShares,
			Size:            parsed.Corps[i].Size,
			Defunct:         parsed.Corps[i].Defunct,
			Tied:            parsed.Corps[i].Tied,
		}
	}
	playerInfo = bots.PlayerData{
		Cash:        parsed.PlayerInfo.Cash,
		OwnedShares: parsed.PlayerInfo.OwnedShares,
	}
	for _, rival := range parsed.RivalsInfo {
		rivalsInfo = append(rivalsInfo, bots.PlayerData{
			Cash:        rival.Cash,
			OwnedShares: rival.OwnedShares,
		})
	}

	st := bots.Status{
		Board:       parsed.Board,
		State:       parsed.State,
		Hand:        hand,
		Corps:       corps,
		PlayerInfo:  playerInfo,
		RivalsInfo:  rivalsInfo,
		IsLastRound: parsed.IsLastRound,
	}
	c.bot.Update(st)
}

func (c *AIClient) encodeResponse(m bots.Message) api.Action {
	switch m.Type {
	case bots.PlayTileResponseType:
		return c.encodePlayTile(m.Params.(bots.PlayTileResponseParams))
	case bots.NewCorpResponseType:
		return c.encodeFoundCorporation(m.Params.(bots.NewCorpResponseParams))
	case bots.BuyResponseType:
		return c.encodeBuyStock(m.Params.(bots.BuyResponseParams))
	case bots.SellTradeResponseType:
		return c.encodeSellTrade(m.Params.(bots.SellTradeResponseParams))
	case bots.UntieMergeResponseType:
		return c.encodeUntieMerge(m.Params.(bots.UntieMergeResponseParams))
	case bots.EndGameResponseType:
		return c.encodeEndGame()
	default:
		return c.encodeOut()
	}
}

func (c *AIClient) encodePlayTile(response bots.PlayTileResponseParams) api.Action {
	params := messages.PlayTile{
		Tile: response.Tile,
	}
	ser, _ := json.Marshal(params)
	return api.Action{
		Type:   messages.TypePlayTile,
		Params: ser,
	}
}

func (c *AIClient) encodeFoundCorporation(response bots.NewCorpResponseParams) api.Action {
	params := messages.NewCorp{
		CorporationIndex: response.CorporationIndex,
	}
	ser, _ := json.Marshal(params)
	return api.Action{
		Type:   messages.TypeFoundCorporation,
		Params: ser,
	}
}

func (c *AIClient) encodeBuyStock(response bots.BuyResponseParams) api.Action {
	params := messages.Buy{
		CorporationsIndexes: response.CorporationsIndexes,
	}
	ser, _ := json.Marshal(params)
	return api.Action{
		Type:   messages.TypeBuyStock,
		Params: ser,
	}
}

func (c *AIClient) encodeSellTrade(response bots.SellTradeResponseParams) api.Action {
	params := messages.SellTrade{
		CorporationsIndexes: map[string]messages.SellTradeAmounts{},
	}
	for k, v := range response.CorporationsIndexes {
		params.CorporationsIndexes[k] = messages.SellTradeAmounts{Sell: v.Sell, Trade: v.Trade}
	}
	ser, _ := json.Marshal(params)
	return api.Action{
		Type:   messages.TypeSellTrade,
		Params: ser,
	}
}

func (c *AIClient) encodeUntieMerge(response bots.UntieMergeResponseParams) api.Action {
	params := messages.UntieMerge{
		CorporationIndex: response.CorporationIndex,
	}
	ser, _ := json.Marshal(params)
	return api.Action{
		Type:   messages.TypeUntieMerge,
		Params: ser,
	}
}

func (c *AIClient) encodeEndGame() api.Action {
	return api.Action{
		Type:   messages.TypeEndGame,
		Params: nil,
	}
}

func (c *AIClient) encodeOut() api.Action {
	params := messages.ClientOut{
		Reason: "fai",
	}
	ser, _ := json.Marshal(params)
	return api.Action{
		Type:   messages.TypeClientOut,
		Params: ser,
	}
}
