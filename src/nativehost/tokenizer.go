package main // too small of a project to have this in separate directory

import (
	"errors"
	"fmt"
	"os"
	"strings"

	// import kagome v2 tokenizer and dictionary package
	"github.com/ikawaha/kagome-dict/dict"
	"github.com/ikawaha/kagome/v2/tokenizer"
)

type Token struct {
	Text   string
	Begin  int
	End    int
	Row    int
	Column int
}

//	func TokenizeText(text string) (string, error) {
//		t := tokenizer.New(ipa.Dict())
//		tokens := t.Tokenize(text)
//		if len(tokens) == 0 {
//			return "", fmt.Errorf("no tokens found in image")
//		}
//
//		var tokenizedText string
//		for _, token := range tokens {
//			tokenizedText += token.Surface + " "
//		}
//
//		return strings.TrimSpace(tokenizedText), nil
//	}
func TokenizeText(text string) (string, []Token, error) {
	// Create a Kagome tokenizer
	theTokenizer, err := prepareIPADict()
	if err != nil {
		panic(err)
	}

	// Tokenize the text
	tokens, err := tokenizeText(theTokenizer, text)
	if len(tokens) == 0 {
		return "", nil, errors.New("no tokens found")
	}

	// Map tokens to their corresponding coordinates
	tokenCoords := make([]Token, len(tokens))
	for i, token := range tokens {
		tokenCoords[i] = Token{
			Text:   token.Surface,
			Begin:  token.Start,
			End:    token.End,
			Row:    token.Index,
			Column: token.Position,
		}
	}

	// Print the tokens with coordinates
	var tokenizedText string
	for _, token := range tokenCoords {
		fmt.Printf("Token: %s, Begin: %d, End: %d, Row: %d, Column: %d\n", token.Text, token.Begin, token.End, token.Row, token.Column)
		tokenizedText += token.Text + " "
	}

	return tokenizedText, tokenCoords, nil
}

// tokenizeText tokenizes the given text using the provided Kagome tokenizer.
func tokenizeText(t *tokenizer.Tokenizer, text string) ([]tokenizer.Token, error) {
	var tokens []tokenizer.Token

	for _, mode := range []tokenizer.TokenizeMode{tokenizer.Normal, tokenizer.Search, tokenizer.Extended} {
		tokenized := t.Analyze(text, mode)
		tokens = append(tokens, tokenized...)
		fmt.Printf("---%s---\n", mode)
		for _, token := range tokens {
			if token.Class == tokenizer.DUMMY {
				// BOS: Begin Of Sentence, EOS: End Of Sentence.
				fmt.Printf("%s\n", token.Surface)
				continue
			}
			features := strings.Join(token.Features(), ",")
			fmt.Printf("%s\t%v\n", token.Surface, features)
		}
	}
	return tokens, nil
}

// getDictionaryDirectory returns the path to the Kagome dictionary directory based on the platform
func getDictionaryDirectory(platform string) string {
	switch platform {
	case "linux":
		// Linux dictionary directory
		return "./"
	case "windows":
		// Windows dictionary directory
		return ".\\"
	default:
		// Unsupported platform
		fmt.Println("Unsupported platform")
		os.Exit(1)
		return ""
	}
}

//// getPlatform returns the current platform as a string
//func getPlatform() string {
//	switch os := runtime.GOOS; os {
//	case "linux":
//		return "linux"
//	case "windows":
//		return "windows"
//	default:
//		return "unsupported"
//
//	}
//}

// prepareTestDict creates a test dictionary.
func prepareIPADict() (*tokenizer.Tokenizer, error) {
	plat := getPlatform()
	dicDir := getDictionaryDirectory(plat)

	dicContents, err := dict.LoadDictFile(dicDir + "ipa.dict")
	if err != nil {
		return nil, err
	}
	theTokenizer, err := tokenizer.New(dicContents, tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}
	return theTokenizer, nil
}

func initTokenizer() error {

	return nil
}
