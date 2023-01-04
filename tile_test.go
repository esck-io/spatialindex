package spatialindex

import "testing"

func TestAddInsertsNodeIntoValues(t *testing.T) {
	sut := newTile[any]()
	node := &Node[any]{}

	sut.add(node)

	res := sut.values()

	if len(res) != 1 {
		t.Errorf("expected tile to have %d values, but had %d values", 1, len(res))
	}

	if res[0] != node {
		t.Errorf("expected tile to contain test node")
	}
}

func TestRemoveRemovesNodeFromValues(t *testing.T) {
	sut := newTile[any]()
	node := &Node[any]{}
	sut.add(node)

	sut.remove(node)

	res := sut.values()

	if len(res) != 0 {
		t.Errorf("expected tile to have %d values, but had %d values", 0, len(res))
	}
}

func TestTransferAddsToDestinationAndRemovesFromSelf(t *testing.T) {
	sutSrc := newTile[any]()
	sutDst := newTile[any]()
	node := &Node[any]{}
	sutSrc.add(node)

	sutSrc.transferTo(node, sutDst)

	srcRes := sutSrc.values()
	dstRes := sutDst.values()

	if len(srcRes) != 0 {
		t.Errorf("expected source tile to have %d values, but had %d values", 0, len(srcRes))
	}

	if len(dstRes) != 1 {
		t.Errorf("expected dest tile to have %d values, but had %d values", 1, len(dstRes))
	}

	if dstRes[0] != node {
		t.Errorf("expected destination tile to contain test node")
	}
}

func TestValuesAreImmutable(t *testing.T) {
	sut := newTile[any]()
	node := &Node[any]{}

	res := sut.values()
	sut.add(node)

	if len(res) != 0 {
		t.Errorf("expected value set to have %d values, but had %d values", 0, len(res))
	}
}
