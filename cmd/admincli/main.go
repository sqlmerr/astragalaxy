package main

import (
	"os"

	"github.com/sqlmerr/astragalaxy/internal/transport/admincli"
)

func main() {
	err := admincli.Execute()
	if err != nil {
		os.Exit(1)
	}
}
