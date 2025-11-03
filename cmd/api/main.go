package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dhquero/stress-test-go/internal/infra/web"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", web.Handler)

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
