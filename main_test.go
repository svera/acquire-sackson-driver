package main

import (
	"encoding/json"
	"testing"

	"github.com/svera/sackson-server/api"
)

func TestParseNonExistingTypeMessage(t *testing.T) {
	driver := New().(*AcquireDriver)
	err := driver.Execute(api.Action{PlayerName: "Test client", Type: "err", Params: json.RawMessage{}})
	if err == nil {
		t.Errorf("Driver must return an error when receiving a non-existing message type")
	}
}

func TestParseWrongTypeMessage(t *testing.T) {
	driver := New().(*AcquireDriver)
	playerNames := map[int]string{0: "test1", 1: "test2", 2: "test3"}

	driver.StartGame(playerNames)
	data := []byte(`{"aaa": "bbb"}`)
	raw := (json.RawMessage)(data)

	err := driver.Execute(api.Action{PlayerName: "Test client", Type: "ply", Params: raw})
	if err == nil {
		t.Errorf("Driver must return an error when receiving a malformed message")
	}

	err = driver.Execute(api.Action{PlayerName: "Test client", Type: "ncp", Params: raw})
	if err == nil {
		t.Errorf("Driver must return an error when receiving a malformed message")
	}

	err = driver.Execute(api.Action{PlayerName: "Test client", Type: "buy", Params: raw})
	if err == nil {
		t.Errorf("Driver must return an error when receiving a malformed message")
	}

	err = driver.Execute(api.Action{PlayerName: "Test client", Type: "sel", Params: raw})
	if err == nil {
		t.Errorf("Driver must return an error when receiving a malformed message")
	}

	err = driver.Execute(api.Action{PlayerName: "Test client", Type: "unt", Params: raw})
	if err == nil {
		t.Errorf("Driver must return an error when receiving a malformed message")
	}

	err = driver.Execute(api.Action{PlayerName: "Test client", Type: "end", Params: raw})
	if err == nil {
		t.Errorf("Driver must return an error when receiving a malformed message")
	}
}

func TestCurrentPlayerNumbersWithoutGameStarted(t *testing.T) {
	driver := New().(*AcquireDriver)
	if _, err := driver.CurrentPlayersNumbers(); err == nil {
		t.Errorf("Driver must return an error when trying to get the current players numbers without a game started")
	}
}

func TestCurrentPlayerNumbersWithGameStarted(t *testing.T) {
	driver := New().(*AcquireDriver)
	playerNames := map[int]string{0: "test1", 1: "test2", 2: "test3"}

	driver.StartGame(playerNames)

	if _, err := driver.CurrentPlayersNumbers(); err != nil {
		t.Errorf("Driver must not return an error when trying to get the current players numbers of a started game")
	}
}

func TestStartGameWithNotEnoughPlayers(t *testing.T) {
	driver := New().(*AcquireDriver)
	playerNames := map[int]string{0: "test1", 1: "test2"}

	if err := driver.StartGame(playerNames); err == nil {
		t.Errorf("Driver must return an error when trying to start a game with not enough players")
	}
}

func TestStartGameWithEnoughPlayers(t *testing.T) {
	driver := New().(*AcquireDriver)
	playerNames := map[int]string{0: "test1", 1: "test2", 2: "test3"}

	if err := driver.StartGame(playerNames); err != nil {
		t.Errorf("Driver must not return an error when trying to start a game with enough players")
	}
}

func TestStatusWithGameStarted(t *testing.T) {
	driver := New().(*AcquireDriver)
	playerNames := map[int]string{0: "test1", 1: "test2", 2: "test3"}

	driver.StartGame(playerNames)

	if _, err := driver.Status(0); err != nil {
		t.Errorf("Driver must not return an error when trying to get the status of a started game")
	}
}

func TestStatusWithGameNotStarted(t *testing.T) {
	driver := New().(*AcquireDriver)
	if _, err := driver.Status(0); err == nil {
		t.Errorf("Driver must return an error when trying to get the status of a non started game")
	}
}

func TestStatusForNonexistentPlayer(t *testing.T) {
	driver := New().(*AcquireDriver)
	playerNames := map[int]string{0: "test1", 1: "test2", 2: "test3"}

	driver.StartGame(playerNames)
	if _, err := driver.Status(9); err == nil {
		t.Errorf("Driver must return an error when trying to get the game status of an nonexistent player")
	}
}
