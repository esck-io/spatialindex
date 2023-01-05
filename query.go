package spatialindex

import "github.com/ungerik/go3d/float64/vec3"

type SphereQuerier[T any] struct {
	Center [3]float64
	Radius float64
}

func (q *SphereQuerier[T]) ListTiles(tileSize float64) []tileId {

	return tilesInRadius(tileSize, q.Center, q.Radius)
}

func (q *SphereQuerier[T]) Contains(n *Node[T]) bool {
	var cVec vec3.T = q.Center
	var nVec vec3.T = n.pos

	return vec3.Distance(&cVec, &nVec) <= q.Radius
}

type CylinderQuerier[T any] struct {
	Center       [3]float64
	Radius       float64
	TopOffset    float64
	BottomOffset float64
}

func (q *CylinderQuerier[T]) ListTiles(tileSize float64) []tileId {

	return tilesInRadius(tileSize, q.Center, q.Radius)
}

func (q *CylinderQuerier[T]) Contains(n *Node[T]) bool {
	var cVec vec3.T = q.Center
	cVec[1] = 0
	var nVec vec3.T = n.pos
	nVec[1] = 0

	return vec3.Distance(&cVec, &nVec) <= q.Radius &&
		n.pos[1] <= q.Center[1]+q.TopOffset &&
		n.pos[1] >= q.Center[1]+q.BottomOffset
}

func tilesInRadius(tileSize float64, center [3]float64, radius float64) []tileId {
	xLow := int((center[0] - radius) / tileSize)
	xHigh := int((center[0] + radius) / tileSize)
	zLow := int((center[2] - radius) / tileSize)
	zHigh := int((center[2] + radius) / tileSize)

	numTiles := (xHigh - xLow) * (zHigh - zLow)
	tiles := make([]tileId, 0, numTiles)

	for x := xLow; x <= xHigh; x++ {
		for z := zLow; z <= zHigh; z++ {
			tiles = append(tiles, tileId{x, z})
		}
	}
	return tiles
}
