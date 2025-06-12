package handler

import (
	"fmt"
	"fsender/internal/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"html/template"
)

func (h *Handler) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	subdir, err := utils.GenerateKey()
	if err != nil {
		http.Error(w, "Failed to generate directory", http.StatusInternalServerError)
		return
	}

	saveDir := filepath.Join(h.cfg.FTP.RootPath, subdir)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}


		filePath := filepath.Join(saveDir, fileHeader.Filename)
		dst, err := os.Create(filePath)
		if err != nil {
			file.Close()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			file.Close()
			dst.Close()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		file.Close()
		dst.Close()
	}
	ip := h.cfg.Server.Addr
	fmt.Fprintf(w, "Files uploaded successfully. Link: %s/f/%s", ip, subdir)
}

func (h *Handler) HandleIndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl.Execute(w, nil)
}

func (h *Handler) GetFileByLink(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	path := filepath.Join(h.cfg.FTP.RootPath + key)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "Page Not Found 404", http.StatusNotFound)
		return
	}

	funcMap := template.FuncMap{
		"bytesToMB": func(bytes int64) string {
			return fmt.Sprintf("%.2f", float64(bytes)/1024/1024)
		},
		"bytesToKB": func(bytes int64) string {
			return fmt.Sprintf("%.2f", float64(bytes)/1024)
		},
	}

	tmpl := template.New("load_page.html").Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles("templates/load_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files, err := utils.ReadDir(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Files []os.FileInfo
		Key   string
	}{
		Files: files,
		Key:   key,
	}
	tmpl.Execute(w, data)
}

func (h *Handler) ServeFile(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	filename := r.PathValue("filename")
	filePath := filepath.Join(h.cfg.FTP.RootPath, key, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(w, r, filePath)
}
