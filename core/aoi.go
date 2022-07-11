package core

import (
	"fmt"
	//	"log"
)

type AOI interface {
	String() string
	GetNeighborsByGId(gid int) []Grid
	GetGIdByPos(x, y float32) int
	GetNeighborsByPos(x, y float32) []Grid
	GetNeighborPlayersByPos(x, y float32) []int
	AddPlayerById(pid, gid int)
	RemovePlayerById(pid, gid int)
	GetPlayersByGid(gid int) []int
	AddPlayerByPos(pid int, x, y float32)
	RemovePlayerByPos(pid int, x, y float32)
	UpdatePlayer(p Player)
}

type aoi struct {
	minx  int
	maxx  int
	cntx  int
	miny  int
	maxy  int
	cnty  int
	grids map[int]Grid
}

func NewAOI(minx, maxx, cntx, miny, maxy, cnty int) AOI {
	v := &aoi{
		minx: minx, maxx: maxx, cntx: cntx,
		miny: miny, maxy: maxy, cnty: cnty,
		grids: make(map[int]Grid),
	}

	w := v.gridwidth()
	l := v.gridlen()
	for y := 0; y < v.cnty; y++ {
		for x := 0; x < v.cntx; x++ {
			id := x + y*v.cntx
			v.grids[id] = NewGrid(id,
				v.minx+x*w,
				v.minx+x*w+x,
				v.miny+y*l,
				v.miny+y*l+y,
			)
		}
	}

	return v
}

func (a *aoi) gridwidth() int {
	return (a.maxx - a.minx) / a.cntx
}

func (a *aoi) gridlen() int {
	return (a.maxy - a.miny) / a.cnty
}

func (a *aoi) String() string {
	s := fmt.Sprintf("AOI: minx:%d, maxx:%d, cntx:%d, miny:%d, maxy:%d, cnty:%d, grids:\n",
		a.minx, a.maxx, a.cntx, a.miny, a.maxy, a.cnty)

	for _, v := range a.grids {
		s += fmt.Sprintf("%s\n", v)
	}
	return s
}

func (a *aoi) GetNeighborsByGId(gid int) []Grid {
	if _, ok := a.grids[gid]; !ok {
		return nil
	}

	idy := gid / a.cntx
	idx := gid % a.cntx
	var s []Grid
	// ul, up, ur, lf, g, rt, dl, dn, dr
	for y := -1; y < 2; y++ {
		if idy+y < 0 || idy+y >= a.cnty {
			continue
		}
		for x := -1; x < 2; x++ {
			if idx+x < 0 || idx+x >= a.cntx {
				continue
			}
			id := y*a.cntx + x + gid
			// it must exists
			s = append(s, a.grids[id])
		}
	}

	return s
}

func (a *aoi) GetGIdByPos(x, y float32) int {
	idx := (int(x) - a.minx) / a.cntx
	idy := (int(y) - a.miny) / a.cnty
	return idy*a.cntx + idx
}

func (a *aoi) GetGIdByPlayer(p Player) int {
	for _, grid := range a.grids {
		pids := grid.GetPlayers()
		for _, id := range pids {
			if id == int(p.Id()) {
				return grid.Id()
			}
		}
	}

	return -1
}

func (a *aoi) GetNeighborsByPos(x, y float32) []Grid {
	gid := a.GetGIdByPos(x, y)
	return a.GetNeighborsByGId(gid)
}

func (a *aoi) GetNeighborPlayersByPos(x, y float32) []int {
	grids := a.GetNeighborsByPos(x, y)
	var s []int
	for _, g := range grids {
		s = append(s, g.GetPlayers()...)
	}
	//log.Printf("GetNeighborPlayersByPos: pos: [%f, %f], grids: %v, players: %v", x, y, grids, s)
	return s
}

func (a *aoi) AddPlayerById(pid, gid int) {
	a.grids[gid].AddPlayer(pid)
}

func (a *aoi) RemovePlayerById(pid, gid int) {
	//log.Printf("RemovePlayerById: pid[%d], gid[%d]", pid, gid)
	a.grids[gid].RemovePlayer(pid)
}

func (a *aoi) GetPlayersByGid(gid int) []int {
	return a.grids[gid].GetPlayers()
}

func (a *aoi) AddPlayerByPos(pid int, x, y float32) {
	gid := a.GetGIdByPos(x, y)
	a.AddPlayerById(pid, gid)
}

func (a *aoi) RemovePlayerByPos(pid int, x, y float32) {
	gid := a.GetGIdByPos(x, y)
	a.RemovePlayerById(pid, gid)
}

func (a *aoi) UpdatePlayer(p Player) {
	newgid := a.GetGIdByPos(p.Position().X, p.Position().Z)
	oldgid := a.GetGIdByPlayer(p)
	if newgid != oldgid {
		//log.Printf("UpdatePlayer: player[%d] move grid[%d] ---> grid[%d]", p.Id(), oldgrid, newgrid)
		a.RemovePlayerById(int(p.Id()), oldgid)
		a.AddPlayerById(int(p.Id()), newgid)
	}
}
