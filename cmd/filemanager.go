package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Struct to hold file info
type FileInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// GET: /files
func GetAllTemplateFiles(w http.ResponseWriter, r *http.Request) {
	dir := filepath.Join(homeDir, "Templates/ThemeM0d/Templates")
	files, err := os.ReadDir(dir)
	if err != nil {
		http.Error(w, "Unable to get list of template files", http.StatusInternalServerError)
		log.Printf("An error occurred reading filesystem: %v", err)
		return
	}

	var fileList []FileInfo
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, FileInfo{
				Name: file.Name(),
				Path: filepath.Join(dir, file.Name()),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileList)
}

// GET: /read?file=<filename>
func ReadFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
		return
	}

	path := filepath.Join(homeDir, "Templates/ThemeM0d/Templates", fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		log.Printf("Error reading file %s: %v", path, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}

// POST: /update?file=<filename>
// Body = new file contents
func UpdateFile(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
		return
	}

	path := filepath.Join(homeDir, "Templates/ThemeM0d/Templates", fileName)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = os.WriteFile(path, body, 0644)
	if err != nil {
		http.Error(w, "Unable to write file", http.StatusInternalServerError)
		log.Printf("Error writing file %s: %v", path, err)
		return
	}

	fmt.Fprintf(w, "File %s updated successfully", fileName)
	BuildTemplates(nil, nil)
}
