package cli

import (
	"fmt"
	"os"
	"ember/cli"
	"ember/http/rpc"
	"ember/structuring/master"
)

func Run(args []string) {
	addr, args := cli.PopArg("addr", "http://127.0.0.1:9000", args)
	//addr, args := cli.PopArg("addr", "http://127.0.0.1:8888", args)
	client, err := NewClient(addr)
	cli.Check(err)

	cmds := cli.NewCmds()
	cmds.Reg("api", "call api", client.CmdCall)

	cmds.Run(os.Args[1:])
}

func (p *Client) CmdCall(args []string) {
	ret, err := p.Rpc.Call(args)
	cli.Check(err)
	fmt.Println(ret)
}

func NewClient(addr string) (p *Client, err error) {
	p = &Client{Rpc: rpc.NewClient(addr)}
	err = p.Rpc.Reg(p, &master.Master{})
	return
}

type Client struct {
	Rpc *rpc.Client
	Fetch func(url string) error
	Search func(key string) (string,error)
	Slaves func() ([]string, error)
	Dones func() ([]string, error)
}
