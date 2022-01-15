package Game

import (
	"fmt"
	"sync"
)

// 一个aoi地图中的格子类型

type Grid struct {
	// 格子ID
	GID int
	// 格子左边边界坐标
	MinX int
	// 格子右边边界坐标
	MaxX int
	// 格子上边边界坐标
	MinY int
	// 格子下边边界坐标
	MaxY int
	// 当前格子内的玩家
	playerIDs map[int]bool
	// 保护当前集合的锁
	pIDLock sync.RWMutex
}

// NewGrid 初始化当前的格子的方法
func NewGrid(gID, minX, maxY, MinY, MaxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxY,
		MinY:      MinY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// Add 给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// Remove 从格子中删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// GetPlayerIDs 得到当前格子中所有的玩家
func (g *Grid) GetPlayerIDs() (playersIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playersIDs = append(playersIDs, k)
	}
	return
}

// 调试使用---打印出格子的基本信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d, minX: %d, maxX: %d, minY:%d,maxY:%d,playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
