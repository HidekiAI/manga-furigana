package main // too small of a project to have this in separate directory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"os"
)

// message is the JSON structure that is passed between the extension (via sendNativeMessage() in background.js) and the native host.
type Message struct {
	TabID int    `json:"tabId"`
	Image []byte `json:"image"`
}

type TokenizedImage struct {
	URL           string `json:"url"`
	TokenizedText string `json:"tokenized_text"`
}

// The JSON structure is defined in the manifest.json file.

// processing background.go code from main.go
func doBackground() {
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)

	for {
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Decode the image data
		perImage, _, err := image.Decode(bytes.NewReader(msg.Image))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Process the image
		tokenizedText, _, err := TokenizeImage(perImage, enc)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// Send the tokenized image back to the JavaScript code
		if err := enc.Encode(TokenizedImage{URL: "", TokenizedText: tokenizedText}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func TokenizeImage(img image.Image, enc *json.Encoder) (string, []Token, error) {
	// first, scan for Japanese text in the image using OCR
	textsPerImage, err := PerformOCRProxy(img, enc)
	if err != nil {
		return "", nil, err
	}

	// Tokenize the image using Kagome
	tokenizedText, tokens, err := TokenizeText(textsPerImage)
	if err != nil {
		return "", nil, err
	}

	return tokenizedText, tokens, nil
}
func InitBackground() error {
	return nil
}
