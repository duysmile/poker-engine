package hand

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/notnil/joker/util"
)

// Sorting is the sorting used to determine which hand is
// selected.
type Sorting int

const (
	// SortingHigh is a sorting method that will return the "high hand"
	SortingHigh Sorting = iota + 1

	// SortingLow is a sorting method that will return the "low hand"
	SortingLow
)

// Ordering is used to order the output of the Sort function
type Ordering int

const (
	// ASC is ascending order
	ASC Ordering = iota + 1

	// DESC is ascending order
	DESC
)

// Config represents the configuration options for hand selection
type Config struct {
	sorting         Sorting
	ignoreStraights bool
	ignoreFlushes   bool
	aceIsLow        bool
	gameType        GameType
}

type configJSON struct {
	Sorting         Sorting `json:"sorting"`
	IgnoreStraights bool    `json:"ignoreStraights"`
	IgnoreFlushes   bool    `json:"ignoreFlushes"`
	AceIsLow        bool    `json:"aceIsLow"`
}

// MarshalJSON implements the json.Marshaler interface.
func (c *Config) MarshalJSON() ([]byte, error) {
	m := &configJSON{
		Sorting:         c.sorting,
		IgnoreStraights: c.ignoreStraights,
		IgnoreFlushes:   c.ignoreFlushes,
		AceIsLow:        c.aceIsLow,
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *Config) UnmarshalJSON(b []byte) error {
	m := &configJSON{}
	if err := json.Unmarshal(b, m); err != nil {
		return err
	}
	c.sorting = m.Sorting
	c.ignoreStraights = m.IgnoreStraights
	c.ignoreFlushes = m.IgnoreFlushes
	c.aceIsLow = m.AceIsLow
	return nil
}

// Low configures NewHand to select the lowest hand in which aces
// are high and straights and flushes are counted.
func Low(c *Config) {
	c.sorting = SortingLow
}

// AceToFiveLow configures NewHand to select the lowest hand in which
// aces are low and straights and flushes aren't counted.
func AceToFiveLow(c *Config) {
	c.sorting = SortingLow
	c.aceIsLow = true
	c.ignoreStraights = true
	c.ignoreFlushes = true
}

func ShortDeck(c *Config) {
	c.gameType = GameTypeShortDeck
}

// A Hand is the highest poker hand derived from five or more cards.
type Hand struct {
	ranking     Ranking
	cards       []Card
	description string
	config      *Config
}

// New forms a hand from the given cards and configuration
// options.  If there are more than five cards, New will return
// the winning hand out of all five card combinations.  If there are
// less than five cards, the best ranking will be calculated for the
// cards given.
func New(cards []Card, options ...func(*Config)) *Hand {
	c := &Config{}
	for _, option := range options {
		option(c)
	}
	combos := cardCombos(cards)
	hands := []*Hand{}
	for _, combo := range combos {
		hand := handForFiveCards(combo, *c)
		hands = append(hands, hand)
	}
	hands = Sort(c.sorting, DESC, hands...)
	hands[0].config = c
	return hands[0]
}

// Ranking returns the hand ranking of the hand.
func (h *Hand) Ranking() Ranking {
	return h.ranking
}

// Cards returns the five cards used in the best hand ranking for the hand.
func (h *Hand) Cards() []Card {
	return append([]Card{}, h.cards...)
}

// Description returns a user displayable description of the hand such as
// "full house kings full of sixes".
func (h *Hand) Description() string {
	return h.description
}

// String returns the description followed by the cards used.
func (h *Hand) String() string {
	return fmt.Sprintf("%s %v", h.Description(), h.Cards())
}

// CompareTo returns a positive value if this hand beats the other hand, a
// negative value if this hand loses to the other hand, and zero if the hands
// are equal.
func (h *Hand) CompareTo(o *Hand) int {
	if h.Ranking() != o.Ranking() {
		return int(h.Ranking()) - int(o.Ranking())
	}
	hCards := h.Cards()
	oCards := o.Cards()
	for i := 0; i < 5; i++ {
		hCard, oCard := hCards[i], oCards[i]
		hIndex, oIndex := hCard.Rank(), oCard.Rank()
		if hIndex != oIndex {
			return int(hIndex) - int(oIndex)
		}
	}
	return 0
}

type handJSON struct {
	Ranking     Ranking `json:"ranking"`
	Cards       []Card  `json:"cards"`
	Description string  `json:"description"`
	Config      *Config `json:"config"`
}

// MarshalJSON implements the json.Marshaler interface.
// The json format is:
// {"ranking":10,"cards":["A♠","K♠","Q♠","J♠","T♠"],"description":"royal flush","config":{"sorting":1,"ignoreStraights":false,"ignoreFlushes":false,"aceIsLow":false}}
func (h *Hand) MarshalJSON() ([]byte, error) {
	m := &handJSON{
		Ranking:     h.ranking,
		Cards:       h.cards,
		Description: h.description,
		Config:      h.config,
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
//  The json format is:
// {"ranking":10,"cards":["A♠","K♠","Q♠","J♠","T♠"],"description":"royal flush","config":{"sorting":1,"ignoreStraights":false,"ignoreFlushes":false,"aceIsLow":false}}
func (h *Hand) UnmarshalJSON(b []byte) error {
	m := &handJSON{}
	if err := json.Unmarshal(b, m); err != nil {
		return err
	}
	f := func(c *Config) {
		c.sorting = m.Config.sorting
		c.ignoreStraights = m.Config.ignoreStraights
		c.ignoreFlushes = m.Config.ignoreFlushes
		c.aceIsLow = m.Config.aceIsLow
	}
	cp := New(m.Cards, f)
	h.ranking = cp.ranking
	h.cards = cp.cards
	h.description = cp.description
	h.config = cp.config
	return nil
}

// Sort returns a list of hands sorted by the given sorting
func Sort(s Sorting, o Ordering, hands ...*Hand) []*Hand {
	handsCopy := make([]*Hand, len(hands))
	copy(handsCopy, hands)

	high := (o == ASC && s == SortingHigh) || (o == DESC && s == SortingLow)
	if high {
		sort.Sort(byHighHand(handsCopy))
	} else {
		sort.Sort(sort.Reverse(byHighHand(handsCopy)))
	}
	return handsCopy
}

// ByHighHand is a slice of hands sort in ascending value
type byHighHand []*Hand

// Len implements the sort.Interface interface.
func (a byHighHand) Len() int { return len(a) }

// Swap implements the sort.Interface interface.
func (a byHighHand) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less implements the sort.Interface interface.
func (a byHighHand) Less(i, j int) bool {
	iHand, jHand := a[i], a[j]
	return iHand.CompareTo(jHand) < 0
}

func handForFiveCards(cards []Card, c Config) *Hand {
	cards = formCards(cards, c)
	rankings := getRankingsByType(c.gameType)
	for _, r := range rankings {
		if r.vFunc(cards, c) {
			return &Hand{
				ranking:     r.r,
				cards:       cards,
				description: r.dFunc(cards),
			}
		}
	}
	panic("unreachable")
}

func cardCombos(cards []Card) [][]Card {
	cCombo := [][]Card{}
	l := 5
	if len(cards) < 5 {
		l = len(cards)
	}
	indexCombos := util.Combinations(len(cards), l)

	for _, combo := range indexCombos {
		cCards := []Card{}
		for _, i := range combo {
			cCards = append(cCards, cards[i])
		}
		cCombo = append(cCombo, cCards)
	}
	return cCombo
}

func formLowStraight(cards []Card) []Card {
	if len(cards) < 5 {
		return cards
	}
	has := cards[0].Rank() == Ace &&
		cards[1].Rank() == Five &&
		cards[2].Rank() == Four &&
		cards[3].Rank() == Three &&
		cards[4].Rank() == Two
	if has {
		cards = []Card{cards[1], cards[2], cards[3], cards[4], cards[0]}
	}
	return cards
}

func formLowSDStraight(cards []Card) []Card {
	if len(cards) < 5 {
		return cards
	}
	has := cards[0].Rank() == Ace &&
		cards[1].Rank() == Nine &&
		cards[2].Rank() == Eight &&
		cards[3].Rank() == Seven &&
		cards[4].Rank() == Six
	if has {
		cards = []Card{cards[1], cards[2], cards[3], cards[4], cards[0]}
	}
	return cards
}

func formCards(cards []Card, c Config) []Card {
	var ranks []Rank
	if c.aceIsLow {
		// sort cards staring w/ king
		sort.Sort(sort.Reverse(byAceLow(cards)))
		// sort ranks starting w/ king
		ranks = allRanks()
		sort.Sort(sort.Reverse(byAceLowRank(ranks)))
	} else {
		// sort cards staring w/ ace
		sort.Sort(sort.Reverse(byAceHigh(cards)))
		// sort ranks starting w/ ace
		ranks = allRanks()
		sort.Sort(sort.Reverse(byAceHighRank(ranks)))
	}

	// form cards starting w/ most paired
	formed := []Card{}
	for i := 4; i > 0; i-- {
		for _, r := range ranks {
			rCards := cardsForRank(cards, r)
			if len(rCards) == i {
				formed = append(formed, rCards...)
			}
		}
	}
	// check for low straight
	switch c.gameType {
	case GameTypeShortDeck:
		return formLowSDStraight(formed)
	case GameTypeStandard:
		fallthrough
	default:
		return formLowStraight(formed)
	}
}

func cardsForRank(cards []Card, r Rank) []Card {
	rCards := []Card{}
	for _, c := range cards {
		if c.Rank() == r {
			rCards = append(rCards, c)
		}
	}
	return rCards
}
