package main

import (
	"flag"
	"fmt"
	"os"
)

func printUsage() {
	fmt.Printf("%s is a tool for updating Gandi domain IP address.\n\n", os.Args[0])
	flag.PrintDefaults()
}
