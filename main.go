package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "test server")
	})
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
