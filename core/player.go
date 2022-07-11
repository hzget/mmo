package core

import (
	"fmt"
	"github.com/hzget/mmo/core/pb"
	"github.com/hzget/tcpserver"
	"google.golang.org/protobuf/proto"
	"log"
	"math/rand"
	"sync/atomic"
)

var pidgen int32 = 0

type Player interface {
	Id() int32
	SendMsg(id uint32, data proto.Message)
	SyncPlayerId()
	SyncPlayerPos()
	SyncPlayersPos()
	Position() Pos
	Talk(content string)
	UpdatePos(pos Pos)
	Offline()
	String() string
}

type player struct {
	id   int32
	conn tcpserver.Conn
	Pos
}

type Pos struct {
	X float32
	Y float32
	Z float32
	V float32
}

func NewPlayer(conn tcpserver.Conn) Player {
	p := &player{
		id:   atomic.AddInt32(&pidgen, 1),
		conn: conn,
	}
	p.X = float32(160 + rand.Intn(10))
	p.Z = float32(140 + rand.Intn(20))
	return p
}

func (p *player) String() string {
	return fmt.Sprintf("player-%d, conn-%d, pos-%v", p.id, p.conn.ConnId(), p.Pos)
}

func (p *player) Id() int32 {
	return p.id
}

func (p *player) Position() Pos {
	return p.Pos
}

func (p *player) SendMsg(id uint32, data proto.Message) {
	mdata, err := proto.Marshal(data)
	if err != nil {
		log.Printf("fail to send msg to player[%d]: %v", p.id, err)
		return
	}

	log.Printf("conn [%d] SendMsg to player[%d]: msgid=%d, data=%v", p.conn.ConnId(), p.id, id, data)

	msg := tcpserver.NewMessage(id, mdata)
	p.conn.SendMsg(msg)
}

// SyncPlayerId: sync its id to the player
func (p *player) SyncPlayerId() {
	p.SendMsg(MSG_S_SyncPid, &pb.SyncPid{Pid: p.id})
}

// SyncPlayerPos sync its pos to the player
func (p *player) SyncPlayerPos() {
	p.SendMsg(MSG_S_Broadcast, p.packData(BroadCast_Pos, p.Pos))
}

// Talk: send the player's chat content to all online players
func (p *player) Talk(content string) {
	data := p.packData(BroadCast_Talk, content)
	players := WorldMgr.GetPlayers()
	for _, p := range players {
		p.SendMsg(MSG_S_Broadcast, data)
	}
}

func (p *player) GetNeighborPlayers() []Player {
	pids := WorldMgr.Aoi().GetNeighborPlayersByPos(p.X, p.Z)
	var players []Player
	for _, pid := range pids {
		if pid == int(p.id) {
			continue
		}
		v := WorldMgr.GetPlayer(int32(pid))
		if v == nil {
			log.Printf("fail to get player[%d] in GetNeighborPlayers", pid)
			continue
		}
		players = append(players, v)
	}
	return players
}

// SyncPlayersPos: sync pos: this player <----> neighbors
func (p *player) SyncPlayersPos() {
	neighbors := p.GetNeighborPlayers()
	if neighbors == nil {
		return
	}

	// sync this player's pos to neighbors
	p.broadcast(MSG_S_Broadcast, p.packData(BroadCast_Pos, p.Pos), neighbors)

	// sync neighbors' pos to this player
	neighborsData := p.getPlayersData(neighbors)
	p.SendMsg(MSG_S_SyncNeighborMsg, &pb.SyncPlayers{Ps: neighborsData})
}

func (*player) getPlayersData(players []Player) []*pb.Player {
	playersData := make([]*pb.Player, 0, len(players))
	for _, v := range players {
		pos := v.Position()
		p := &pb.Player{
			Pid: v.Id(),
			P: &pb.Position{
				X: pos.X,
				Y: pos.Y,
				Z: pos.Z,
				V: pos.V,
			},
		}
		playersData = append(playersData, p)
	}
	return playersData
}

func (p *player) packData(tp int32, data interface{}) proto.Message {
	pmsg := &pb.BroadCast{
		Pid: p.id,
		Tp:  tp,
	}

	switch tp {
	case BroadCast_Talk:
		content := data.(string)
		pmsg.Data = &pb.BroadCast_Content{
			Content: content,
		}
	case BroadCast_Pos:
		fallthrough
	case BroadCast_UpdatePos:
		pos := data.(Pos)
		pmsg.Data = &pb.BroadCast_P{
			P: &pb.Position{
				X: pos.X,
				Y: pos.Y,
				Z: pos.Z,
				V: pos.V,
			},
		}
	}

	return pmsg
}

func (*player) broadcast(msgid uint32, data proto.Message, dest []Player) {
	if dest == nil {
		return
	}
	for _, p := range dest {
		if p == nil {
			continue
		}
		p.SendMsg(msgid, data)
	}
}

// UpdatePos: update the player's pos to players inside the AOI
func (p *player) UpdatePos(pos Pos) {
	p.Pos = pos
	players := p.GetNeighborPlayers()
	p.broadcast(MSG_S_Broadcast, p.packData(BroadCast_UpdatePos, p.Pos), players)
}

// Offline: send offline info to neighbors
func (p *player) Offline() {
	players := p.GetNeighborPlayers()
	log.Printf("Offline: get neighbors: %v", players)
	p.broadcast(MSG_S_ClientOffline, &pb.SyncPid{Pid: p.id}, players)
}
