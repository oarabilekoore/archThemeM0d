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

	// (Dark Theme context)
	Surface          color.RGBA // App backgrounds
	SurfaceVariant   color.RGBA // Cards, dialogs
	OnSurface        color.RGBA // Text on Surface
	OnSurfaceVariant color.RGBA // Secondary text
	PrimaryFixed     color.RGBA // A primary color that doesn't change
	OnPrimaryFixed   color.RGBA // Text on PrimaryFixed
}

// colorMetrics is a helper struct for color analysis and sorting.
type colorMetrics struct {
	Color      color.RGBA
	Luminance  float64
	Saturation float64
	Index      int
}

var templateFillCmd = &cobra.Command{
	Use:   "build",
	Short: "build - fill the templates provided with theme data for use.",
	Run:   buildTemplates,
}

func init() {
	rootCmd.AddCommand(templateFillCmd)
}

// rgbToHsl converts an RGB color to HSL values.
// We use this to determine saturation and luminance for classification.
func rgbToHsl(c color.RGBA) (h, s, l float64) {
	r, g, b := float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0
	max, min := math.Max(r, math.Max(g, b)), math.Min(r, math.Min(g, b))
	l = (max + min) / 2.0
	if max == min {
		return 0, 0, l // Achromatic (gray)
	}
	d := max - min
	if l > 0.5 {
		s = d / (2.0 - max - min)
	} else {
		s = d / (max + min)
	}
	switch max {
	case r:
		h = (g - b) / d
		if g < b {
			h += 6
		}
	case g:
		h = (b-r)/d + 2
	case b:
		h = (r-g)/d + 4
	}
	h /= 6
	return h, s, l
}

// blendColor mixes two colors together based on a ratio.
func blendColor(c1, c2 color.RGBA, ratio float64) color.RGBA {
	ratio = math.Max(0, math.Min(1, ratio)) // Clamp ratio between 0 and 1
	return color.RGBA{
		R: uint8(float64(c1.R)*(1.0-ratio) + float64(c2.R)*ratio),
		G: uint8(float64(c1.G)*(1.0-ratio) + float64(c2.G)*ratio),
		B: uint8(float64(c1.B)*(1.0-ratio) + float64(c2.B)*ratio),
		A: 255,
	}
}

// generateTonalPalette creates a full 13-step tonal ramp from a single seed color.
func generateTonalPalette(seed color.RGBA) TonalPalette {
	tones := make(map[int]color.RGBA)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{R: 0, G: 0, B: 0, A: 255}

	// For a dark theme, the seed color is considered Tone 80.
	tones[80] = seed
	tones[100] = white
	tones[0] = black

	// Generate lighter tones by blending the seed with white.
	tones[90] = blendColor(seed, white, 0.50)
	tones[95] = blendColor(seed, white, 0.75)
	tones[99] = blendColor(seed, white, 0.95)

	// Generate darker tones by blending the seed with black.
	tones[70] = blendColor(seed, black, 0.15)
	tones[60] = blendColor(seed, black, 0.30)
	tones[50] = blendColor(seed, black, 0.45)
	tones[40] = blendColor(seed, black, 0.60)
	tones[30] = blendColor(seed, black, 0.75)
	tones[20] = blendColor(seed, black, 0.85)
	tones[10] = blendColor(seed, black, 0.95)

	return TonalPalette{Tones: tones}
}

// classifyPalette analyzes a raw color palette and generates a full Material-style theme.
func classifyPalette(palette []color.RGBA) ClassifiedTheme {
	if len(palette) < 4 {
		log.Fatal("ERROR: Palette must have at least 4 colors for Material generation.")
	}

	metrics := make([]colorMetrics, len(palette))
	for i, c := range palette {
		_, s, l := rgbToHsl(c)
		metrics[i] = colorMetrics{Color: c, Luminance: l, Saturation: s, Index: i}
	}

	// Sort by saturation to find the most vibrant colors for key roles.
	sort.SliceStable(metrics, func(i, j int) bool {
		return metrics[i].Saturation > metrics[j].Saturation
	})

	// Assign seed colors for Primary, Secondary, and Tertiary roles.
	primarySeed := metrics[0].Color
	secondarySeed := metrics[1].Color
	tertiarySeed := metrics[2].Color
	usedIndices := map[int]bool{
		metrics[0].Index: true,
		metrics[1].Index: true,
		metrics[2].Index: true,
	}

	// Find a low-saturation color for the Neutral seed from the remaining colors.
	var neutralSeed color.RGBA
	sort.SliceStable(metrics, func(i, j int) bool {
		return metrics[i].Saturation < metrics[j].Saturation
	})
	for _, m := range metrics {
		if !usedIndices[m.Index] {
			neutralSeed = m.Color
			break
		}
	}

	// Generate the full tonal palettes from our chosen seed colors.
	primaryPalette := generateTonalPalette(primarySeed)
	secondaryPalette := generateTonalPalette(secondarySeed)
	tertiaryPalette := generateTonalPalette(tertiarySeed)
	neutralPalette := generateTonalPalette(neutralSeed)

	// Assemble the final theme based on Material 3 dark theme conventions.
	return ClassifiedTheme{
		Primary:   primaryPalette,
		Secondary: secondaryPalette,
		Tertiary:  tertiaryPalette,
		Neutral:   neutralPalette,

		Surface:          neutralPalette.Tones[10],
		SurfaceVariant:   neutralPalette.Tones[30],
		OnSurface:        primaryPalette.Tones[90],
		OnSurfaceVariant: neutralPalette.Tones[80],
		PrimaryFixed:     primaryPalette.Tones[90],
		OnPrimaryFixed:   primaryPalette.Tones[10],
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

		// Classify the raw palette into our new structured theme.
		classifiedTheme := classifyPalette(monitorData.Theme.Palletes)

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
