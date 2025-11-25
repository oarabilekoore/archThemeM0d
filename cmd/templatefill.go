package cmd

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var templateFillCmd = &cobra.Command{
	Use:   "build",
	Short: "build - fill the templates provided with theme data for use.",
	Run:   BuildTemplates,
}

func init() {
	rootCmd.AddCommand(templateFillCmd)
}

func BuildTemplates(cmd *cobra.Command, args []string) {
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
		"toRGB": func(c color.RGBA) string {
			return fmt.Sprintf("[%d, %d, %d]", c.R, c.G, c.B)
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
		"replace": func(s, old, new string) string {
			return strings.Replace(s, old, new, -1)
		},
	}

	for _, monitorData := range allMonitorsData {
		monitorOutputDir := filepath.Join(themesOutputDir, monitorData.Monitor)
		if err := os.MkdirAll(monitorOutputDir, 0755); err != nil {
			log.Printf("ERROR: Could not create directory for monitor %s: %v", monitorData.Monitor, err)
			continue
		}
		fmt.Printf("\nProcessing templates for monitor: %s\n", monitorData.Monitor)

		// Use Material 3 classification from color_logic.go
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
				return
			}
			outputFile.Close()
		}
	}
	fmt.Println("\nBuild complete!")
}
