package jokertest_test

import (
	"github.com/notnil/joker/pkg/hand"
	"github.com/notnil/joker/pkg/jokertest"
	"testing"
)

func TestDeck(t *testing.T) {
	cards := jokertest.Cards("Qh", "Ks", "4s")
	actual := []hand.Card{hand.QueenHearts, hand.KingSpades, hand.FourSpades}
	deck := jokertest.Dealer(cards).Deck()

	for i := 0; i < len(actual); i++ {
		card := deck.Pop()
		if actual[i] != card {
			t.Fatalf("Pop() = %s; want %s; i = %d", card, actual[i], i)
		}
	}
}
