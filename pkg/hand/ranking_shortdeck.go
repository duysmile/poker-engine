package hand

import (
	"fmt"
)

const (
	// SDHighCard represents a hand composed of no pairs, straights, or flushes.
	// Ex: A♠ K♠ J♣ 7♥ 5♦
	SDHighCard Ranking = iota + 1

	// SDPair represents a hand composed of a single pair.
	// Ex: A♠ A♣ K♣ J♥ 5♦
	SDPair

	// SDTwoPair represents a hand composed of two pairs.
	// Ex: A♠ A♣ J♣ J♦ 5♦
	SDTwoPair

	// SDThreeOfAKind represents a hand composed of three of the same rank.
	// Ex: A♠ A♣ A♦ J♥ 5♦
	SDThreeOfAKind

	// SDStraight represents a hand composed of five cards of consecutive rank.
	// Ex: A♠ K♣ Q♦ J♥ T♦
	SDStraight

	// SDFullHouse represents a hand composed of three of a kind and a pair.
	// Ex: 4♠ 4♣ 4♦ 2♠ 2♥
	SDFullHouse

	// SDFlush represents a hand composed of five cards that share the same suit.
	// Ex: T♠ 7♠ 4♠ 3♠ 2♠
	SDFlush

	// SDFourOfAKind represents a hand composed of four cards of the same rank.
	// Ex: A♠ A♣ A♦ A♥ 5♥
	SDFourOfAKind

	// SDStraightFlush represents a hand composed of five cards of consecutive
	// rank that share the same suit.
	// Ex: 5♥ 4♥ 3♥ 2♥ A♥
	SDStraightFlush

	// SDRoyalFlush represents a hand composed of ace, king, queen, jack, and ten
	// of the same suit.
	// Ex: A♥ K♥ Q♥ J♥ T♥
	SDRoyalFlush
)

var (
	sdHighCard = NewRanking(
		SDHighCard,
		func(cards []Card, c Config) bool {
			flush := hasFlush(cards)
			straight := hasStraightInShortDeck(cards)
			pairs := hasPairs(cards, []int{1, 1, 1, 1, 1})
			if !c.ignoreStraights {
				pairs = pairs && !straight
			}
			if !c.ignoreFlushes {
				pairs = pairs && !flush
			}
			return pairs
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("high card %v high", r.singularName())
		},
	)

	sdPair = NewRanking(
		SDPair,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{2, 2, 1, 1, 1})
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("pair of %v", r.pluralName())
		},
	)

	sdTwoPair = NewRanking(
		SDTwoPair,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{2, 2, 2, 2, 1})
		},
		func(cards []Card) string {
			r1 := cards[0].Rank()
			r2 := cards[2].Rank()
			return fmt.Sprintf("two pair %v and %v", r1.pluralName(), r2.pluralName())
		},
	)

	sdThreeOfAKind = NewRanking(
		SDThreeOfAKind,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{3, 3, 3, 1, 1})
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("three of a kind %v", r.pluralName())
		},
	)

	sdStraight = NewRanking(
		SDStraight,
		func(cards []Card, c Config) bool {
			if c.ignoreStraights {
				return false
			}
			flush := hasFlush(cards)
			straight := hasStraightInShortDeck(cards)
			return !flush && straight
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("straight %v high", r.singularName())
		},
	)

	sdFlush = NewRanking(
		SDFlush,
		func(cards []Card, c Config) bool {
			if c.ignoreFlushes {
				return false
			}

			flush := hasFlush(cards)
			straight := hasStraightInShortDeck(cards)
			return flush && !straight
		},
		func(cards []Card) string {
			r1 := cards[0].Rank()
			return fmt.Sprintf("flush %v high", r1.singularName())
		},
	)

	sdFullHouse = NewRanking(
		SDFullHouse,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{3, 3, 3, 2, 2})
		},
		func(cards []Card) string {
			r1 := cards[0].Rank()
			r2 := cards[3].Rank()
			return fmt.Sprintf("full house %v full of %v", r1.pluralName(), r2.pluralName())
		},
	)

	sdFourOfAKind = NewRanking(
		SDFourOfAKind,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{4, 4, 4, 4, 1})
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("four of a kind %v", r.pluralName())
		},
	)

	sdStraightFlush = NewRanking(
		SDStraightFlush,
		func(cards []Card, c Config) bool {
			if c.ignoreStraights || c.ignoreFlushes {
				return false
			}
			flush := hasFlush(cards)
			straight := hasStraightInShortDeck(cards)
			return cards[0].Rank() != Ace && flush && straight
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("straight flush %v high", r.singularName())
		},
	)

	sdRoyalFlush = NewRanking(
		SDRoyalFlush,
		func(cards []Card, c Config) bool {
			if c.ignoreStraights || c.ignoreFlushes {
				return false
			}
			flush := hasFlush(cards)
			straight := hasStraightInShortDeck(cards)
			return cards[0].Rank() == Ace && flush && straight
		},
		func(cards []Card) string {
			return "royal flush"
		},
	)
)

func hasStraightInShortDeck(cards []Card) bool {
	return hasStraight(cards) || hasLowSDStraight(cards)
}

func hasLowSDStraight(cards []Card) bool {
	return cards[0].Rank() == Ace &&
		cards[1].Rank() == Nine &&
		cards[2].Rank() == Eight &&
		cards[3].Rank() == Seven &&
		cards[4].Rank() == Six
}
