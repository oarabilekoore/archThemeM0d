package cmd

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

// TonalPalette holds a map of tones for a single color role.
// The key is the tone level (0, 10, 20, ..., 100).
type TonalPalette struct {
	Tones map[int]color.RGBA
}

// ClassifiedTheme is a Material-inspired theme structure.
// It holds full tonal palettes for key roles and specific colors for surfaces.
type ClassifiedTheme struct {
	Primary   TonalPalette
	Secondary TonalPalette
	Tertiary  TonalPalette
	Neutral   TonalPalette

	// Surface colors for dark theme
	Surface          color.RGBA // App backgrounds
	SurfaceVariant   color.RGBA // Cards, dialogs
	OnSurface        color.RGBA // Text on Surface
	OnSurfaceVariant color.RGBA // Secondary text
	PrimaryFixed     color.RGBA // A primary color that doesn't change
	OnPrimaryFixed   color.RGBA // Text on PrimaryFixed
}

// HCT represents a color in Hue, Chroma, Tone space (Material 3's color space)
type HCT struct {
	H float64 // Hue (0-360)
	C float64 // Chroma (0-100+)
	T float64 // Tone (0-100)
}

// colorMetrics is a helper struct for color analysis and sorting.
type colorMetrics struct {
	Color    color.RGBA
	HCT      HCT
	Vibrancy float64 // Combined chroma and tone score
	Index    int
}

var templateFillCmd = &cobra.Command{
	Use:   "build",
	Short: "build - fill the templates provided with theme data for use.",
	Run:   buildTemplates,
}

func init() {
	rootCmd.AddCommand(templateFillCmd)
}

// rgbToHct converts RGB to HCT color space (Material 3's perceptual color space)
func rgbToHct(c color.RGBA) HCT {
	// First convert RGB to XYZ
	r, g, b := float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0

	// Gamma correction
	if r > 0.04045 {
		r = math.Pow((r+0.055)/1.055, 2.4)
	} else {
		r = r / 12.92
	}
	if g > 0.04045 {
		g = math.Pow((g+0.055)/1.055, 2.4)
	} else {
		g = g / 12.92
	}
	if b > 0.04045 {
		b = math.Pow((b+0.055)/1.055, 2.4)
	} else {
		b = b / 12.92
	}

	// Convert to XYZ
	x := r*0.4124564 + g*0.3575761 + b*0.1804375
	y := r*0.2126729 + g*0.7151522 + b*0.0721750
	z := r*0.0193339 + g*0.1191920 + b*0.9503041

	// Convert XYZ to LAB
	xn, yn, zn := 0.95047, 1.0, 1.08883 // D65 illuminant
	fx := labF(x / xn)
	fy := labF(y / yn)
	fz := labF(z / zn)

	L := 116*fy - 16
	A := 500 * (fx - fy)
	B := 200 * (fy - fz)

	// Convert LAB to HCT
	chroma := math.Sqrt(A*A + B*B)
	hue := math.Atan2(B, A) * 180 / math.Pi
	if hue < 0 {
		hue += 360
	}

	return HCT{
		H: hue,
		C: chroma,
		T: L, // In HCT, Tone is equivalent to L* in LAB
	}
}

func labF(t float64) float64 {
	if t > 0.008856 {
		return math.Pow(t, 1.0/3.0)
	}
	return (7.787*t + 16.0/116.0)
}

// hctToRgb converts HCT back to RGB
func hctToRgb(hct HCT) color.RGBA {
	// Convert HCT to LAB
	L := hct.T
	A := hct.C * math.Cos(hct.H*math.Pi/180)
	B := hct.C * math.Sin(hct.H*math.Pi/180)

	// Convert LAB to XYZ
	fy := (L + 16) / 116
	fx := A/500 + fy
	fz := fy - B/200

	var x, y, z float64
	if fx*fx*fx > 0.008856 {
		x = fx * fx * fx
	} else {
		x = (fx - 16.0/116.0) / 7.787
	}
	if fy*fy*fy > 0.008856 {
		y = fy * fy * fy
	} else {
		y = (fy - 16.0/116.0) / 7.787
	}
	if fz*fz*fz > 0.008856 {
		z = fz * fz * fz
	} else {
		z = (fz - 16.0/116.0) / 7.787
	}

	// Scale by illuminant
	x *= 0.95047
	y *= 1.0
	z *= 1.08883

	// Convert XYZ to RGB
	r := x*3.2404542 + y*-1.5371385 + z*-0.4985314
	g := x*-0.9692660 + y*1.8760108 + z*0.0415560
	b := x*0.0556434 + y*-0.2040259 + z*1.0572252

	// Gamma correction
	if r > 0.0031308 {
		r = 1.055*math.Pow(r, 1/2.4) - 0.055
	} else {
		r = 12.92 * r
	}
	if g > 0.0031308 {
		g = 1.055*math.Pow(g, 1/2.4) - 0.055
	} else {
		g = 12.92 * g
	}
	if b > 0.0031308 {
		b = 1.055*math.Pow(b, 1/2.4) - 0.055
	} else {
		b = 12.92 * b
	}

	// Clamp to [0,1] and convert to 0-255
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))

	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: 255,
	}
}

