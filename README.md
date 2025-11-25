ArchThemeM0d

<div align="center">
<img src="archThemeM0d.png" alt="ArchThemeM0d Logo" width="200" />
</div>
<div align="center">
ArchThemeM0d is a dynamic theming engine for Arch Linux with Hyprland that automatically generates stunning, cohesive application themes directly from your current wallpaper.
</div>


<br>
    
[!IMPORTANT]
Project Scope Refactor: This project is shifting from a web-based tool to a native CLI/TUI workflow. The Web IDE has been removed in favor of a Bubble Tea terminal interface. Commands are being renamed to create (was generate) and apply (was build).

Table of Contents

Overview

Features

Real-World Usage

Requirements

Installation

Workflow & Commands

Template System

Color System

Hooks (In Development)

TUI (In Development)

Configuration

Troubleshooting

Contributing

API Reference

Overview

ArchThemeM0d transforms your desktop experience by creating a unified color scheme across all your applications based on your wallpaper.

How It Works

The system operates in a pipeline:

Ingestion: Takes a wallpaper (or set of wallpapers) via Hyprpaper

Analysis: Extracts dominant colors and converts them to HCT (Hue, Chroma, Tone) space

Classification: Assigns colors to Material You roles (Primary, Secondary, etc.) based on vibrancy and hue harmony

Generation: Creates 13-step tonal ramps for every color role

Application: Injects these colors into user-defined templates and runs post-build hooks

Features

🎨 Material Design 3 Logic: Generates vibrant, balanced palettes rather than simple hex extraction

🧠 Intelligent Color Classification: Assigns colors to Primary, Secondary, Tertiary, and Neutral roles

🖥️ Monitor-Based Theming: Unique themes for each monitor based on its specific wallpaper

🖼️ Unified General Mode [Planned]: Create a single coherent theme from a collection of wallpapers with color coherency checks

📟 Interactive TUI [In Development]: A Bubble Tea-based terminal interface for selecting wallpapers and previewing palettes

🪝 Post-Build Hooks [In Development]: Automatically execute shell scripts (e.g., reloading Waybar/Kitty) after applying a theme

🚀 Atomic set Command [In Development]: Change wallpaper and theme in a single command

Real-World Usage

To see ArchThemeM0d integrated into a full production Hyprland environment, check out Workspace57:

👉 github.com/oarabilekoore/workspace57

This repository demonstrates:

How to structure your templates (kitty, waybar, rofi)

How to hook ArchThemeM0d into your startup scripts

Practical examples of dotfile management with dynamic theming

Requirements

System Requirements

OS: Arch Linux

WM: Hyprland

Wallpaper Manager: hyprpaper

Go: Version 1.21+

Dependencies

hyprctl (part of Hyprland)

hyprpaper for wallpaper management

Installation

From Source

