package messages

// These are the types for the messages allowed by the Acquire driver
// and describe actions performed by players in the game:
const (
	TypePlayTile         = "ply"
	TypeFoundCorporation = "ncp"
	TypeBuyStock         = "buy"
	TypeSellTrade        = "sel"
	TypeUntieMerge       = "unt"
	TypeEndGame          = "end"
	TypeClientOut        = "out"
)

// PlayTile is a struct which defines the content of
// incoming play tile messages.
//
// The following is a play tile message example:
//
//   {
//	   "typ": "ply",
//	   "cnt": {
//	     "til": "2A"
//	   }
//   }
type PlayTile struct {
	Tile string `json:"til"`
}

// NewCorp is a struct which defines the content of
// incoming new corporation messages.
//
// The following is a found corporation message example:
//
//   {
//     "typ": "ncp",
//     "cnt": {
//       "cor": "2"
//     }
//   }
type NewCorp struct {
	CorporationIndex int `json:"cor"`
}

// Buy is a struct which defines the content of
// incoming buy messages.
//
// The following is a buy message example:
//
//   {
//     "typ": "buy",
//     "cnt": {
//       "cor": {
//         "0": 3,
//         "1": 0,
//         "2": 2,
//         ...
//       }
//     }
//   }
type Buy struct {
	CorporationsIndexes map[string]int `json:"cor"`
}

// SellTrade is a struct which defines the content of
// incoming sell and trade messages.
//
// The following is a sell and trade message example:
//
//   {
//     "typ": "sel",
//     "cnt": {
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
type SellTrade struct {
	CorporationsIndexes map[string]SellTradeAmounts `json:"cor"`
}

// SellTradeAmounts stores the amount of stock shares to be sold or traded for a corporation
type SellTradeAmounts struct {
	Sell  int `json:"sel"`
	Trade int `json:"tra"`
}

// UntieMerge is a struct which defines the content of
// incoming untie merge messages.
//
// The following is an untie merge message example:
//
//   {
//     "typ": "unt",
//     "cnt": {
//       "cor": "2"
//     }
//   }
type UntieMerge struct {
	CorporationIndex int `json:"cor"`
}

type ClientOut struct {
	Reason string `json:"rea"`
}
