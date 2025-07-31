package main

import (
	"fmt"
	"log"
	"d/GITVIEW/PromeConfig/backend/pkg/password"
)

func main() {
	hash, err := password.Hash("password123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Password hash for 'password123': %s\n", hash)
}