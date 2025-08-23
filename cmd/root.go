package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "archThemeM0d",
	Short: "archThemeM0d helps you to build cohesive themes across your systems based off your wallpaper.",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
