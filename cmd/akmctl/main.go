// Command akmctl reads, checks and changes the records of an AKM repository.
package main

import (
	"fmt"
	"os"

	"github.com/alexsasharegan/akmctl/internal/cli"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "akmctl:", err)
		os.Exit(2)
	}
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, cwd))
}
