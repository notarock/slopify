package video

import (
	"fmt"

	"gopkg.in/gographics/imagick.v3/imagick"
)

func CreateTitleCard(title, filePath string, videoWidth int) string {
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
	// Set text color
	// Set text color
	textColor := imagick.NewPixelWand()
	defer textColor.Destroy()
	textColor.SetColor("rgba(0, 0, 0, 1)")
	draw.SetFillColor(textColor)

	// Split text into lines
	lines := splitTextIntoLines(title, 42) // Adjust the line length as needed

	// Draw each line of text
	lineHeight := 64 // Adjust line height as needed
	for i, line := range lines {
		yPos := float64(lineHeight*(i+1) + 10)
		draw.Annotation(30, yPos, line)
		// Draw to image
		mw.DrawImage(draw)
	}
	// Write the image to a file
	mw.WriteImage(filePath)

	fmt.Println("Image created successfully")

	resized := resizeImage(filePath, videoWidth)
	return resized
}

func resizeImage(filePath string, videoWidth int) string {
	// Create new ImageMagick wand
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// Read the image file
	err := mw.ReadImage(filePath)
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return ""
	}

	// Get the current dimensions of the image
	width := mw.GetImageWidth()
	height := mw.GetImageHeight()

	// Calculate the new dimensions based on the video width
	newWidth := videoWidth
	newHeight := int(float64(height) / float64(width) * float64(newWidth))

	// Resize the image
	fmt.Println("Resizing image to", newWidth, "x", newHeight)
	err = mw.ResizeImage(uint(newWidth), uint(newHeight), imagick.FILTER_LANCZOS)
	if err != nil {
		fmt.Println("Error resizing image:", err)
		return ""
	}

	// Write the resized image to a new file
	newFilePath := filePath[:len(filePath)-4] + "-resized.png"
	fmt.Println("Writing resized image to", newFilePath)
	err = mw.WriteImage(newFilePath)
	fmt.Println("Image resized successfully")
	if err != nil {
		fmt.Println("Error writing resized image:", err)
		return ""
	}

	return newFilePath
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
