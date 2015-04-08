package main

import (
	"os"
	"ember/structuring/master"
)

func main() {
	master.Run(os.Args[1:])
}
