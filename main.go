package main

import (
	"fmt"
	"os"

	"github.com/jmehdipour/gift-card/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
