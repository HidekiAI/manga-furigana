package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	// import google cloud vision package
	// import kagome v2 tokenizer and dictionary package
)

func main() {

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

func init() {
	if getPlatform() == "windows" {
		// Register the native messaging host
		err := exec.Command("regedit", "/s", "./manga-furigana.codemonkeyninja.dev.json").Run()
		if err != nil {
			log.Fatal("Failed to register the native messaging host:", err)
		}
	}

	errbg := initBackground()
	if errbg != nil {
		fmt.Println("Failed to initialize background:", errbg.Error())
	}
	errtkn := initTokenizer()
	if errtkn != nil {
		fmt.Println("Failed to initialize tokenizer:", errtkn.Error())
	}
	errocr := initOCR("./credentials.json")
	if errocr != nil {
		fmt.Println("Failed to initialize OCR:", errocr.Error())
	}
}
