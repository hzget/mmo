package main

import (
	"github.com/hzget/mmo/core"
	"github.com/hzget/tcpserver"
	"log"
)

func OnConnStart(conn tcpserver.Conn) {
	player := core.NewPlayer(conn)
	log.Printf("conn [%d] player[%d] online", conn.ConnId(), player.Id())
	conn.AddProperty("player", player)
	core.WorldMgr.AddPlayer(player)
	player.SyncPlayerId()
	player.SyncPlayerPos()
	player.SyncPlayersPos()
}

func OnConnStop(conn tcpserver.Conn) {
	p, err := conn.GetProperty("player")
	if err != nil {
		log.Printf("GetProperty failed: %v", err)
		return
	}
	player := p.(core.Player)
	player.Offline()
	core.WorldMgr.RemovePlayer(player)
	log.Printf("conn [%d] player[%d] offline", conn.ConnId(), player.Id())
}

func main() {

	s := tcpserver.NewServer()
	// register handlers for different msg from the client
	s.AddRouter(core.MSG_C_Talk, &core.ChatRouter{})
	s.AddRouter(core.MSG_C_Move, &core.MoveRouter{})
	// hook func when the player is online
	s.SetOnConnStart(OnConnStart)
	// hook func when the player is offline
	s.SetOnConnStop(OnConnStop)

	s.Serve()
}
