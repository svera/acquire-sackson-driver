package main

import (
	"encoding/json"
	"errors"

	"github.com/svera/acquire"
	"github.com/svera/acquire-sackson-driver/corporation"
	"github.com/svera/acquire-sackson-driver/player"
	"github.com/svera/acquire/bots"
	acquireInterfaces "github.com/svera/acquire/interfaces"
)

// AcquireDriver implements the driver interface in order to be able to have
// and acquire game through the turn based game server
type AcquireDriver struct {
	game         *acquire.Game
	players      map[int]acquireInterfaces.Player
	corporations [7]acquireInterfaces.Corporation
	history      []i18n
}

// NotEndGame defines the message returned when a player claims wrongly that end game conditions have been met
const NotEndGame = "not_end_game"

// WrongMessage defines the message returned when AcquireDriver receives a bad formed message
const WrongMessage = "message_parsing_error"

// GameAlreadyStarted is an error returned when a player tries to start a game in a hub instance which an already running one
const GameAlreadyStarted = "game_already_started"

// GameNotStarted is an error returned when a player tries to do an action that requires a running game
const GameNotStarted = "game_not_started"

// NonexistentPlayer is an error returned when someone tries to remove or get information of a non existent player
const NonexistentPlayer = "nonexistent_player"

// CorporationNotFound is an error returned when someone tries to use a non existent corporation
const CorporationNotFound = "corporation_not_found"

// New initializes a new AcquireDriver instance
func New() interface{} {
	return &AcquireDriver{
		corporations: defaultCorporations(),
	}
}

// Execute gets an input JSON-encoded message and parses it, executing
// whatever actions are required by it
func (b *AcquireDriver) Execute(clientName string, t string, params json.RawMessage) error {
	var err error
	b.history = nil

	switch t {
	case messageTypePlayTile:
		var parsed playTileMessageParams
		if err = json.Unmarshal(params, &parsed); err == nil {
			err = b.playTile(clientName, parsed)
		}
	case messageTypeFoundCorporation:
		var parsed newCorpMessageParams
		if err = json.Unmarshal(params, &parsed); err == nil {
			err = b.foundCorporation(clientName, parsed)
		}
	case messageTypeBuyStock:
		var parsed buyMessageParams
		if err = json.Unmarshal(params, &parsed); err == nil {
			err = b.buyStock(clientName, parsed)
		}
	case messageTypeSellTrade:
		var parsed sellTradeMessageParams
		if err = json.Unmarshal(params, &parsed); err == nil {
			err = b.sellTrade(clientName, parsed)
		}
	case messageTypeUntieMerge:
		var parsed untieMergeMessageParams
		if err = json.Unmarshal(params, &parsed); err == nil {
			err = b.untieMerge(clientName, parsed)
		}
	case messageTypeEndGame:
		err = b.claimEndGame(clientName)
	default:
		err = errors.New(WrongMessage)
	}

	return err
}

// CurrentPlayersNumbers returns a slice containing the number of each player currently in turn
func (b *AcquireDriver) CurrentPlayersNumbers() ([]int, error) {
	currentPlayersNumbers := []int{}
	if !b.GameStarted() {
		return currentPlayersNumbers, errors.New(GameNotStarted)
	}
	currentPlayersNumbers = append(currentPlayersNumbers, b.game.CurrentPlayer().Number())
	return currentPlayersNumbers, nil
}

// GameStarted returns true if there's a game in progress, false otherwise
func (b *AcquireDriver) GameStarted() bool {
	if b.game == nil {
		return false
	}
	return true
}

func (b *AcquireDriver) isCurrentPlayer(n int) bool {
	if b.game.CurrentPlayer().Number() == n {
		return true
	}
	return false
}

func (b *AcquireDriver) playersShares(playerNumber int) [7]int {
	var data [7]int
	for i, corp := range b.game.Corporations() {
		data[i] = b.players[playerNumber].Shares(corp)
	}
	return data
}

// RemovePlayer removes a player from the game
func (b *AcquireDriver) RemovePlayer(number int) error {
	if _, exists := b.players[number]; !exists {
		return errors.New(NonexistentPlayer)
	}
	playerName := b.players[number].(*player.Player).Name()
	b.game.RemovePlayer(b.players[number])
	delete(b.players, number)
	b.history = append([]i18n{}, i18n{
		Key: "game.history.player_left",
		Arguments: map[string]string{
			"player": playerName,
		},
	})
	return nil
}

// StartGame starts a new Acquire game
func (b *AcquireDriver) StartGame(clientNames map[int]string) error {
	var err error

	if b.GameStarted() {
		err = errors.New(GameAlreadyStarted)
	}

	b.addPlayers(clientNames)

	if b.game, err = acquire.New(b.players, acquire.Optional{Corporations: b.corporations}); err == nil {
		b.history = append(b.history, i18n{
			Key: "game.history.starter_player",
			Arguments: map[string]string{
				"player": b.currentPlayerName(),
			},
		})
	}
	return err
}

// addPlayers adds players to the game
func (b *AcquireDriver) addPlayers(clientNames map[int]string) {
	b.players = make(map[int]acquireInterfaces.Player)

	for n, playerName := range clientNames {
		b.players[n] = player.New(playerName, n)
	}
}

func (b *AcquireDriver) currentPlayerName() string {
	currentPlayerNumber := b.game.CurrentPlayer().Number()
	return b.players[currentPlayerNumber].(*player.Player).Name()
}

// IsGameOver returns true if the game has reached its end or there are not
// enough players to continue playing
func (b *AcquireDriver) IsGameOver() bool {
	if b.GameStarted() {
		return b.game.GameStateName() == acquireInterfaces.EndGameStateName ||
			b.game.GameStateName() == acquireInterfaces.InsufficientPlayersStateName ||
			b.game.GameStateName() == acquireInterfaces.ErrorStateName
	}
	return false
}

// CreateAI create an instance of an AI of the passed level
func (b *AcquireDriver) CreateAI(params interface{}) (interface{}, error) {
	var err error
	var bot acquireInterfaces.Bot
	if level, ok := params.(string); ok {
		if bot, err = bots.Create(level); err == nil {
			return &AIClient{
				bot: bot,
			}, nil
		}
		return nil, err
	}
	panic("Expecting string in CreateAI parameter")
}

func defaultCorporations() [7]acquireInterfaces.Corporation {
	var corporations [7]acquireInterfaces.Corporation
	corpsParams := [7]string{
		"Sackson",
		"Zeta",
		"Hydra",
		"Fusion",
		"America",
		"Phoenix",
		"Quantum",
	}

	for i, corpName := range corpsParams {
		corporations[i] = corporation.New(corpName, i)
	}
	return corporations
}

func main() {}
