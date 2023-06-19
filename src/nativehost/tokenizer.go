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
	Text      string
	Start     int
	End       int
	Index     int
	ByteIndex int
}

var dictPath string

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
			Text:      token.Surface,
			Start:     token.Start,
			End:       token.End,
			Index:     token.Index,
			ByteIndex: token.Position,
		}
	}

	// Print the tokens with coordinates
	var tokenizedText string
	for _, token := range tokenCoords {
		fmt.Printf("Token: %s, Begin: %d, End: %d, Row: %d, Column: %d\n", token.Text, token.Start, token.End, token.Index, token.ByteIndex)
		tokenizedText += token.Text + " "
	}

	return tokenizedText, tokenCoords, nil
}

// tokenizeText tokenizes the given text using the provided Kagome tokenizer.
func tokenizeText(t *tokenizer.Tokenizer, text string) ([]tokenizer.Token, error) {
	// for return, we only need the normal mode
	tokens := t.Analyze(text, tokenizer.Normal)

	// For debugging purposes, print the tokens for each mode
	for _, mode := range []tokenizer.TokenizeMode{tokenizer.Search, tokenizer.Extended} {
		tokenized := t.Analyze(text, mode)
		//tokens = append(tokens, tokenized...)
		fmt.Printf("---%s---\n", mode)
		for _, token := range tokenized {
			if token.Class == tokenizer.DUMMY {
				// BOS: Begin Of Sentence, EOS: End Of Sentence.
				fmt.Printf("%s\n", token.Surface)
				continue
			}
			features := strings.Join(token.Features(), ",")
			fmt.Printf("%s\t%v\n", token.Surface, features)
		}
	}

	fmt.Printf("---%s---\n", tokenizer.Normal)
	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			// BOS: Begin Of Sentence, EOS: End Of Sentence.
			fmt.Printf("%s\n", token.Surface)
			continue
		}
		features := strings.Join(token.Features(), ",")
		fmt.Printf("%s\t%v\n", token.Surface, features)
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
	//plat := getPlatform()
	//dicDir := getDictionaryDirectory(plat)

	dicContents, err := dict.LoadDictFile(dictPath + "ipa.dict")
	if err != nil {
		return nil, err
	}
	theTokenizer, err := tokenizer.New(dicContents, tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}
	return theTokenizer, nil
}

func InitTokenizer(dirctDir string) error {
	dictPath = dirctDir
	if dictPath == "" {
		dictPath = getDictionaryDirectory(getPlatform())
	}
	// verify that the dictionary directory exists
	if _, err := os.Stat(dictPath); os.IsNotExist(err) {
		return fmt.Errorf("dictionary directory does not exist: %s", dictPath)
	}
	return nil
}
