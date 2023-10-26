package main

import (
	"bufio"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/skip2/go-qrcode"
)

func main() {
	// Prompt the user to enter the path to the input folder
	fmt.Print("Enter the path to the input folder: ")
	var inputPath string
	fmt.Scan(&inputPath)

	// Prompt the user to enter the path to the output folder
	fmt.Print("Enter the path to the output folder: ")
	var outputPath string
	fmt.Scan(&outputPath)

	// Ensure that the output folder exists, create it if it doesn't
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		log.Fatalf("Error creating output folder: %v", err)
	}

	// Prompt the user to enter the file name in the input folder
	fmt.Print("Enter the name of the input text file in the input folder: ")
	var inputFileName string
	fmt.Scan(&inputFileName)

	// Construct the full path to the input file
	inputFilePath := filepath.Join(inputPath, inputFileName)

	// Open and read the input text file
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("Error opening input file: %v", err)
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	for lineNum := 1; scanner.Scan(); lineNum++ {
		data := scanner.Text()

		// Remove any leading or trailing whitespace
		data = strings.TrimSpace(data)

		if data == "" {
			log.Printf("Skipping empty line at line %d.", lineNum)
			continue
		}

		qr, err := qrcode.New(data, qrcode.Medium)
		if err != nil {
			log.Printf("Error creating QR code for data at line %d: %v", lineNum, err)
			continue
		}
		qr.ForegroundColor = color.RGBA{0, 0, 0, 255}
		qr.BackgroundColor = color.RGBA{255, 255, 255, 255}

		img := qr.Image(256)

		// Generate a file name for the QR code image based on the line number
		outputFileName := fmt.Sprintf("qrcode_line_%d.png", lineNum)

		// Save the QR code image to the output folder
		outputFilePath := filepath.Join(outputPath, outputFileName)
		file, err := os.Create(outputFilePath)
		if err != nil {
			log.Printf("Error saving QR code for data at line %d: %v", lineNum, err)
			continue
		}
		defer file.Close()

		err = png.Encode(file, img)
		if err != nil {
			log.Printf("Error encoding QR code for data at line %d: %v", lineNum, err)
		}

		fmt.Printf("QR code for data at line %d saved as %s\n", lineNum, outputFilePath)
	}
}
