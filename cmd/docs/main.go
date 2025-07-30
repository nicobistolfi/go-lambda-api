package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := "8081"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	fs := http.FileServer(http.Dir("./docs"))
	http.Handle("/", fs)

	fmt.Printf("Swagger documentation available at http://localhost:%s/\n", port)
	fmt.Println("Press Ctrl+C to stop")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
