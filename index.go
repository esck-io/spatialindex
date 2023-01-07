package spatialindex

import (
	"math"
	"sync"

	"github.com/puzpuzpuz/xsync"
)

type TileId struct {
	x int
	z int
}

type Querier[T any] interface {
	ListTiles(tileSize float64) []TileId
	Contains(*Node[T]) bool
}

type Index[T any] struct {
	TileSize float64
	tiles    *xsync.MapOf[TileId, *tile[T]]
	init     sync.Once
}

func (i *Index[T]) initialize() {
	i.init.Do(func() {
		i.tiles = xsync.NewTypedMapOf[TileId, *tile[T]](func(id TileId) uint64 {
			return uint64(31*id.x + id.z)
		})
	})
}

func (i *Index[T]) Create(val T, pos [3]float64) *Node[T] {
	i.initialize()
	tileId := i.tileFor(pos)
	n := &Node[T]{pos, tileId, val}

	tile, _ := i.tiles.LoadOrStore(i.tileFor(pos), &tile[T]{})
	tile.add(n)
	return n
}

func (i *Index[T]) Update(n *Node[T], pos [3]float64) {
	i.initialize()
	oldTileId := n.prevTile
	newTileId := i.tileFor(pos)

	if oldTileId == newTileId {
		return
	}

	oldTile, _ := i.tiles.LoadOrStore(oldTileId, &tile[T]{})
	newTile, _ := i.tiles.LoadOrStore(newTileId, &tile[T]{})

	oldTile.transferTo(n, newTile)

	n.pos = pos
	n.prevTile = newTileId
}

func (i *Index[T]) Remove(n *Node[T]) {
	i.initialize()
	oldTileId := n.prevTile
	oldTile, _ := i.tiles.LoadOrStore(oldTileId, &tile[T]{})
	oldTile.remove(n)
}

func (i *Index[T]) Query(q Querier[T], out []*Node[T]) []*Node[T] {
	i.initialize()
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

func (i *Index[T]) tileFor(pos [3]float64) TileId {
	return TileId{
		x: int(math.Floor(pos[0] / i.TileSize)),
		z: int(math.Floor(pos[2] / i.TileSize)),
	}
}
