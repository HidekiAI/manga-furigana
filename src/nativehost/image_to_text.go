package main // too small of a project to have this in separate directory
import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"
	"strings"

	visionpb "cloud.google.com/go/vision/v2/apiv1/visionpb"
	"github.com/disintegration/imaging"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"

	// import google cloud vision package
	vision "cloud.google.com/go/vision/apiv1"
)

var visionClient *vision.ImageAnnotatorClient
var bgContext context.Context

func PerformOCR(imageFromMessage image.Image, enc *json.Encoder) (string, error) {
	// Resize the image to a smaller size to speed up processing
	scaledAndGreyScaledImage := imaging.Resize(imageFromMessage, 800, 0, imaging.Lanczos)

	// Convert the image to grayscale
	scaledAndGreyScaledImage = imaging.Grayscale(scaledAndGreyScaledImage)

	// other suggestions includes adjusting the contrast, binarizing, deskewing, noise reduction, and edge detection
	scaledAndGreyScaledImage = imaging.AdjustContrast(scaledAndGreyScaledImage, 20) // as a starter, adjust the contrast to +20%

	// Process the image
	tokenizedText, err := performOCR(scaledAndGreyScaledImage, enc)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "", err
	}

	// Join the slice of strings into a single string
	joinedText := strings.Join(tokenizedText, " ")

	// Send the tokenized image back to the JavaScript code
	if err := enc.Encode(TokenizedImage{URL: "", TokenizedText: joinedText}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return joinedText, nil
}

// Implement the logic to call Google Cloud Vision OCR API and extract the Japanese text from the image URL
func performOCR(img image.Image, enc *json.Encoder) ([]string, error) {
	// let's make sure initOCR() was called
	if visionClient == nil || bgContext == nil {
		// fail gracefully, but inform the caller that initOCR() needs to be called in the error message
		return []string{}, fmt.Errorf("visionClient or bgContext is nil, did you call initOCR()?")
	}

	// Convert the image to format acceptable by Google Cloud Vision API, Convert the image to a []byte slice
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return []string{}, fmt.Errorf("failed to encode image: %v", err)
	}
	imageBytes := buf.Bytes()

	// Create a new image annotation request
	image, err := vision.NewImageFromReader(bytes.NewReader(imageBytes))
	if err != nil {
		return []string{}, fmt.Errorf("failed to create image object: %v", err)
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
	response, err := visionClient.AnnotateImage(bgContext, request)
	if err != nil {
		return []string{}, fmt.Errorf("failed to perform image annotation: %v", err)
	}

	// create array of strings to store the extracted text
	texts := []string{}
	for _, annotation := range response.TextAnnotations {
		texts = append(texts, annotation.Description)
	}

	return texts, nil
}

func initOCR(pathToCredentials_json string) error {
	// Initialize the Google Cloud Vision API client
	// Set up a context and create a new client
	bgContext = context.Background()
	// Create a new Google Cloud Vision client
	creds, err := transport.Creds(bgContext, option.WithCredentialsFile(pathToCredentials_json))
	if err != nil {
		fmt.Printf("failed to load credentials for '%s': %v\n", pathToCredentials_json, err)
		os.Exit(1)
	}
	visionClient, err = vision.NewImageAnnotatorClient(bgContext, option.WithCredentials(creds))
	if err != nil {
		fmt.Println("failed to create Vision client:", err)
		return fmt.Errorf("failed to create Vision client: %v", err)
	}
	return nil
}
