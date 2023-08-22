package nlp_test

import (
	"fmt"

	"github.com/madgeniusblink/nlp"
)

func ExampleTokenize() {
	text := "who's on first?"
	tokens := nlp.Tokenize(text)
	fmt.Println(tokens)

	//output:
	// [who on first]
}
