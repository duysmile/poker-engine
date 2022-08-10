package hand

import (
	"fmt"
)

const (
	// HighCard represents a hand composed of no pairs, straights, or flushes.
	// Ex: A♠ K♠ J♣ 7♥ 5♦
	StdHighCard Ranking = iota + 1

	// Pair represents a hand composed of a single pair.
	// Ex: A♠ A♣ K♣ J♥ 5♦
	StdPair

	// TwoPair represents a hand composed of two pairs.
	// Ex: A♠ A♣ J♣ J♦ 5♦
	StdTwoPair

	// ThreeOfAKind represents a hand composed of three of the same rank.
	// Ex: A♠ A♣ A♦ J♥ 5♦
	StdThreeOfAKind

	// Straight represents a hand composed of five cards of consecutive rank.
	// Ex: A♠ K♣ Q♦ J♥ T♦
	StdStraight

	// Flush represents a hand composed of five cards that share the same suit.
	// Ex: T♠ 7♠ 4♠ 3♠ 2♠
	StdFlush

	// FullHouse represents a hand composed of three of a kind and a pair.
	// Ex: 4♠ 4♣ 4♦ 2♠ 2♥
	StdFullHouse

	// FourOfAKind represents a hand composed of four cards of the same rank.
	// Ex: A♠ A♣ A♦ A♥ 5♥
	StdFourOfAKind

	// StraightFlush represents a hand composed of five cards of consecutive
	// rank that share the same suit.
	// Ex: 5♥ 4♥ 3♥ 2♥ A♥
	StdStraightFlush

	// RoyalFlush represents a hand composed of ace, king, queen, jack, and ten
	// of the same suit.
	// Ex: A♥ K♥ Q♥ J♥ T♥
	StdRoyalFlush
)

var (
	stdHighCard = NewRanking(
		StdHighCard,
		func(cards []Card, c Config) bool {
			flush := hasFlush(cards)
			straight := hasStraight(cards)
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

	stdPair = NewRanking(
		StdPair,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{2, 2, 1, 1, 1})
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("pair of %v", r.pluralName())
		},
	)

	stdTwoPair = NewRanking(
		StdTwoPair,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{2, 2, 2, 2, 1})
		},
		func(cards []Card) string {
			r1 := cards[0].Rank()
			r2 := cards[2].Rank()
			return fmt.Sprintf("two pair %v and %v", r1.pluralName(), r2.pluralName())
		},
	)

	stdThreeOfAKind = NewRanking(
		StdThreeOfAKind,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{3, 3, 3, 1, 1})
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("three of a kind %v", r.pluralName())
		},
	)

	stdStraight = NewRanking(
		StdStraight,
		func(cards []Card, c Config) bool {
			if c.ignoreStraights {
				return false
			}
			flush := hasFlush(cards)
			straight := hasStraight(cards)
			return !flush && straight
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("straight %v high", r.singularName())
		},
	)

	stdFlush = NewRanking(
		StdFlush,
		func(cards []Card, c Config) bool {
			if c.ignoreFlushes {
				return false
			}

			flush := hasFlush(cards)
			straight := hasStraight(cards)
			return flush && !straight
		},
		func(cards []Card) string {
			r1 := cards[0].Rank()
			return fmt.Sprintf("flush %v high", r1.singularName())
		},
	)

	stdFullHouse = NewRanking(
		StdFullHouse,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{3, 3, 3, 2, 2})
		},
		func(cards []Card) string {
			r1 := cards[0].Rank()
			r2 := cards[3].Rank()
			return fmt.Sprintf("full house %v full of %v", r1.pluralName(), r2.pluralName())
		},
	)

	stdFourOfAKind = NewRanking(
		StdFourOfAKind,
		func(cards []Card, c Config) bool {
			return hasPairs(cards, []int{4, 4, 4, 4, 1})
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("four of a kind %v", r.pluralName())
		},
	)

	stdStraightFlush = NewRanking(
		StdStraightFlush,
		func(cards []Card, c Config) bool {
			if c.ignoreStraights || c.ignoreFlushes {
				return false
			}
			flush := hasFlush(cards)
			straight := hasStraight(cards)
			return cards[0].Rank() != Ace && flush && straight
		},
		func(cards []Card) string {
			r := cards[0].Rank()
			return fmt.Sprintf("straight flush %v high", r.singularName())
		},
	)

	stdRoyalFlush = NewRanking(
		StdRoyalFlush,
		func(cards []Card, c Config) bool {
			if c.ignoreStraights || c.ignoreFlushes {
				return false
			}
			flush := hasFlush(cards)
			straight := hasStraight(cards)
			return cards[0].Rank() == Ace && flush && straight
		},
		func(cards []Card) string {
			return "royal flush"
		},
	)
)
