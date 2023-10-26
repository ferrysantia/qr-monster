package main

import (
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"

	"github.com/skip2/go-qrcode"
)

func main() {
	// Prompt the user to enter the data for the QR code
	fmt.Print("Enter the data for the QR code: ")
	var data string
	fmt.Scan(&data)

	// Prompt the user to enter the path for the output folder
	fmt.Print("Enter the path for the output folder: ")
	var outputPath string
	fmt.Scan(&outputPath)

	// Ensure that the output folder exists, create it if it doesn't
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		log.Fatalf("Error creating output folder: %v", err)
	}

	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}
	qr.ForegroundColor = color.RGBA{0, 0, 0, 255}
	qr.BackgroundColor = color.RGBA{255, 255, 255, 255}

	img := qr.Image(256)

	// Save the QR code image to the specified output folder
	outputFilePath := outputPath + "/qrcode.png"
	file, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("QR code saved as %s\n", outputFilePath)
}
