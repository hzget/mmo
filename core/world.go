package core

import (
	"sync"
)

type WorldManager interface {
	Aoi() AOI
	AddPlayer(p Player)
	RemovePlayer(p Player)
	GetPlayer(id int32) Player
	GetPlayers() []Player
	UpdatePlayer(p Player)
}

type worldmanager struct {
	aoi     AOI
	players map[int32]Player
	sync.RWMutex
}

func NewWorldManager() WorldManager {
	aoi := NewAOI(
		AOI_MIN_X, AOI_MAX_X, AOI_CNT_X,
		AOI_MIN_Y, AOI_MAX_Y, AOI_CNT_Y,
	)
	players := make(map[int32]Player)
	return &worldmanager{aoi: aoi, players: players}
}

func (m *worldmanager) Aoi() AOI {
	return m.aoi
}

func (m *worldmanager) AddPlayer(p Player) {
	m.Lock()
	m.players[p.Id()] = p
	m.Unlock()
	m.aoi.AddPlayerByPos(int(p.Id()), p.Position().X, p.Position().Z)
}

func (m *worldmanager) RemovePlayer(p Player) {
	m.aoi.RemovePlayerByPos(int(p.Id()), p.Position().X, p.Position().Z)
	m.Lock()
	delete(m.players, p.Id())
	m.Unlock()
}

func (m *worldmanager) GetPlayer(id int32) Player {
	m.RLock()
	defer m.RUnlock()
	return m.players[id]
}

func (m *worldmanager) GetPlayers() []Player {
	m.RLock()
	defer m.RUnlock()
	var s []Player
	for _, v := range m.players {
		s = append(s, v)
	}
	return s
}

func (m *worldmanager) UpdatePlayer(p Player) {
	m.aoi.UpdatePlayer(p)
}
