package core

import (
	"fmt"
	"github.com/hzget/mmo/core/pb"
	"github.com/hzget/tcpserver"
	"google.golang.org/protobuf/proto"
)

type MoveRouter struct {
	tcpserver.BaseRouter
}

func (m *MoveRouter) Handle(req tcpserver.Request) error {
	conn := req.Conn()
	p, err := conn.GetProperty("player")
	if err != nil {
		return fmt.Errorf("GetProperty failed: %v", err)
	}

	pos := &pb.Position{}
	if err = proto.Unmarshal(req.Msg().Data(), pos); err != nil {
		return fmt.Errorf("GetProperty failed: %v", err)
	}

	player := p.(Player)
	player.UpdatePos(Pos{pos.X, pos.Y, pos.Z, pos.V})

	// !!! Attention !!!
	// this func implementation needs a mutex
	WorldMgr.UpdatePlayer(player)
	return nil
}
