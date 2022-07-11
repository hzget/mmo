package core

import (
	"fmt"
	"log"
	"sync"
)

type Grid interface {
	AddPlayer(pid int)
	RemovePlayer(pid int)
	GetPlayers() (s []int)
	String() string
	Id() int
}

type grid struct {
	id   int
	minx int
	maxx int
	miny int
	maxy int
	pids map[int]bool
	m    sync.Mutex
}

func NewGrid(id, minx, maxx, miny, maxy int) Grid {
	return &grid{id: id, minx: minx, maxx: maxx,
		miny: miny, maxy: maxy,
		pids: make(map[int]bool),
	}
}

func (g *grid) Id() int {
	return g.id
}

func (g *grid) AddPlayer(pid int) {
	g.m.Lock()
	defer g.m.Unlock()
	g.pids[pid] = true
	log.Printf("grid[%d] add player[%d]", g.id, pid)
}

func (g *grid) RemovePlayer(pid int) {
	g.m.Lock()
	defer g.m.Unlock()
	delete(g.pids, pid)
	log.Printf("grid[%d] remove player[%d]", g.id, pid)
}

func (g *grid) GetPlayers() []int {
	g.m.Lock()
	defer g.m.Unlock()
	var s []int
	for k, _ := range g.pids {
		s = append(s, k)
	}
	return s
}

func (g *grid) String() string {
	return fmt.Sprintf("Grid id:%d, minx:%d, maxx:%d, miny:%d, maxy:%d, pids: %v",
		g.id, g.minx, g.maxx, g.miny, g.maxy, g.pids)
}
