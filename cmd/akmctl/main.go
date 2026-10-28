// Command akmctl reads, checks and changes the records of an AKM repository.
package main

import (
	"fmt"
	"os"
)

const version = "0.0.0-dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("akmctl", version)
		return
	}
	fmt.Fprintln(os.Stderr, "akmctl: not implemented yet; see README.md")
	os.Exit(2)
}
