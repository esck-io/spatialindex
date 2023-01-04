package spatialindex

import (
	"math"
	"sync"

	"github.com/puzpuzpuz/xsync"
)

type tileId struct {
	x int
	z int
}

type Querier[T any] interface {
	ListTiles(tileSize float64) []tileId
	Contains(*Node[T]) bool
}

type Index[T any] struct {
	TileSize float64
	tiles    *xsync.MapOf[tileId, *tile[T]]
	init     sync.Once
}

func (i *Index[T]) initialize() {
	i.init.Do(func() {
		i.tiles = xsync.NewTypedMapOf[tileId, *tile[T]](func(id tileId) uint64 {
			return uint64(31*id.x + id.z)
		})
	})
}

func (i *Index[T]) Create(val T, pos [3]float64) *Node[T] {
	tileId := i.tileFor(pos)
	n := &Node[T]{pos, tileId}

	tile, _ := i.tiles.LoadOrStore(i.tileFor(pos), &tile[T]{})
	tile.add(n)
	return n
}

func (i *Index[T]) Update(n *Node[T], pos [3]float64) {
	oldTileId := n.prevTile
	newTileId := i.tileFor(pos)

	oldTile, _ := i.tiles.LoadOrStore(oldTileId, &tile[T]{})
	newTile, _ := i.tiles.LoadOrStore(newTileId, &tile[T]{})

	oldTile.transferTo(n, newTile)

	n.pos = pos
	n.prevTile = newTileId
}

func (i *Index[T]) Remove(n *Node[T]) {
	oldTileId := n.prevTile
	oldTile, _ := i.tiles.LoadOrStore(oldTileId, &tile[T]{})
	oldTile.remove(n)
}

func (i *Index[T]) Query(q Querier[T], out []*Node[T]) []*Node[T] {
	tiles := q.ListTiles(i.TileSize)

	for _, id := range tiles {
		tile, _ := i.tiles.Load(id)
		vals := tile.values()

		for _, val := range vals {
			if q.Contains(val) {
				out = append(out, val)
			}
		}
	}

	return out
}

func (i *Index[T]) tileFor(pos [3]float64) tileId {
	return tileId{
		x: int(math.Floor(pos[0] / i.TileSize)),
		z: int(math.Floor(pos[2] / i.TileSize)),
	}
}
