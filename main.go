
package main

import (
	"log"
	"os"

	"github.com/MohamedBabker/project1-jwks/internal/httpserver"
	"github.com/MohamedBabker/project1-jwks/internal/keystore"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	ks, err := keystore.NewDefaultStore()
	if err != nil {
		log.Fatalf("keystore init: %v", err)
	}
	srv := httpserver.New(ks)
	log.Printf("listening on :%s", port)
	if err := srv.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
