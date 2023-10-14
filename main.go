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
	// String to encode in the QR code
	data := "some_string"

	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}
	qr.ForegroundColor = color.RGBA{0, 0, 0, 255}
	qr.BackgroundColor = color.RGBA{255, 255, 255, 255}

	img := qr.Image(256)

	// Save the QR code image to a file
	file, err := os.Create("qrcode.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("QR code saved as qrcode.png")
}
