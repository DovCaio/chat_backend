package main

import (
	"log"
	"net/http"

	handler "github.com/DovCaio/chat_backend/handle/chat"
	chat "github.com/DovCaio/chat_backend/internal/chat"
)

func main() {
	hub := chat.NewHub()
	wsHandler := handler.NewWebSocketHandler(hub)

	mux := http.NewServeMux()

	mux.Handle("/ws", wsHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "web/index.html")
	})

	mux.Handle("GET /app.js", http.FileServer(http.Dir("web")))
	mux.Handle("GET /style.css", http.FileServer(http.Dir("web")))

	log.Println("Server avalible in http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
