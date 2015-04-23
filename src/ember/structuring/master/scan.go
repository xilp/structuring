package master

import (
	"fmt"
	"os"
	"time"
)

func (p *Master) scan() {
	go func() {
		for {
			for k, v := range p.doings {
				if time.Since(time.Unix(time.Unix(0, v.Created).Unix(), 0)) >= time.Minute {
					p.Push("repush", v)
					delete(p.doings, k)
				}
			}
			p.save()
			time.Sleep(time.Minute)
			fmt.Fprintf(os.Stderr, "scan every minute\n")
		}
	}()
}
