package core

const (
	MSG_S_SyncPid         = 1
	MSG_C_Talk            = 2
	MSG_C_Move            = 3
	MSG_S_Broadcast       = 200
	MSG_S_ClientOffline   = 201
	MSG_S_SyncNeighborMsg = 202
)

const (
	BroadCast_Talk      int32 = 1
	BroadCast_Pos             = 2
	BroadCast_Action          = 3
	BroadCast_UpdatePos       = 4
)

const (
	AOI_MIN_X = 85
	AOI_MAX_X = 410
	AOI_CNT_X = 10
	AOI_MIN_Y = 75
	AOI_MAX_Y = 400
	AOI_CNT_Y = 20
)

var WorldMgr WorldManager

func initGlobals() {
	WorldMgr = NewWorldManager()
}

func init() {
	initGlobals()
}
