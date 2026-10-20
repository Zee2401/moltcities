package canvas

import (
	"bytes"
	"image/color"
	"io"
	"testing"
)

func TestHexToColor(t *testing.T) {
	tests := []struct {
		input    string
		expected color.RGBA
		wantErr  bool
	}{
		{"#FF5733", color.RGBA{255, 87, 51, 255}, false},
		{"0080FF", color.RGBA{0, 128, 255, 255}, false},
		{"#000000", color.RGBA{0, 0, 0, 255}, false},
		{"FFFFFF", color.RGBA{255, 255, 255, 255}, false},
		{"invalid", color.RGBA{}, true},
		{"#12345", color.RGBA{}, true},
		{"#1234567", color.RGBA{}, true},
		{"#GGGGGG", color.RGBA{}, true},
	}

	for _, tt := range tests {
		got, err := HexToColor(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("HexToColor(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.expected {
			t.Errorf("HexToColor(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestRender(t *testing.T) {
	pixels := map[[2]int]string{
		{0, 0}: "#FF0000",
		{1, 1}: "#00FF00",
	}

	var buf bytes.Buffer
	err := Render(pixels, &buf)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Errorf("Expected non-empty output PNG buffer")
	}
}

func TestRenderRegion(t *testing.T) {
	pixels := [][]string{
		{"#FF0000", "#00FF00"},
		{"#0000FF", "#FFFFFF"},
	}

	var buf bytes.Buffer
	err := RenderRegion(pixels, &buf)
	if err != nil {
		t.Fatalf("RenderRegion failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Errorf("Expected non-empty output PNG buffer")
	}
}

func BenchmarkHexToColor(b *testing.B) {
	hexes := []string{"#FF5733", "0080FF", "#123456", "ABCDEF"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HexToColor(hexes[i%len(hexes)])
	}
}

func BenchmarkRender(b *testing.B) {
	pixels := make(map[[2]int]string)
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			pixels[[2]int{x, y}] = "#FF5733"
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Render(pixels, io.Discard)
	}
}
