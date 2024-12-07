package main

import (
	"fmt"
	"net/http"
)

// This is an empty site that is needed
// to debug the parser.
// Use only in development mode

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<html><body> 
								<a href="/1">1</a>
								<a href="/2">2</a>
								<a href="/3">3</a>
							</body></html>`)
	})

	http.HandleFunc("/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<html><body> 
								<a href="/6">6</a>
							</body></html>`)
	})

	http.HandleFunc("/2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<html><body> 
							</body></html>`)
	})

	http.HandleFunc("/3", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<html><body> 
								<a href="/4">4</a>
							</body></html>`)
	})

	http.HandleFunc("/4", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<html><body> 
								<a href="/5">5</a>
							</body></html>`)
	})

	http.HandleFunc("/5", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<html><body> 
								<a href="/1">1</a>
								<a href="/2">2</a>
							</body></html>`)
	})

	http.HandleFunc("/6", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<html><body> 
								<a href="/4">4</a>
							</body></html>`)
	})

	fmt.Println("Сервер запущен на http://localhost:9090")
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
