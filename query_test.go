package spatialindex

import (
	"math"
	"testing"
)

func TestSphereReturnsCorrectTilesWhenCenteredOnTileIntersection(t *testing.T) {
	sut := SphereQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 49.0}

	res := sut.ListTiles(50)

	if len(res) != 4 {
		t.Errorf("Expected %d tiles to be returned, but got %d", 4, len(res))
	}

	assertContainsTile(t, res, TileId{1, 5})
	assertContainsTile(t, res, TileId{2, 5})
	assertContainsTile(t, res, TileId{2, 5})
	assertContainsTile(t, res, TileId{2, 6})
}

func TestSphereReturnsCorrectTilesWhenContainedWithinTile(t *testing.T) {
	sut := SphereQuerier[any]{[3]float64{150.0, 200.0, 350.0}, 49.0}

	res := sut.ListTiles(200)

	if len(res) != 1 {
		t.Errorf("Expected %d tiles to be returned, but got %d", 1, len(res))
	}

	assertContainsTile(t, res, TileId{0, 1})
}

func TestSphereReturnsCorrectTilesForNegativePosition(t *testing.T) {
	sut := SphereQuerier[any]{[3]float64{-100.0, -200.0, -300.0}, 49.0}

	res := sut.ListTiles(50)

	if len(res) != 4 {
		t.Errorf("Expected %d tiles to be returned, but got %d", 4, len(res))
	}

	assertContainsTile(t, res, TileId{-1, -5})
	assertContainsTile(t, res, TileId{-2, -5})
	assertContainsTile(t, res, TileId{-2, -5})
	assertContainsTile(t, res, TileId{-2, -6})
}

func assertContainsTile(t *testing.T, tiles []TileId, expected TileId) {
	for _, tile := range tiles {
		if tile == expected {
			return
		}
	}

	t.Errorf("Expected tile (%d,%d) in result", expected.x, expected.z)
}

func TestSphereContainsReturnsTrueForSphereCenter(t *testing.T) {
	p := [3]float64{100.0, 200.0, 300.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := SphereQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0}

	res := sut.Contains(n)

	if res != true {
		t.Errorf("Expected SphereQuerier to Contain it's own center")
	}
}

func TestSphereContainsReturnsTrueForPointInRadius(t *testing.T) {
	p := [3]float64{110.0, 210.0, 310.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := SphereQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0}

	res := sut.Contains(n)

	if res != true {
		t.Errorf("Expected SphereQuerier to Contain a point within it's radius")
	}
}

func TestSphereContainsReturnsFalseForPointOutsideRadius(t *testing.T) {
	p := [3]float64{150.0, 250.0, 350.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := SphereQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0}

	res := sut.Contains(n)

	if res != false {
		t.Errorf("Expected SphereQuerier to not Contain a point outside it's radius")
	}
}

func TestSphereContainsReturnsTrueForPointOnSurface(t *testing.T) {
	p := [3]float64{100.0, 200.0, 350.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := SphereQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0}

	res := sut.Contains(n)

	if res != true {
		t.Errorf("Expected SphereQuerier to Contain a point on it's surface")
	}
}

func TestCylinderContainsReturnsTrueForCylinderCenter(t *testing.T) {
	p := [3]float64{100.0, 200.0, 300.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := CylinderQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0, math.Inf(1), math.Inf(-1)}

	res := sut.Contains(n)

	if res != true {
		t.Errorf("Expected CylinderQuerier to Contain it's own center")
	}
}

func TestCylinderContainsReturnsTrueForPointInRadius(t *testing.T) {
	p := [3]float64{110.0, 210.0, 310.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := CylinderQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0, math.Inf(1), math.Inf(-1)}

	res := sut.Contains(n)

	if res != true {
		t.Errorf("Expected CylinderQuerier to Contain a point within it's radius")
	}
}

func TestCylinderContainsReturnsFalseForPointOutsideRadius(t *testing.T) {
	p := [3]float64{150.0, 250.0, 350.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := CylinderQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0, math.Inf(1), math.Inf(-1)}

	res := sut.Contains(n)

	if res != false {
		t.Errorf("Expected CylinderQuerier to not Contain a point outside it's radius")
	}
}

func TestCylinderContainsReturnsFalseForPointAboveTop(t *testing.T) {
	p := [3]float64{100.0, 400.0, 300.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := CylinderQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0, 100, math.Inf(-1)}

	res := sut.Contains(n)

	if res != false {
		t.Errorf("Expected CylinderQuerier to not Contain a point above it's top")
	}
}

func TestCylinderContainsReturnsFalseForPointBelowBottom(t *testing.T) {
	p := [3]float64{150.0, -250.0, 350.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := CylinderQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0, math.Inf(1), -100}

	res := sut.Contains(n)

	if res != false {
		t.Errorf("Expected CylinderQuerier to not Contain a point above it's top")
	}
}

func TestCylinderContainsReturnsTrueForPointOnSideSurface(t *testing.T) {
	p := [3]float64{100.0, 10000.0, 350.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := CylinderQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0, math.Inf(1), math.Inf(-1)}

	res := sut.Contains(n)

	if res != true {
		t.Errorf("Expected CylinderQuerier to Contain a point on it's side surface")
	}
}

func TestCylinderContainsReturnsTrueForPointOnTopSurface(t *testing.T) {
	p := [3]float64{100.0, 300.0, 300.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := CylinderQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0, 100, math.Inf(-1)}

	res := sut.Contains(n)

	if res != true {
		t.Errorf("Expected CylinderQuerier to Contain a point on it's top surface")
	}
}

func TestCylinderContainsReturnsTrueForPointOnBottomSurface(t *testing.T) {
	p := [3]float64{100.0, 100.0, 300.0}
	n := &Node[any]{p, TileId{}, nil}
	sut := CylinderQuerier[any]{[3]float64{100.0, 200.0, 300.0}, 50.0, math.Inf(1), -100}

	res := sut.Contains(n)

	if res != true {
		t.Errorf("Expected CylinderQuerier to Contain a point on it's bottom surface")
	}
}
