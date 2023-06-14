package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"google.golang.org/api/option"
	"google.golang.org/api/transport"

	// import google cloud vision package
	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"

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

func main() {
	// Create a new Google Cloud Storage client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println("Failed to create Google Cloud Storage client:", err)
		os.Exit(1)
	}

	// Open the image file
	imageFile, err := client.Bucket("my-bucket").Object("path/to/image.png").NewReader(ctx)
	if err != nil {
		fmt.Println("Failed to open image file:", err)
		os.Exit(1)
	}
	defer imageFile.Close()

	// Read the image file into memory
	imageBytes, err := ioutil.ReadAll(imageFile)
	if err != nil {
		fmt.Println("Failed to read image file:", err)
		os.Exit(1)
	}

	// Create a new Google Cloud Vision client
	creds, err := transport.Creds(ctx, option.WithCredentialsFile("path/to/credentials.json"))
	if err != nil {
		fmt.Println("Failed to load credentials:", err)
		os.Exit(1)
	}
	visionClient, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentials(creds))
	if err != nil {
		fmt.Println("Failed to create Google Cloud Vision client:", err)
		os.Exit(1)
	}

	// Create a new image annotation request
	image, err := vision.NewImageFromReader(bytes.NewReader(imageBytes))
	if err != nil {
		fmt.Println("Failed to create image object:", err)
		os.Exit(1)
	}
	request := &visionpb.AnnotateImageRequest{
		Image: image,
		Features: []*visionpb.Feature{
			{
				Type: visionpb.Feature_TEXT_DETECTION,
			},
		},
	}

	// Perform the image annotation request
	response, err := visionClient.AnnotateImage(ctx, request)
	if err != nil {
		fmt.Println("Failed to perform image annotation:", err)
		os.Exit(1)
	}

	// Get the OCR output
	text := response.TextAnnotations[0].Description

	// Create a Kagome tokenizer
	theTokenizer, err := prepareIPADict()
	if err != nil {
		panic(err)
	}

	// Tokenize the text
	tokens := tokenizeText(theTokenizer, text)
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
	for _, token := range tokenCoords {
		fmt.Printf("Token: %s, Begin: %d, End: %d, Row: %d, Column: %d\n", token.Text, token.Begin, token.End, token.Row, token.Column)
	}
}

// tokenizeText tokenizes the given text using the provided Kagome tokenizer.
func tokenizeText(t *tokenizer.Tokenizer, text string) []tokenizer.Token {
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
	return tokens
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

// getPlatform returns the current platform as a string
func getPlatform() string {
	switch os := runtime.GOOS; os {
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return "unsupported"

	}
}

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
