package main

import (
	"fmt"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func main() {
	// Initialize imagick
	imagick.Initialize()
	defer imagick.Terminate()

	// Create new ImageMagick wand
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Set image size
	transparent := imagick.NewPixelWand()
	defer transparent.Destroy()
	transparent.SetColor("transparent")
	mw.NewImage(1200, 300, transparent)

	// Create a gradient fill for the rectangle
	innerTriangle := imagick.NewPixelWand()
	defer innerTriangle.Destroy()
	innerTriangle.SetColor("rgba(235, 235, 235, 1)") // Semi-transparent black

	// Draw a semi-transparent rectangle behind the text
	draw := imagick.NewDrawingWand()
	defer draw.Destroy()
	draw.SetStrokeOpacity(0)
	draw.SetFillColor(innerTriangle)
	draw.RoundRectangle(10, 10, 1200, 300, 50, 50) // Adjust rectangle size and position as needed

	// Set font properties
	draw.SetFont("DejaVu-Sans")
	draw.SetFontSize(56)

	// Define text
	text := "What is currently in it's \"Golden age\", but not enough people know about it?"

	// Set text color
	// Set text color
	textColor := imagick.NewPixelWand()
	defer textColor.Destroy()
	textColor.SetColor("rgba(0, 0, 0, 1)")
	draw.SetFillColor(textColor)

	// Split text into lines
	lines := splitTextIntoLines(text, 42) // Adjust the line length as needed

	// Draw each line of text
	lineHeight := 64 // Adjust line height as needed
	for i, line := range lines {
		yPos := float64(lineHeight*(i+1) + 10)
		fmt.Println(yPos)
		fmt.Println(line)
		draw.Annotation(30, yPos, line)
		// Draw to image
		mw.DrawImage(draw)
	}
	// Write the image to a file
	mw.WriteImage("title_image.png")

	fmt.Println("Image created successfully")

}

func splitTextIntoLines(text string, maxLineLength int) []string {
	var lines []string
	for len(text) > maxLineLength {
		splitIndex := maxLineLength
		for splitIndex > 0 && text[splitIndex] != ' ' {
			splitIndex--
		}
		lines = append(lines, text[:splitIndex])
		text = text[splitIndex+1:]
	}
	lines = append(lines, text)
	return lines
}
