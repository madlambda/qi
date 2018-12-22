package main

import (
	"log"
	"net/http"
)

// Version of yang
const Version = "v0.1"

func wasmHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/wasm")
	http.ServeFile(w, r, "yin.wasm")
}
func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("/yin.wasm", wasmHandler)
	log.Fatal(http.ListenAndServe(":3000", mux))
}
