package spatialindex

import "sync"

type Node[T any] struct {
	pos      [3]float64
	prevTile tileId
}

type nodeTransfer[T any] struct {
	n *Node[T]
	t *tile[T]
}

type tile[T any] struct {
	create   chan *Node[T]
	transfer chan nodeTransfer[T]
	delete   chan *Node[T]
	read     chan []*Node[T]
	init     sync.Once
}

func (t *tile[T]) initialize() {
	t.init.Do(func() {
		t.create = make(chan *Node[T])
		t.transfer = make(chan nodeTransfer[T])
		t.delete = make(chan *Node[T])
		t.read = make(chan []*Node[T])

		go t.run()
	})
}

func (t *tile[T]) run() {
	data := make([]*Node[T], 0, 1)
	immut := []*Node[T]{}

	for {
		select {
		case n := <-t.create:
			data = append(data, n)
			immut = clone(data)

		case n := <-t.delete:
			data = deleteNode(data, n)
			immut = clone(data)

		case nt := <-t.transfer:
			nt.t.initialize()
			nt.t.create <- nt.n
			data = deleteNode(data, nt.n)
			immut = clone(data)

		case t.read <- immut:
		}
	}
}

func (t *tile[T]) add(n *Node[T]) {
	t.initialize()
	t.create <- n
}

func (t *tile[T]) remove(n *Node[T]) {
	t.initialize()
	t.delete <- n
}

func (t *tile[T]) transferTo(n *Node[T], other *tile[T]) {
	t.initialize()
	t.transfer <- nodeTransfer[T]{n, other}
}

func (t *tile[T]) values() []*Node[T] {
	t.initialize()
	return <-t.read
}

func clone[T any](src []*Node[T]) []*Node[T] {
	dst := make([]*Node[T], len(src))
	copy(dst, src)
	return dst
}

func deleteNode[T any](data []*Node[T], n *Node[T]) []*Node[T] {
	for i, dn := range data {
		if dn == n {
			data[i] = data[len(data)-1]
			data = data[:len(data)-1]
		}
	}

	return data
}
