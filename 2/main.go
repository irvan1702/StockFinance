package main

import (
	"net/http"
)

// cons(
// 	text := <p>Hello World</p>
// )

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
