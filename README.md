
# **ArchThemeM0d**

<div align="center">
<img src="archThemeM0d.png" alt="ArchThemeM0d Logo" width="200" />

<p align="center">


  <!-- Arch Linux -->
  <img src="https://img.shields.io/badge/Arch%20Linux-supported-1793D1?style=for-the-badge&logo=arch-linux&logoColor=white" />

  <!-- Hyprland -->
  <img src="https://img.shields.io/badge/Hyprland-required-00A6FF?style=for-the-badge&logo=hyprland&logoColor=white" />

  <!-- Go -->
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" />

  <!-- License -->
  <img src="https://img.shields.io/badge/license-MIT-yellow?style=for-the-badge" />

</p>

</div>

<div align="center">
ArchThemeM0d is a dynamic theming engine for Arch Linux with Hyprland that automatically generates stunning, cohesive application themes directly from your current wallpaper.
</div>

<br>
    
> **IMPORTANT**
> **Project Scope Refactor:**
> This project is shifting from a web-based tool to a native CLI/TUI workflow.
> The Web IDE has been removed in favor of a Bubble Tea terminal interface.
> Commands are being renamed to **create** (was *generate*) and **apply** (was *build*).



# **Table of Contents**

* [Overview](#overview)
* [Features](#features)
* [Real-World Usage](#real-world-usage)
* [Requirements](#requirements)
* [Installation](#installation)
* [Workflow & Commands](#workflow--commands)
* [Template System](#template-system)
* [Color System](#color-system)
* [Hooks (In Development)](#hooks-in-development)
* [TUI (In Development)](#tui-in-development)
* [Configuration](#configuration)
* [Troubleshooting](#troubleshooting)
* [Contributing](#contributing)
* [API Reference](#api-reference)



# **Overview**

ArchThemeM0d transforms your desktop experience by creating a unified color scheme across your entire system based on your wallpaper.

## **How It Works**

The system operates in a pipeline:

1. **Ingestion:** Takes wallpapers from Hyprpaper
2. **Analysis:** Extracts dominant colors and converts to HCT
3. **Classification:** Assigns Material You roles
4. **Generation:** Produces 13-step tonal ramps
5. **Application:** Injects colors into templates and runs hooks



# **Features**

* 🎨 **Material Design 3 Logic**
* 🧠 **Intelligent Color Classification**
* 🖥️ **Monitor-Based Theming**
* 🖼️ **Unified General Mode** *(Planned)*
* 📟 **Interactive TUI** *(In Development)*
* 🪝 **Post-Build Hooks** *(In Development)*
* 🚀 **Atomic `set` Command** *(In Development)*



# **Real-World Usage**

See a full Hyprland integration here:

👉 **github.com/oarabilekoore/workspace57**

This demonstrates:

* Template structure (kitty, waybar, rofi)
* Startup script integration
* Dynamic theming in production environments



# **Requirements**

### **System Requirements**

* Arch Linux
* Hyprland
* hyprpaper
* Go 1.21+

### **Dependencies**

* `hyprctl`
* `hyprpaper`



# **Installation**

### **From Source**

```bash
# Clone the repository
git clone https://github.com/oarabilekoore/archThemeM0d
cd archThemeM0d

# Build the binary
go build -o archThemeM0d main.go

# Install to your PATH (optional)
sudo mv archThemeM0d /usr/local/bin/
```



# **Directory Structure**

Create the following after installation:

```
~/Templates/ThemeM0d/
├── Templates/          # .tmpl files
├── Themes/             # Generated themes
├── hooks/              # Post-build scripts
└── currenttheme.tm0d   # Auto-generated palette
```



# **Workflow & Commands**

## **create**

*(Formerly `generate`)*

```bash
archThemeM0d create
```

**What it does:**

* Queries hyprpaper for active wallpapers
* Extracts dominant colors
* Classifies them using Material You rules
* Writes palette → `currenttheme.tm0d`



## **apply**

*(Formerly `build`)*

```bash
archThemeM0d apply
```

**What it does:**

* Reads palette data
* Generates tonal palettes
* Processes `.tmpl` files
* Runs hooks



## **set** *(In Development)*

```bash
archThemeM0d set <path/to/image>
```

**What it does:**

* Sets wallpaper
* Runs `create`
* Runs `apply`
* Executes hooks



# **Template System**

Templates use **Go's text/template**.

### **Data Structure**

Available under `.Theme`:

* `.Primary`
* `.Secondary`
* `.Tertiary`
* `.Neutral`
* `.Surface`, `.OnSurface`, etc.

### **Functions**

#### `toHex`

```go
{{ .Theme.Primary | tone 80 | toHex }}
```

#### `toRGB`

```go
{{ .Theme.Surface | toRGB }}
```

#### `tone`

```go
{{ .Theme.Primary | tone 50 }}
{{ .Theme.Neutral | tone 10 }}
```



# **Color System**

### **Design Principles**

* HCT Color Space
* Intelligent Classification
* Hue Harmony
* 13-tone ramps

### **Color Roles**

* **Primary**
* **Secondary**
* **Tertiary**
* **Neutral**

### **Tone Levels**

| Tone | Description |
| - | -- |
| 0    | Black       |
| 10   | Very Dark   |
| 20   | Dark        |
| 50   | Medium      |
| 80   | Light       |
| 90   | Very Light  |
| 100  | White       |

### **Advanced Classification Algorithm**

```
Score = (chroma / 100 * 0.7) + ((1 - abs(tone-50)/50) * 0.3)
```



# **Hooks** *(In Development)*

Store scripts in:

```
~/Templates/ThemeM0d/hooks/
```

Example hook:

```bash
#!/bin/bash
pkill waybar
waybar &
```



# **TUI** *(In Development)*

Planned features:

* Terminal wallpaper selector
* Palette coherency checks (DeltaE)
* ANSI palette preview



# **Configuration**

### **Multi-Monitor Example**

```json
[
  {
    "monitor": "DP-1",
    "theme": {
      "wallpaper_location": "/path/to/img1.jpg",
      "palletes": []
    }
  },
  {
    "monitor": "HDMI-A-1",
    "theme": {
      "wallpaper_location": "/path/to/img2.jpg",
      "palletes": []
    }
  }
]
```



# **Troubleshooting**

### **"This only works with arch hyprland"**

Ensure `HYPRLAND_INSTANCE_SIGNATURE` is set.

### **"No Wallpapers Found"**

Check hyprpaper:

```bash
hyprctl hyprpaper listactive
```



# **Contributing**

### **Development Setup**

```bash
git clone https://github.com/oarabilekoore/archThemeM0d
cd archThemeM0d
go mod tidy
go test ./...
go build -o archThemeM0d main.go
```



# **API Reference**

### **ClassifiedTheme**

```go
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
```

### **TonalPalette**

```go
type TonalPalette struct {
    Tones map[int]color.RGBA
}
```

### **HCT**

```go
type HCT struct {
    H float64
    C float64
    T float64
}
```
