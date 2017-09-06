package main

// These are the types for the messages allowed by the Acquire driver
// and describe actions performed by players in the game:
const (
	messageTypePlayTile         = "ply"
	messageTypeFoundCorporation = "ncp"
	messageTypeBuyStock         = "buy"
	messageTypeSellTrade        = "sel"
	messageTypeUntieMerge       = "unt"
	messageTypeEndGame          = "end"
)

// playTileMessageParams is a struct which defines the format of the params of
// incoming play tile messages.
//
// The following is a play tile message example:
//
//   {
//	   "typ": "ply",
//	   "par": {
//	     "til": "2A"
//	   }
//   }
type playTileMessageParams struct {
	Tile string `json:"til"`
}

// newCorpMessageParams is a struct which defines the format of the params of
// incoming new corporation messages.
//
// The following is a found corporation message example:
//
//   {
//     "typ": "ncp",
//     "par": {
//       "cor": "2"
//     }
//   }
type newCorpMessageParams struct {
	CorporationIndex int `json:"cor"`
}

// buyMessageParams is a struct which defines the format of the params of
// incoming buy messages.
//
// The following is a buy message example:
//
//   {
//     "typ": "buy",
//     "par": {
//       "cor": {
//         "0": 3,
//         "1": 0,
//         "2": 2,
//         ...
//       }
//     }
//   }
type buyMessageParams struct {
	CorporationsIndexes map[string]int `json:"cor"`
}

// sellTradeMessageParams is a struct which defines the format of the params of
// incoming sell and trade messages.
//
// The following is a sell and trade message example:
//
//   {
//     "typ": "sel",
//     "par": {
//       "cor": {
//         "0": {
//           "sel": 2,
//           "tra": 0
//         },
//         "1": {
//           "sel": 0,
//           "tra": 2
//         },
//         ...
//       }
//     }
//   }
type sellTradeMessageParams struct {
	CorporationsIndexes map[string]sellTrade `json:"cor"`
}

type sellTrade struct {
	Sell  int `json:"sel"`
	Trade int `json:"tra"`
}

// untieMergeMessageParams is a struct which defines the format of the params of
// incoming untie merge messages.
//
// The following is an untie merge message example:
//
//   {
//     "typ": "unt",
//     "par": {
//       "cor": "2"
//     }
//   }
type untieMergeMessageParams struct {
	CorporationIndex int `json:"cor"`
}
