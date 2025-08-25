package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve command allows you to use an interactive editor to modify & build your templates",
	Run:   StartThemeIDEServer,
}

var (
	port int
)

func init() {
	rootCmd.AddCommand(ServeCmd)
	ServeCmd.Flags().IntVar(&port, "port", 8080, "the port to serve the IDE")
}

func StartThemeIDEServer(cmd *cobra.Command, args []string) {
	addr := fmt.Sprintf(":%d", port)
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.Handle("/assets/", http.StripPrefix("/assets",
		http.FileServer(http.Dir("../web/dist/assets"))))

	fmt.Printf("Starting server at http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}

}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("../web/dist/index.html")
	if err != nil {
		http.Error(w, "Failed to serve index.html", http.StatusInternalServerError)
		log.Fatalf("Failed to serve index.html: %v", err)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(file)
}