# Clone the repository
git clone [https://github.com/oarabilekoore/archThemeM0d](https://github.com/oarabilekoore/archThemeM0d)
cd archThemeM0d

# Build the binary
go build -o archThemeM0d main.go

# Install to your PATH (optional)
sudo mv archThemeM0d /usr/local/bin/


Directory Structure

After installation, create the following directory structure in your home directory:

~/Templates/ThemeM0d/
├── Templates/          # Your .tmpl files go here
├── Themes/             # Generated themes (auto-created)
├── hooks/              # Post-build scripts (In Development)
└── currenttheme.tm0d   # Generated palette data (auto-created)


Workflow & Commands

create

(Formerly generate)

Extracts color palette from your current wallpaper(s).

archThemeM0d create


What it does:

Queries hyprpaper for active wallpapers

Extracts 12 dominant colors per wallpaper using advanced algorithms

Classifies colors using Material You principles

Saves palette data to currenttheme.tm0d

apply

(Formerly build)

Processes templates using the generated color palette.

archThemeM0d apply


What it does:

Reads palette data from currenttheme.tm0d

Classifies colors into design roles using HCT color space

Generates complete tonal palettes (13 tones per role)

Processes all .tmpl files in Templates directory

Runs Hooks (See Hooks section)

set [In Development]

A one-shot command that orchestrates the entire pipeline.

archThemeM0d set <path/to/image>


What it does:

Sets the wallpaper via hyprpaper

Generates the palette (create)

Applies templates (apply)

Runs hooks

Template System

Templates use Go's text/template syntax.

Available Data Structure

Your templates have access to the .Theme object:

Palettes: .Theme.Primary, .Theme.Secondary, .Theme.Tertiary, .Theme.Neutral

Fixed Roles: .Theme.Surface, .Theme.OnSurface, .Theme.SurfaceVariant, etc.

Template Functions

toHex

Converts a color to hexadecimal format.

{{.Theme.Primary | tone 80 | toHex}}
// Output: #ff5722


toRGB [New]

Converts a color to an integer array (useful for JSON/Chrome manifests).

{{.Theme.Surface | toRGB}}
// Output: [20, 20, 20]


tone

Extracts a specific tone level (0-100) from a tonal palette.

{{.Theme.Primary | tone 50}}        // Mid-tone
{{.Theme.Neutral | tone 10}}        // Very dark


Color System

Design Principles

ArchThemeM0d follows Material You design principles with advanced color science:

HCT Color Space: Uses Hue, Chroma, Tone for perceptually uniform colors

Intelligent Classification: Analyzes vibrancy, hue relationships, and chroma

Harmonious Relationships: Ensures complementary, triadic, and analogous color harmony

13-Tone Ramps: Complete tonal palettes from dark (0) to light (100)

Color Roles

Primary: Most prominent brand/accent color (highest vibrancy)

Secondary: Supporting color with harmonious hue relationship

Tertiary: Additional accent for balance and variety

Neutral: Low-chroma colors for backgrounds and surfaces

Tone Levels

Tone

Usage

Description

0

Black

Pure black

10

Very Dark

Dark surfaces, text on light

20

Dark

Dark elements

50

Medium

Disabled text

80

Light

Seed color placement

90

Very Light

Text on dark surfaces

100

White

Pure white

Advanced Classification Algorithm

HCT Conversion: All colors converted to perceptually uniform HCT color space.

Vibrancy Calculation: Material You formula combining chroma and tone.

Score = (chroma / 100.0 * 0.7) + ((1.0 - abs(tone-50)/50.0) * 0.3)


Hue Analysis: Identifies complementary (~180°), triadic (~120°), and analogous (~30°) relationships.

Role Assignment: Assigns Primary to highest vibrancy, then searches for harmonious hues for Secondary/Tertiary.

Hooks [In Development]

You can add shell scripts to ~/Templates/ThemeM0d/hooks/ to automatically reload your applications after a theme is applied.

Example reload-waybar.sh:

#!/bin/bash
pkill waybar
waybar &


Note: Ensure scripts are executable (chmod +x).

TUI [In Development]

We are replacing the web-based IDE with a native Terminal User Interface (TUI) built with Bubble Tea.

Planned Features:

Visual File Picker: Browse and select wallpapers directly in the terminal

Coherency Check: Select multiple wallpapers for a unified theme; the system will warn you if their color palettes clash (high DeltaE distance)

Palette Preview: View the generated 13-tone ramps as ANSI color blocks before applying

Configuration

Multi-Monitor Setup

ArchThemeM0d automatically handles multiple monitors via currenttheme.tm0d:

[
  {
    "monitor": "DP-1",
    "theme": {
      "wallpaper_location": "/path/to/img1.jpg",
      "palletes": [...]
    }
  },
  {
    "monitor": "HDMI-A-1",
    "theme": {
      "wallpaper_location": "/path/to/img2.jpg",
      "palletes": [...]
    }
  }
]


Troubleshooting

"This only works with arch hyprland"

Cause: Not running in Hyprland environment
Fix: Ensure HYPRLAND_INSTANCE_SIGNATURE environment variable is set.

"No Wallpapers Found"

Cause: hyprpaper not running or no wallpapers set.
Fix:

hyprctl hyprpaper listactive


Contributing

Development Setup

# Clone and setup
git clone [https://github.com/oarabilekoore/archThemeM0d](https://github.com/oarabilekoore/archThemeM0d)
cd archThemeM0d
go mod tidy

# Run tests
go test ./...

# Build
go build -o archThemeM0d main.go


API Reference

Data Structures

ClassifiedTheme

The master structure containing the generated theme data.

type ClassifiedTheme struct {
    Primary   TonalPalette
    Secondary TonalPalette
    Tertiary  TonalPalette
    Neutral   TonalPalette

    Surface          color.RGBA
    SurfaceVariant   color.RGBA
    OnSurface        color.RGBA
    OnSurfaceVariant color.RGBA
    PrimaryFixed     color.RGBA
    OnPrimaryFixed   color.RGBA
}


TonalPalette

Represents a color scale with tones ranging from 0 (Black) to 100 (White).

type TonalPalette struct {
    Tones map[int]color.RGBA
}


HCT

Perceptual color space representation.

type HCT struct {
    H float64 // Hue (0-360)
    C float64 // Chroma (0-100+)
    T float64 // Tone (0-100)
}
