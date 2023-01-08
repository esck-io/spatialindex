# Overview

A lightweight, thread-safe, tile-based spatial index for nodes in three dimensions. 

# Usage

```go
// create some data to store
type MySpatialType struct {
	Name     string
	Position [3]float64
}
myThing := MySpatialType{
	Name:     "Testy McTestFace",
	Position: [3]float64{100, 200, 300},
}

// create the index
index := spatialindex.Index[MySpatialType]{TileSize: 200}

//create a node for your data
node := index.Create(myThing, myThing.Position)

//change the position of your data
myThing.Position = [3]float64{1000, 900, 200}

//update the indexed position
//you can update the position in your struct without updating the index if
//you can tolerate some innacuracy in your queries
index.Update(node, myThing.Position)

//query the index to retrieve all of the nodes in a given region
myQuery := spatialindex.SphereQuerier[MySpatialType]{
	Center: [3]float64{1000, 900, 200},
	Radius: 100,
}
out := []*spatialindex.Node[MySpatialType]{}
out = index.Query(&myQuery, out)

//get your data out of the query result
myThing == out[0].Value //true
```
