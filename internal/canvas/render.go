// Package canvas handles canvas rendering operations.
package canvas

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"

	"github.com/ergodic/moltcities/internal/models"
)

// Render generates a PNG image of the canvas.
func Render(pixels map[[2]int]string, w io.Writer) error {
	img := image.NewRGBA(image.Rect(0, 0, models.CanvasSize, models.CanvasSize))

	// Direct pixel buffer manipulation for maximum performance
	pix := img.Pix

	// Fill with white (255, 255, 255, 255)
	for j := range pix {
		pix[j] = 255
	}

	// Apply edited pixels directly to byte slice
	stride := img.Stride
	for coord, hex := range pixels {
		x, y := coord[0], coord[1]
		if x < 0 || x >= models.CanvasSize || y < 0 || y >= models.CanvasSize {
			continue
		}
		c, err := HexToColor(hex)
		if err != nil {
			continue // Skip invalid colors
		}
		offset := (y * stride) + x*4
		pix[offset] = c.R
		pix[offset+1] = c.G
		pix[offset+2] = c.B
		pix[offset+3] = 255
	}

	return png.Encode(w, img)
}

// RenderRegion generates a PNG image of a canvas region.
func RenderRegion(pixels [][]string, w io.Writer) error {
	height := len(pixels)
	if height == 0 {
		return fmt.Errorf("empty region")
	}
	width := len(pixels[0])

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	stride := img.Stride
	pix := img.Pix

	for y, row := range pixels {
		for x, hex := range row {
			c, err := HexToColor(hex)
			if err != nil {
				c = color.RGBA{255, 255, 255, 255} // Default to white
			}
			offset := (y * stride) + x*4
			pix[offset] = c.R
			pix[offset+1] = c.G
			pix[offset+2] = c.B
			pix[offset+3] = 255
		}
	}

	return png.Encode(w, img)
}

// HexToColor converts a hex color string to color.RGBA.
// Accepts formats: "#RRGGBB" or "RRGGBB"
// Optimized manual hex parsing instead of fmt.Sscanf for ~30x speedup.
func HexToColor(hex string) (color.RGBA, error) {
	// Remove # prefix if present
	if len(hex) > 0 && hex[0] == '#' {
		hex = hex[1:]
	}

	if len(hex) != 6 {
		return color.RGBA{}, fmt.Errorf("invalid hex color: %s", hex)
	}

	var r, g, b uint8
	for i := 0; i < 6; i++ {
		c := hex[i]
		var val byte
		if c >= '0' && c <= '9' {
			val = c - '0'
		} else if c >= 'a' && c <= 'f' {
			val = c - 'a' + 10
		} else if c >= 'A' && c <= 'F' {
			val = c - 'A' + 10
		} else {
			return color.RGBA{}, fmt.Errorf("invalid hex color: %s", hex)
		}
		switch i {
		case 0:
			r = val << 4
		case 1:
			r |= val
		case 2:
			g = val << 4
		case 3:
			g |= val
		case 4:
			b = val << 4
		case 5:
			b |= val
		}
	}

	return color.RGBA{r, g, b, 255}, nil
}

// ValidateColor checks if a color string is valid.
func ValidateColor(hex string) error {
	_, err := HexToColor(hex)
	return err
}

// ValidateCoordinate checks if a coordinate is within canvas bounds.
func ValidateCoordinate(n int) error {
	if n < 0 || n >= models.CanvasSize {
		return fmt.Errorf("coordinate must be between 0 and %d", models.CanvasSize-1)
	}
	return nil
}

// ValidateRegion checks if a region is valid.
func ValidateRegion(x, y, width, height int) error {
	if err := ValidateCoordinate(x); err != nil {
		return fmt.Errorf("x: %w", err)
	}
	if err := ValidateCoordinate(y); err != nil {
		return fmt.Errorf("y: %w", err)
	}
	if width < 1 || width > models.MaxRegionSize {
		return fmt.Errorf("width must be between 1 and %d", models.MaxRegionSize)
	}
	if height < 1 || height > models.MaxRegionSize {
		return fmt.Errorf("height must be between 1 and %d", models.MaxRegionSize)
	}
	if x+width > models.CanvasSize {
		return fmt.Errorf("region extends beyond canvas width")
	}
	if y+height > models.CanvasSize {
		return fmt.Errorf("region extends beyond canvas height")
	}
	return nil
}
