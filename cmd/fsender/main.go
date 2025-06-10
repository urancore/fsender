package main

import (
	"fsender/config"
	"fsender/internal/handler"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	hand := handler.NewHandler(cfg)
	fs := http.FileServer(http.Dir("static"))

	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.Handle("HEAD /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", hand.HandleIndexPage)
	mux.HandleFunc("POST /", hand.FileUploadHandler)
	mux.HandleFunc("GET /f/{key}", hand.GetFileByLink) // http://localhost:1212/f/2SKM-K70IHxC8
	mux.HandleFunc("GET /f/{key}/{filename}", hand.ServeFile)

	http.ListenAndServe(cfg.Server.Addr, mux)
}
