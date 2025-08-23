# ArchThemeM0d 🎨

*ArchThemeM0d* is a dynamic theming engine for Arch with Hyprland that automatically generates a stunning and cohesive set of application themes directly from your current wallpaper.

Powered by Go, it intelligently extracts a color palette and builds a full, Material You-inspired tonal system. This ensures that your Waybar, Rofi, and other supported applications always look beautiful and perfectly matched to your desktop environment.

## ✨ Core Features

- Wallpaper-Based Generation: Automatically creates a rich color palette from your current wallpaper.
- Intelligent Color System: Goes beyond simple color extraction by classifying colors into Primary, Secondary, Tertiary, and Neutral roles.
- Material You-Inspired Tonal Palettes: Generates a full ramp of 13 tones (from dark to light) for each key color, providing a complete and flexible design system.
- Template-Driven: Easily extend support to any application that uses text-based configuration files (like Kitty, Dunst, Alacritty, etc.) by creating simple templates.
- Cohesive Desktop Experience: Ensures all themed applications share a single, harmonious, and professionally designed color scheme.
- Fast & Efficient: Written in Go for excellent performance.


## 🚀 How It Works

The process is split into two simple commands:


1. archThemeM0d generate: This command finds your current wallpaper, extracts a 16-color palette, and saves it to a Templates/ThemeM0d/currenttheme.tm0d (it's json) file. This step only needs to be run when you change your wallpaper.


2. archThemeM0d build: This is where the magic happens.

	- It reads the raw palette from currenttheme.tm0d.
	- It analyzes the colors and assigns the most vibrant ones to be the "seed" for the Primary, Secondary, and Tertiary roles. A neutral color is also selected.
	- It uses these seeds to generate full tonal palettes.
	- Finally, it processes all .tmpl files in your Templates directory, injecting the generated theme to produce the final configuration files in the Themes directory.
