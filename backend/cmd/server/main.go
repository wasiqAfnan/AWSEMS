package main

import (
	"fmt"
	"log"
	"net/http"

	"awsems/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "AWSEMS server is running")
	})

	log.Printf("Server running on port %s", cfg.AppPort)

	err = http.ListenAndServe(":"+cfg.AppPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