// calculateVibrancy calculates Material 3 style vibrancy score
func calculateVibrancy(hct HCT) float64 {
	// Material 3 considers both chroma and tone for vibrancy
	// Higher chroma = more vibrant, tone around 50 = most vibrant
	chromaWeight := hct.C / 100.0               // Normalize chroma
	toneWeight := 1.0 - math.Abs(hct.T-50)/50.0 // Peak at tone 50
	return chromaWeight*0.7 + toneWeight*0.3
}

// calculateHueDistance calculates the shortest distance between two hues
func calculateHueDistance(h1, h2 float64) float64 {
	diff := math.Abs(h1 - h2)
	if diff > 180 {
		diff = 360 - diff
	}
	return diff
}

// isHarmoniousHue checks if two hues are harmoniously related
func isHarmoniousHue(primary, secondary float64) bool {
	distance := calculateHueDistance(primary, secondary)
	// Harmonious relationships: complementary (~180°), triadic (~120°), analogous (~30°)
	return (distance >= 25 && distance <= 35) || // Analogous
		(distance >= 115 && distance <= 125) || // Triadic
		(distance >= 175 && distance <= 185) // Complementary
}

// generateTonalPaletteHct creates a Material 3 compliant tonal palette using HCT
func generateTonalPaletteHct(seedHct HCT) TonalPalette {
	tones := make(map[int]color.RGBA)

	// Material 3 tone levels
	toneLevels := []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 95, 99, 100}

	for _, tone := range toneLevels {
		// Keep hue and chroma, only change tone
		newHct := HCT{
			H: seedHct.H,
			C: seedHct.C,
			T: float64(tone),
		}

		// For very dark/light tones, reduce chroma to avoid impossible colors
		if tone <= 20 || tone >= 90 {
			newHct.C = seedHct.C * 0.8
		}
		if tone <= 10 || tone >= 95 {
			newHct.C = seedHct.C * 0.5
		}

		tones[tone] = hctToRgb(newHct)
	}

	return TonalPalette{Tones: tones}
}

// classifyPaletteMaterial3 analyzes a raw color palette using Material 3 principles
func classifyPaletteMaterial3(palette []color.RGBA) ClassifiedTheme {
	if len(palette) < 4 {
		log.Fatal("ERROR: Palette must have at least 4 colors for Material 3 generation.")
	}

	// Convert all colors to HCT and calculate metrics
	metrics := make([]colorMetrics, len(palette))
	for i, c := range palette {
		hct := rgbToHct(c)
		vibrancy := calculateVibrancy(hct)
		metrics[i] = colorMetrics{
			Color:    c,
			HCT:      hct,
			Vibrancy: vibrancy,
			Index:    i,
		}
	}

	// Sort by vibrancy (Material 3's approach)
	sort.SliceStable(metrics, func(i, j int) bool {
		return metrics[i].Vibrancy > metrics[j].Vibrancy
	})

	// Select Primary as the most vibrant color
	primarySeed := metrics[0]

	// Find Secondary with harmonious hue relationship to Primary
	var secondarySeed colorMetrics
	found := false
	for i := 1; i < len(metrics); i++ {
		if isHarmoniousHue(primarySeed.HCT.H, metrics[i].HCT.H) {
			secondarySeed = metrics[i]
			found = true
			break
		}
	}
	if !found {
		// Fallback to second most vibrant if no harmonious hue found
		secondarySeed = metrics[1]
	}

	// Find Tertiary that's different from both Primary and Secondary
	var tertiarySeed colorMetrics
	for i := 1; i < len(metrics); i++ {
		if metrics[i].Index != secondarySeed.Index {
			hue := metrics[i].HCT.H
			primaryDist := calculateHueDistance(primarySeed.HCT.H, hue)
			secondaryDist := calculateHueDistance(secondarySeed.HCT.H, hue)

			// Ensure it's sufficiently different from both Primary and Secondary
			if primaryDist > 60 && secondaryDist > 60 {
				tertiarySeed = metrics[i]
				found = true
				break
			}
		}
	}
	if !found {
		tertiarySeed = metrics[2] // Fallback
	}

	// Find a low-chroma color for Neutral
	var neutralSeed colorMetrics
	sort.SliceStable(metrics, func(i, j int) bool {
		return metrics[i].HCT.C < metrics[j].HCT.C // Sort by chroma (ascending)
	})

	usedIndices := map[int]bool{
		primarySeed.Index:   true,
		secondarySeed.Index: true,
		tertiarySeed.Index:  true,
	}

	for _, m := range metrics {
		if !usedIndices[m.Index] {
			neutralSeed = m
			break
		}
	}

	// Generate tonal palettes using HCT color space
	primaryPalette := generateTonalPaletteHct(primarySeed.HCT)
	secondaryPalette := generateTonalPaletteHct(secondarySeed.HCT)
	tertiaryPalette := generateTonalPaletteHct(tertiarySeed.HCT)
	neutralPalette := generateTonalPaletteHct(neutralSeed.HCT)

	// Assemble the final theme based on Material 3 dark theme specifications
	return ClassifiedTheme{
		Primary:   primaryPalette,
		Secondary: secondaryPalette,
		Tertiary:  tertiaryPalette,
		Neutral:   neutralPalette,

		// Material 3 dark surface colors
		Surface:          neutralPalette.Tones[6],  // Very dark neutral
		SurfaceVariant:   neutralPalette.Tones[30], // Slightly lighter
		OnSurface:        neutralPalette.Tones[90], // Light text on dark surface
		OnSurfaceVariant: neutralPalette.Tones[80], // Secondary text
		PrimaryFixed:     primaryPalette.Tones[90], // Fixed primary for consistency
		OnPrimaryFixed:   primaryPalette.Tones[10], // Text on fixed primary
	}
}

