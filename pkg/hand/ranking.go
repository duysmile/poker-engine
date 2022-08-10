package hand

// A Ranking is one of the ten possible hand rankings that determine the
// value of a hand.  Hand rankings are composed of different arrangments of
// pairs, straights, and flushes.
type Ranking int

type validFunc func([]Card, Config) bool
type descFunc func([]Card) string

type ranking struct {
	r     Ranking
	vFunc validFunc
	dFunc descFunc
}

func NewRanking(r Ranking, vFunc validFunc, dFunc descFunc) ranking {
	return ranking{
		r:     r,
		vFunc: vFunc,
		dFunc: dFunc,
	}
}

func hasFlush(cards []Card) bool {
	if len(cards) != 5 {
		return false
	}
	suit := cards[0].Suit()
	has := true
	for _, c := range cards {
		has = has && c.Suit() == suit
	}
	return has
}

func hasStraight(cards []Card) bool {
	if len(cards) != 5 {
		return false
	}
	lastIndex := cards[0].Rank()
	straight := true
	for i := 1; i < 5; i++ {
		index := cards[i].Rank()
		straight = straight && (lastIndex == index+1)
		lastIndex = index
	}
	return straight || hasLowStraight(cards)
}

func hasLowStraight(cards []Card) bool {
	return cards[0].Rank() == Five &&
		cards[1].Rank() == Four &&
		cards[2].Rank() == Three &&
		cards[3].Rank() == Two &&
		cards[4].Rank() == Ace
}

func hasPairs(cards []Card, pairNums []int) bool {
	for i := 0; i < 5; i++ {
		num := pairNums[i]
		if i >= len(cards) {
			return num == 1
		}
		card := cards[i]
		if num != len(cardsForRank(cards, card.Rank())) {
			return false
		}
	}
	return true
}

func getRankingsByType(gameType GameType) []ranking {
	switch gameType {
	case GameTypeShortDeck:
		return []ranking{sdHighCard, sdPair, sdTwoPair, sdThreeOfAKind,
			sdStraight, sdFlush, sdFullHouse, sdFourOfAKind, sdStraightFlush, sdRoyalFlush}
	case GameTypeStandard:
		return []ranking{stdHighCard, stdPair, stdTwoPair, stdThreeOfAKind,
			stdStraight, stdFlush, stdFullHouse, stdFourOfAKind, stdStraightFlush, stdRoyalFlush}
	default:
		return []ranking{stdHighCard, stdPair, stdTwoPair, stdThreeOfAKind,
			stdStraight, stdFlush, stdFullHouse, stdFourOfAKind, stdStraightFlush, stdRoyalFlush}
	}

}
