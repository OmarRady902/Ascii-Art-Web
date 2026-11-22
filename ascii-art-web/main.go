package main

import (
	"ascii-art-web/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.HomePageHandler)
	http.HandleFunc("/ascii-art", handler.AsciiArtHandler)

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	fmt.Printf("Server running on http://localhost:\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
