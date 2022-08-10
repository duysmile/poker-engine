joker
========

Poker hand evaluation and ranking written in go (golang)

To install run:

```
go get github.com/notnil/joker/hand
```

```go
package main

import (
	"fmt"

	"github.com/notnil/joker/hand"
)

func main() {
	r := rand.New(rand.NewSource(0))
	deck := hand.NewDealer(r, hand.GameTypeStandard).Deck()
	h1 := hand.New(deck.PopMulti(5))
	h2 := hand.New(deck.PopMulti(5))

	fmt.Println(h1)
	fmt.Println(h2)

	hands := hand.Sort(hand.SortingHigh, hand.DESC, h1, h2)
	fmt.Println("Winner is:", hands[0].Cards())
}

```
