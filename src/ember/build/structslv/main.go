package main

import (
	"os"
	"ember/structuring/slave"
)

func main() {
	slave.Run(os.Args[1:])
}
