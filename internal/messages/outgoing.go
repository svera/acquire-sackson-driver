package messages

// This file specifies messages sent from the hub to the clients, basically notifying
// about the status of the game after a player action.

// Status is a struct which contains the status of the game at the moment
// it is issued. It is sent to each player after every action made by one of them.
//
//   {
//      "typ": "upd", // Type: update
//      "cnt": {
//        "brd": { // Board state
//          "1A": "empty", // Board cell 1A is empty
//          "1B": "unincorporated", // Board cell 1B is unincorporated
//          "1C": "empty",
//          "1D": "0", // Board cell 1A belongs to corporation 0
//          ...
//        }
//        "sta": "PlayTile",
//        "hnd": {
//          "1A": true, // Player has tile 1A and it is playable
//        },
//        "cor": [
//          {
//            "nam": "Hilton",
//            "prc": 100, // Corporation stock price
//            "maj": 400, // Corporation majority bonus
//            "min": 200, // Corporation minority bonus
//            "rem": 20,  // Remaining stock shares
//            "siz": 2,   // Corporation size
//            "def": false, // Is corporation defunct? (in corporation merges)
//            "tie": false, // Is corporation part of a tied merge?
//          },
//          ...
//        ],
//        "ply": {
//          "nam": "John",
//          "trn": true, // Is player currently in turn?
//          "csh": 6000, // Player cash
//          "own": [     // Player owned shares per corporation
//            0: 2,
//            1: 0,
//            ...
//          ]
//        },
//        "riv": [
//          {
//            "nam": "Doe",
//            "trn": false,
//            "csh": 6000,
//            "own": [
//              0: 2,
//              1: 0,
//              ...
//            ]
//          },
//          ...
//        ],
//        "rnd": 3, // Round number
//        "lst": false, // Is last round?
//        "his": [ // History log (i18n enabled)
//          {
//            "key": "translation_key",
//            "arg": {
//              "argument_name": "argument_value",
//              ...
//            }
//          },
//          ...
//        ]
//      }
//   }
type Status struct {
	Board       map[string]string `json:"brd"`
	State       string            `json:"sta"`
	Hand        map[string]bool   `json:"hnd"`
	Corps       [7]CorpData       `json:"cor"`
	PlayerInfo  PlayerData        `json:"ply"`
	RivalsInfo  []PlayerData      `json:"riv"`
	RoundNumber int               `json:"rnd"`
	IsLastRound bool              `json:"lst"`
	History     []I18n            `json:"his"`
}

// CorpData stores all corporation information
type CorpData struct {
	Name            string `json:"nam"`
	Price           int    `json:"prc"`
	MajorityBonus   int    `json:"maj"`
	MinorityBonus   int    `json:"min"`
	RemainingShares int    `json:"rem"`
	Size            int    `json:"siz"`
	Defunct         bool   `json:"def"`
	Tied            bool   `json:"tie"`
}

// PlayerData stores all player information
type PlayerData struct {
	Name        string `json:"nam"`
	InTurn      bool   `json:"trn"`
	Cash        int    `json:"csh"`
	OwnedShares [7]int `json:"own"`
}

// I18n stores strings to be translated by the frontend, as well as related variables.
type I18n struct {
	Key       string            `json:"key"`
	Arguments map[string]string `json:"arg"`
}
