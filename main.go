package main

import (
	"fmt"
	"os"

	"kadai/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}

}

func run() error {
	ew, err := server.New("localhost:8080")
	fmt.Println("Server is running at http://localhost:8080")

	if err != nil {
		return err
	}

	if err := ew.Start(); err != nil {
		return err
	}

	return nil
}
