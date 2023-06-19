package main

import (
	"bytes"
	"image"
	"os"
	"testing"
)

func TestTokenizeText(t *testing.T) {
	// NOTE: Calling InitTokenizer() here is redundant, since init() (private one) already calls it, but it's here for completeness
	errtkn := InitTokenizer("../") // one directory up from this file
	if errtkn != nil {
		t.Errorf("Failed to initialize tokenizer: %v", errtkn)
	}

	// Tokenize the text
	tokenizedText, tokens, err := TokenizeText("これはテストです")
	if err != nil {
		t.Errorf("Failed to tokenize text: %v", err)
	}
	if len(tokens) == 0 {
		t.Fatal("no tokens found")
	}
	t.Logf("Tokenized text: %s", tokenizedText)

	// Declare constants expected values and compare them to the actual values
	expectedTokens := []Token{
		{
			Text:      "これ",
			Start:     0,
			End:       2,
			Index:     0,
			ByteIndex: 0,
		},
		{
			Text:      "は",
			Start:     2,
			End:       3,
			Index:     1,
			ByteIndex: 6,
		},
		{
			Text:      "テスト",
			Start:     3,
			End:       6,
			Index:     2,
			ByteIndex: 9,
		},
		{
			Text:      "です",
			Start:     6,
			End:       8,
			Index:     3,
			ByteIndex: 18,
		},
	}
	if len(tokens) != len(expectedTokens) {
		t.Errorf("Expected %d tokens, but got %d", len(expectedTokens), len(tokens))
		return
	}
	for i, token := range tokens {
		t.Logf("Token Index=%d: Text:'%s'", i, token.Text)
		if token.Text != expectedTokens[i].Text {
			t.Errorf("Expected token %d to have text %s, but got %s", i, expectedTokens[i].Text, token.Text)
		}
		if token.Index != expectedTokens[i].Index {
			t.Errorf("Expected token %d to have Index %d, but got %d", i, expectedTokens[i].Index, token.Index)
		}
		if token.ByteIndex != expectedTokens[i].ByteIndex {
			t.Errorf("Expected token %d to have ByteIndex %d, but got %d", i, expectedTokens[i].ByteIndex, token.ByteIndex)
		}
		if token.Start != expectedTokens[i].Start {
			t.Errorf("Expected token %d to have Start %d, but got %d", i, expectedTokens[i].Start, token.Start)
		}
		if token.End != expectedTokens[i].End {
			t.Errorf("Expected token %d to have End %d, but got %d", i, expectedTokens[i].End, token.End)
		}
	}
}

// Test OCR on a sample image
func TestOCR(t *testing.T) {
	// NOTE: Calling InitOCR() here is redundant, since init() (private one) already calls it, but it's here for completeness
	// NOTE: Assume we're inside "NativeHost" directory, so reference files one directory up via "../"
	errocr := InitOCR("../credentials.json") // NOTE: Init() expects the path to the credentials file to be on SAME directory as the binary wher init() gets called
	if errocr != nil {
		t.Errorf("Failed to initialize OCR: %v", errocr)
	}

	// create image.Image from file
	pngFile := "../ubunchu01_ja/ubunchu01_02.png" // _01.png is the cover page, so we use _02.png where there is real text
	t.Logf("Reading image file: %s", pngFile)
	bytesImage, err := os.ReadFile(pngFile)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s was %d bytes./nDecoding image file: %s", pngFile, len(bytesImage), pngFile)
	theImage, _, err := image.Decode(bytes.NewReader(bytesImage))
	if err != nil {
		t.Error(err)
	}

	// Perform OCR on the image
	textsOfTestImage, err := PerformOCR(theImage)
	if err != nil {
		t.Error(err)
	}

	for _, text := range textsOfTestImage {
		if text == "" {
			t.Error("text is empty")
		}
	}
}
