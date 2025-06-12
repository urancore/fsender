package main

import (
	"fsender/config"
	"fsender/internal/handler"
	"fsender/internal/middleware"
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

	middl := middleware.NewMiddleware()


	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	mux.Handle("HEAD /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /", hand.HandleIndexPage)
	mux.HandleFunc("POST /", hand.FileUploadHandler)
	mux.HandleFunc("GET /f/{key}", hand.GetFileByLink)
	mux.HandleFunc("GET /f/{key}/{filename}", hand.ServeFile)

	wrappedMux := middl.LoggingMiddleware(mux)

	http.ListenAndServe(cfg.Server.Addr, wrappedMux)
}
