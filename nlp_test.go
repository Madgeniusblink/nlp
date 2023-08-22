package nlp

import (
	"log"
	"os"
	"strings"

	"testing"

	"github.com/stretchr/testify/require"

	"github.com/BurntSushi/toml"
)

// define the structure of the TOML file
type tokenizeCase struct {
	Text   string
	Tokens []string
}

func loadTokenizeCases(t *testing.T) []tokenizeCase {
	data, err := os.ReadFile("testdata/tokenize_cases.toml")
	require.NoError(t, err, "failed to load tokenize cases")
	if err != nil {
		log.Fatal(err)
	}
	var testCases struct {
		Cases []tokenizeCase `toml:"cases"`
	}
	err = toml.Unmarshal(data, &testCases)
	require.NoError(t, err, "failed to parse tokenize cases")
	return testCases.Cases

}

// a two step process to load the test cases, usually more flexible if you try to use the data else where.
func TestTokenizeTable(t *testing.T) {
	for _, tc := range loadTokenizeCases(t) {
		t.Run(tc.Text, func(t *testing.T) {
			tokens := Tokenize(tc.Text)
			require.Equal(t, tc.Tokens, tokens)
		})
	}
}

// straight forward way to load the test cases
func TestTokenizeToml(t *testing.T) {

	var data struct {
		cases []tokenizeCase `toml:"cases"`
	}

	_, err := toml.DecodeFile("testdata/tokenize_cases.toml", &data)
	require.NoError(t, err)
	if err != nil {
		log.Fatal(err)
	}

	for _, tc := range data.cases {
		t.Run(tc.Text, func(t *testing.T) {
			tokens := Tokenize(tc.Text)
			require.Equal(t, tc.Tokens, tokens)
		})
	}
}

func TestTokenize(t *testing.T) {
	text := "What's on second?"
	expected := []string{"what", "on", "second"}
	tokens := Tokenize(text)
	require.Equal(t, expected, tokens)
	/*
		if !reflect.DeepEqual(tokens, expected) {
			t.Fatalf("Expected %v, got %v", expected, tokens)
		}
	*/
}

func FuzzTokenize(f *testing.F) {
	f.Fuzz(func(t *testing.T, text string) {
		tokens := Tokenize(text)
		lText := strings.ToLower(text)
		for _, tok := range tokens {
			if !strings.Contains(lText, tok) {
				t.Fatal(tok)
			}
		}
	})
}