func buildTemplates(cmd *cobra.Command, args []string) {
	themeFilePath := filepath.Join(homeDir, themeFileDir)

	if _, err := os.Stat(themeFilePath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("\nISSUE: Theme file not found.")
			fmt.Println("FIX: Running the 'generate' command first...")
			GenerateThemeFile(cmd, args)
			fmt.Println("---")
		} else {
			fmt.Printf("\nERROR: Could not stat theme file: %v\n", err)
			return
		}
	}

	data, err := os.ReadFile(themeFilePath)
	if err != nil {
		fmt.Printf("\nERROR: Could not read theme file: %v\n", err)
		return
	}

	var allMonitorsData []MonitorInfo
	if err = json.Unmarshal(data, &allMonitorsData); err != nil {
		fmt.Printf("\nERROR: Could not unmarshal theme JSON: %v\n", err)
		return
	}

	appDir := filepath.Join(homeDir, tm0dDir)
	templatesDir := filepath.Join(appDir, "Templates")
	themesOutputDir := filepath.Join(appDir, "Themes")

	templateFiles, err := os.ReadDir(templatesDir)
	if err != nil {
		fmt.Printf("ERROR: Failed to read templates directory '%s': %v\n", templatesDir, err)
		return
	}

	fmt.Printf("Preparing output directory: %s\n", themesOutputDir)
	_ = os.RemoveAll(themesOutputDir)
	if err := os.MkdirAll(themesOutputDir, 0755); err != nil {
		fmt.Printf("ERROR: Could not create output directory: %v\n", err)
		return
	}

	funcMap := template.FuncMap{
		"toHex": func(c color.RGBA) string {
			return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
		},
		"toRgba": func(c color.RGBA, alpha string) string {
			return fmt.Sprintf("rgba(%d, %d, %d, %s)", c.R, c.G, c.B, alpha)
		},
		// Helper to easily access a tone from a palette in the template.
		"tone": func(p TonalPalette, level int) color.RGBA {
			if c, ok := p.Tones[level]; ok {
				return c
			}
			// Return a bright pink for debugging if a tone is missing.
			return color.RGBA{R: 255, G: 0, B: 255, A: 255}
		},
	}

	for _, monitorData := range allMonitorsData {
		monitorOutputDir := filepath.Join(themesOutputDir, monitorData.Monitor)
		if err := os.MkdirAll(monitorOutputDir, 0755); err != nil {
			log.Printf("ERROR: Could not create directory for monitor %s: %v", monitorData.Monitor, err)
			continue
		}
		fmt.Printf("\nProcessing templates for monitor: %s\n", monitorData.Monitor)

		// Use Material 3 classification instead of simple saturation sorting
		classifiedTheme := classifyPaletteMaterial3(monitorData.Theme.Palletes)

		templateData := struct {
			Monitor string
			Theme   ClassifiedTheme
		}{
			Monitor: monitorData.Monitor,
			Theme:   classifiedTheme,
		}

		for _, file := range templateFiles {
			if file.IsDir() {
				continue
			}

			templateName := file.Name()
			finalFileName := strings.TrimSuffix(templateName, ".tmpl")
			templatePath := filepath.Join(templatesDir, templateName)
			outputPath := filepath.Join(monitorOutputDir, finalFileName)

			fmt.Printf("  -> Rendering %s\n", templateName)

			templateContent, err := os.ReadFile(templatePath)
			if err != nil {
				log.Printf("ERROR: Failed to read template file %s: %v", templateName, err)
				continue
			}

			tmpl, err := template.New(templateName).Funcs(funcMap).Parse(string(templateContent))
			if err != nil {
				log.Printf("ERROR: Failed to parse template %s: %v", templateName, err)
				continue
			}

			outputFile, err := os.Create(outputPath)
			if err != nil {
				log.Printf("ERROR: Failed to create output file %s: %v", outputPath, err)
				continue
			}

			err = tmpl.Execute(outputFile, templateData)
			if err != nil {
				log.Printf("ERROR: Failed to execute template %s: %v", templateName, err)
			}
			outputFile.Close()
		}
	}
	fmt.Println("\nBuild complete!")
}
