package main

import (
	"fmt"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
