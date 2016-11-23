package main

import (
	"fmt"
	"github.com/sethmcl/gofrosty/lib"
	"log"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	args := os.Args[2:]

	switch os.Args[1] {
	case "install":
		err := lib.InstallCmdRun(args)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("usage: frosty <command> [args]")
	fmt.Println("  frosty install -- Install dependencies from npm-shrinkwrap.json")
}
