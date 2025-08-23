package cmd

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	// Blank imports for image decoding
	_ "image/jpeg"
	_ "image/png"

	"github.com/cascax/colorthief-go"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate a json file with the color pallete from your wallpaper",
	Run:   GenerateThemeFile,
}

type WallpaperInfo struct {
	WallpaperPath string       `json:"wallpaper_location"`
	Palletes      []color.RGBA `json:"palletes"`
}

type MonitorInfo struct {
	Monitor string        `json:"monitor"`
	Theme   WallpaperInfo `json:"theme"`
}

var homeDir string

const tm0dDir string = "Templates/ThemeM0d"

var themeFileDir = filepath.Join(tm0dDir, "currenttheme.tm0d")

func init() {
	homeDir = os.Getenv("HOME")
	rootCmd.AddCommand(generateCmd)
}

func getWallpaper() (map[string]string, error) {
	cmd := exec.Command("hyprctl", "hyprpaper", "listactive")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("EROR: Could not get wallpaper: %w", err)
	}

	wallpapers := make(map[string]string)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		monitorName := strings.TrimSpace(parts[0])
		wallpaperPath := strings.TrimSpace(parts[1])

		wallpapers[monitorName] = wallpaperPath
	}

	if len(wallpapers) == 0 {
		return nil, fmt.Errorf("No Wallpapers Found")
	}

	return wallpapers, nil
}

func getDominantColors(imagePath string) ([]color.Color, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	palette, err := colorthief.GetPalette(img, 12)
	if err != nil {
		return nil, fmt.Errorf("failed to get palette: %w", err)
	}

	return palette, nil
}

func DoesThemeM0dFolderExist() (bool, error) {
	info, err := os.Stat(filepath.Join(homeDir, "Templates/ThemeM0d"))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func GenerateThemeFile(cmd *cobra.Command, args []string) {
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") == "" {
		fmt.Println("This only works with arch hyprland")
		return
	}

	wallpapers, err := getWallpaper()
	if err != nil {
		log.Fatalf("ERROR: Could not get wallpaper: %s", err)
	}

	themeDir := filepath.Join(homeDir, "Templates/ThemeM0d")
	if exists, err := DoesThemeM0dFolderExist(); err != nil {
		log.Fatalf("Error checking folder: %v", err)
	} else if !exists {
		err := os.MkdirAll(themeDir, 0755)
		if err != nil {
			log.Fatalf("ERROR: An eror occured trying to make ThemeM0d Directory: %s", err)
		}
	}

	var allMonitorsInfo []MonitorInfo

	for monitor, path := range wallpapers {
		fmt.Printf("Processing wallpaper for monitor %s: %s\n", monitor, path)
		colors, err := getDominantColors(path)

		// We need to convert color.Color to color.RGBA
		rgbaPalette := make([]color.RGBA, 0, len(colors))
		for _, c := range colors {
			if rgba, ok := c.(color.RGBA); ok {
				rgbaPalette = append(rgbaPalette, rgba)
			}
		}

		if err != nil {
			log.Printf("Could not process wallpaper %s: %v. Skipping.", path, err)
			continue
		}

		info := MonitorInfo{
			Monitor: monitor,
			Theme: WallpaperInfo{
				WallpaperPath: path,
				Palletes:      rgbaPalette,
			},
		}
		allMonitorsInfo = append(allMonitorsInfo, info)
	}

	jsonData, err := json.MarshalIndent(allMonitorsInfo, "", "  ")
	if err != nil {
		log.Fatalf("ERROR: Failed to generate JSON: %v", err)
	}

	outputFile := filepath.Join(themeDir, "currenttheme.tm0d")
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("ERROR: Failed to write theme file: %v", err)
	}

	fmt.Printf("Successfully generated theme file at: %s\n", outputFile)
}
