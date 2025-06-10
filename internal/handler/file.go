package handler

import (
	"fmt"
	"fsender/internal/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

func (h *Handler) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// (максимум 10 МБ)
	r.ParseMultipartForm(10 << 20)

	subdir_save_file, err := utils.GenerateKey()
	if err != nil {
		http.Error(w, "File create error", http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// создаем полный путь к целевой директории
	save_dir := filepath.Join(h.cfg.FTP.RootPath, subdir_save_file)

	// создаем директорию со всеми родительскими папками при необходимости
	if err := os.MkdirAll(save_dir, 0755); err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	// создаем полный путь к файлу
	filePath := filepath.Join(save_dir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// копируем содержимое файла
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s / link: http://26.199.90.194:1212/f/%s", handler.Filename, subdir_save_file)
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
		http.Error(w, "Page Not Found", http.StatusNotFound)
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
	}

	files, err := utils.ReadDir(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Проверка существования файла
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовки для скачивания
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	http.ServeFile(w, r, filePath)
}
