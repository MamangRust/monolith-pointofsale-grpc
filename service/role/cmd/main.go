package main

import "github.com/MamangRust/monolith-point-of-sale-role/internal/apps"

func main() {
	server, err := apps.NewServer()

	if err != nil {
		panic(err)
	}

	server.Run()
}
