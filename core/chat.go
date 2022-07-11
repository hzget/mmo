package core

import (
	"fmt"
	"github.com/hzget/tcpserver"
)

type ChatRouter struct {
	tcpserver.BaseRouter
}

func (c *ChatRouter) Handle(req tcpserver.Request) error {
	conn := req.Conn()
	player, err := conn.GetProperty("player")
	if err != nil {
		return fmt.Errorf("GetProperty failed: %v", err)
	}
	player.(Player).Talk(string(req.Msg().Data()))
	return nil
}
