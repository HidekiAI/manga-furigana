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

	// import google cloud vision package
	vision "cloud.google.com/go/vision/apiv1"
	"golang.org/x/oauth2"

	// google-auth-library-go package
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	"google.golang.org/grpc/metadata"
)

var visionClient *vision.ImageAnnotatorClient
var bgContext context.Context

// Mainly, just a proxy to the performOCR() function in which it'll encode/decode JSON from the JavaScript code (browser extension)
func PerformOCRProxy(imageFromMessage image.Image, enc *json.Encoder) (string, error) {
	// Resize the image to a smaller size to speed up processing
	scaledAndGreyScaledImage := imaging.Resize(imageFromMessage, 800, 0, imaging.Lanczos)

	// Convert the image to grayscale
	scaledAndGreyScaledImage = imaging.Grayscale(scaledAndGreyScaledImage)

	// other suggestions includes adjusting the contrast, binarizing, deskewing, noise reduction, and edge detection
	scaledAndGreyScaledImage = imaging.AdjustContrast(scaledAndGreyScaledImage, 20) // as a starter, adjust the contrast to +20%

	// Process the image
	tokenizedText, err := PerformOCR(scaledAndGreyScaledImage)
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
func PerformOCR(img image.Image) ([]string, error) {
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

func createVisionClient(token string) (*vision.ImageAnnotatorClient, error) {
	// Create a context and set the token as a header
	ctx := context.Background()
	header := fmt.Sprintf("Bearer %s", token)
	ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", header)

	// Create an oauth2.TokenSource with the token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// Create an option.ClientOption with the oauth2.TokenSource
	clientOption := option.WithTokenSource(ts)

	// Create a transport with the clientOption
	httpClient, _, err := transport.NewHTTPClient(ctx, clientOption)
	if err != nil {
		return nil, err
	}

	// Create a vision.ImageAnnotatorClient with the transport
	visionClient, err := vision.NewImageAnnotatorClient(ctx, clientOption, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return visionClient, nil
}

// NOTE: Because this is client-to-server, we need to use OAuth2 to authenticate the client
// rather than using a service account key file (which is used for server-to-server)
func InitOCR() error {
	// Use the chrome.identity.getAuthToken() function to obtain an OAuth2 access token (this is a client-to-server)
	// based chrome extension, so we need to use OAuth2 to authenticate the client
	return nil
}
