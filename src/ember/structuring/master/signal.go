package master

import (
	"fmt"
	"os"
	"os/signal"
)

func (p *Master) catchSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("received ctrl+c(%v)\n", sig)
			p.save()
			os.Exit(0)
		}
	}()
}

