package spatialindex

import "testing"

type DumpQuerier struct {
	tilesToDump []TileId
}

func (q *DumpQuerier) ListTiles(tileSize float64) []TileId {
	return q.tilesToDump
}

func (q *DumpQuerier) Contains(n *Node[any]) bool {
	return true
}

func TestIndexNodeCanBeQueriedAfterCreate(t *testing.T) {
	sut := Index[any]{TileSize: 200}

	node := sut.Create(nil, [3]float64{100, 200, 300})
	res := []*Node[any]{}
	res = sut.Query(&DumpQuerier{tilesToDump: []TileId{{0, 1}}}, res)

	if len(res) != 1 {
		t.Errorf("Expected %d nodes in query result, but got %d", 1, len(res))
	}

	if res[0] != node {
		t.Errorf("Expected query result to contain created node")
	}
}

func TestIndexNodeCanBeQueriedAfterUpdate(t *testing.T) {
	sut := Index[any]{TileSize: 200}
	node := sut.Create(nil, [3]float64{100, 200, 300})

	sut.Update(node, [3]float64{1100, 0, 1300})
	res := []*Node[any]{}
	res = sut.Query(&DumpQuerier{tilesToDump: []TileId{{5, 6}}}, res)

	if len(res) != 1 {
		t.Errorf("Expected %d nodes in query result, but got %d", 1, len(res))
	}

	if res[0] != node {
		t.Errorf("Expected query result to contain created node")
	}
}

func TestIndexNodeCannotBeQueriedAfterRemove(t *testing.T) {
	sut := Index[any]{TileSize: 200}
	node := sut.Create(nil, [3]float64{100, 200, 300})

	sut.Remove(node)
	res := []*Node[any]{}
	res = sut.Query(&DumpQuerier{tilesToDump: []TileId{{0, 1}}}, res)

	if len(res) != 0 {
		t.Errorf("Expected %d nodes in query result, but got %d", 0, len(res))
	}
}

type EqualityQuerier struct {
	tilesToSearch []TileId
	nodeToMatch   *Node[any]
}

func (q *EqualityQuerier) ListTiles(tileSize float64) []TileId {
	return q.tilesToSearch
}

func (q *EqualityQuerier) Contains(n *Node[any]) bool {
	return n == q.nodeToMatch
}

func TestIndexQueryFiltersUsingContains(t *testing.T) {
	sut := Index[any]{TileSize: 200}
	node1 := sut.Create(nil, [3]float64{100, 200, 300})
	sut.Create(nil, [3]float64{100, 200, 300})

	res := []*Node[any]{}
	res = sut.Query(&EqualityQuerier{tilesToSearch: []TileId{{0, 1}}, nodeToMatch: node1}, res)

	if len(res) != 1 {
		t.Errorf("Expected %d nodes in query result, but got %d", 1, len(res))
	}

	if res[0] != node1 {
		t.Errorf("Expected query result to contain created node")
	}
}
