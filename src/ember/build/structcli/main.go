package main

import (
	"os"
	"ember/structuring/cli"
)

func main() {
	cli.Run(os.Args[1:])
}
