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

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
)

var outputFolder string // Define outputFolder as a global variable

func main() {
	app := fiber.New()

	// Serve the HTML form
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("index.html")
	})

	// Handle form submission and QR code generation
	app.Post("/generate", func(c *fiber.Ctx) error {
		inputFolder := c.FormValue("inputFolder")
		outputFolder = c.FormValue("outputFolder") // Update the global outputFolder
		inputFileName := c.FormValue("inputFileName")

		// Ensure that the output folder exists, create it if it doesn't
		if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
			log.Printf("Error creating output folder: %v", err)
			return c.SendString("Error creating output folder")
		}

		// Construct the full path to the input file
		inputFilePath := filepath.Join(inputFolder, inputFileName)

		// Open and read the input text file
		inputFile, err := os.Open(inputFilePath)
		if err != nil {
			log.Printf("Error opening input file: %v", err)
			return c.SendString("Error opening input file")
		}
		defer inputFile.Close()

		scanner := bufio.NewScanner(inputFile)

		// Create a list to store the generated QR code file paths
		qrCodeFilePaths := []string{}

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
			outputFilePath := filepath.Join(outputFolder, outputFileName)
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

			qrCodeFilePaths = append(qrCodeFilePaths, outputFilePath)
		}

		// Respond with a JSON array containing the generated file names and paths
		return c.JSON(qrCodeFilePaths)
	})

	// Serve the QR code files for download
	app.Get("/download", func(c *fiber.Ctx) error {
		fileName := c.Query("file")
		filePath := filepath.Join(outputFolder, fileName)

		return c.SendFile(filePath)
	})

	// Start the Fiber web server
	log.Fatal(app.Listen(":3000"))
}
