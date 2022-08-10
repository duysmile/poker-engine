package main

import (
	"fmt"
	"github.com/notnil/joker/pkg/hand"
	"math/rand"
)

func main() {
	r := rand.New(rand.NewSource(0))
	deck := hand.NewDealer(r, hand.GameTypeShortDeck).Deck()
	h1 := hand.New(deck.PopMulti(5), hand.ShortDeck)
	h2 := hand.New(deck.PopMulti(5), hand.ShortDeck)

	fmt.Println(h1)
	fmt.Println(h2)

	hands := hand.Sort(hand.SortingHigh, hand.DESC, h1, h2)
	fmt.Println("Winner is:", hands[0].Cards())
}
