package octree

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// From https://github.com/benbjohnson/testing
// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestBoxContainsPoints(t *testing.T) {
	b := Box{
		min: Vector3f{0, 0, 0},
		max: Vector3f{1, 1, 1},
	}

	equals(t, true, b.ContainsPoint(&Vector3f{0, 0, 0}))
	equals(t, true, b.ContainsPoint(&Vector3f{1, 0, 0}))
	equals(t, true, b.ContainsPoint(&Vector3f{0, 0, 1}))
	equals(t, true, b.ContainsPoint(&Vector3f{0.5, 0.5, 0.5}))
	equals(t, true, b.ContainsPoint(&Vector3f{-0, 0, 0}))
	equals(t, false, b.ContainsPoint(&Vector3f{-0.000001, 0.5, 0.5}))
	equals(t, false, b.ContainsPoint(&Vector3f{0.5, -0.000001, 0.5}))
	equals(t, false, b.ContainsPoint(&Vector3f{0.5, 0.5, -0.000001}))
	equals(t, false, b.ContainsPoint(&Vector3f{1.000001, 0.5, 0.5}))
	equals(t, false, b.ContainsPoint(&Vector3f{0.5, 1.000001, 0.5}))
	equals(t, false, b.ContainsPoint(&Vector3f{0.5, 0.5, 1.000001}))
}

func TestBoxContainsBox(t *testing.T) {
	b := Box{
		min: Vector3f{0, 0, 0},
		max: Vector3f{1, 1, 1},
	}
	var b2 Box

	// contains equal box, symmetrically
	b2 = Box{Vector3f{0, 0, 0}, Vector3f{1, 1, 1}}
	equals(t, true, b2.Contains(&b))
	equals(t, true, b.Contains(&b2))
	equals(t, true, b2.IsContainedIn(&b))
	equals(t, true, b.IsContainedIn(&b2))

	// contained on edge
	b2 = Box{Vector3f{0, 0, 0}, Vector3f{0.5, 1, 1}}
	equals(t, true, b.Contains(&b2))
	equals(t, false, b2.Contains(&b))
	equals(t, false, b.IsContainedIn(&b2))
	equals(t, true, b2.IsContainedIn(&b))

	// contained away from edges
	b2 = Box{Vector3f{0.1, 0.1, 0.1}, Vector3f{0.9, 0.9, 0.9}}
	equals(t, true, b.Contains(&b2))
	equals(t, false, b2.Contains(&b))
	equals(t, false, b.IsContainedIn(&b2))
	equals(t, true, b2.IsContainedIn(&b))

	// 1 corner inside
	b2 = Box{Vector3f{-0.1, -0.1, -0.1}, Vector3f{0.9, 0.9, 0.9}}
	equals(t, false, b.Contains(&b2))
	equals(t, false, b2.Contains(&b))
	equals(t, false, b.IsContainedIn(&b2))
	equals(t, false, b2.IsContainedIn(&b))

	b2 = Box{Vector3f{0.9, 0.9, 0.9}, Vector3f{1.1, 1.1, 1.1}}
	equals(t, false, b.Contains(&b2))
	equals(t, false, b2.Contains(&b))
	equals(t, false, b.IsContainedIn(&b2))
	equals(t, false, b2.IsContainedIn(&b))
}

func TestBoxIntersectsBox(t *testing.T) {
	b := Box{
		min: Vector3f{0, 0, 0},
		max: Vector3f{1, 1, 1},
	}
	var b2 Box

	// not intersecting box above or below in each dimension
	b2 = Box{Vector3f{1.1, 0, 0}, Vector3f{2, 1, 1}}
	equals(t, false, b.Intersects(&b2))
	b2 = Box{Vector3f{-1, 0, 0}, Vector3f{-0.1, 1, 1}}
	equals(t, false, b.Intersects(&b2))
	b2 = Box{Vector3f{0, 1.1, 0}, Vector3f{1, 2, 1}}
	equals(t, false, b.Intersects(&b2))
	b2 = Box{Vector3f{0, -1, 0}, Vector3f{1, -0.1, 1}}
	equals(t, false, b.Intersects(&b2))
	b2 = Box{Vector3f{0, 0, 1.1}, Vector3f{1, 1, 2}}
	equals(t, false, b.Intersects(&b2))
	b2 = Box{Vector3f{0, 0, -1}, Vector3f{1, 1, -0.1}}
	equals(t, false, b.Intersects(&b2))

	// intersects equal box, symmetrically
	b2 = Box{Vector3f{0, 0, 0}, Vector3f{1, 1, 1}}
	equals(t, true, b.Intersects(&b2))
	equals(t, true, b2.Intersects(&b))

	// intersects containing and contained
	b2 = Box{Vector3f{0.1, 0.1, 0.1}, Vector3f{0.9, 0.9, 0.9}}
	equals(t, true, b.Intersects(&b2))
	equals(t, true, b2.Intersects(&b))

	// intersects partial containment on each corner
	b2 = Box{Vector3f{0.9, 0.9, 0.9}, Vector3f{2, 2, 2}}
	equals(t, true, b.Intersects(&b2))
	b2 = Box{Vector3f{-1, 0.9, 0.9}, Vector3f{0.1, 2, 2}}
	equals(t, true, b.Intersects(&b2))
	b2 = Box{Vector3f{0.9, -1, 0.9}, Vector3f{2, 0.1, 2}}
	equals(t, true, b.Intersects(&b2))
	b2 = Box{Vector3f{-1, -1, 0.9}, Vector3f{0.1, 0.1, 2}}
	equals(t, true, b.Intersects(&b2))
	b2 = Box{Vector3f{0.9, 0.9, -1}, Vector3f{2, 2, 0.1}}
	equals(t, true, b.Intersects(&b2))
	b2 = Box{Vector3f{-1, 0.9, -1}, Vector3f{0.1, 2, 0.1}}
	equals(t, true, b.Intersects(&b2))
	b2 = Box{Vector3f{0.9, -1, -1}, Vector3f{2, 0.1, 0.1}}
	equals(t, true, b.Intersects(&b2))
	b2 = Box{Vector3f{-1, -1, -1}, Vector3f{0.1, 0.1, 0.1}}
	equals(t, true, b.Intersects(&b2))

	// intersects 'beam'; where no corners inside
	// other but some contained
	b2 = Box{Vector3f{-1, 0.1, 0.1}, Vector3f{2, 0.9, 0.9}}
	equals(t, true, b.Intersects(&b2))
	b2 = Box{Vector3f{0.1, -1, 0.1}, Vector3f{0.9, 2, 0.9}}
	equals(t, true, b.Intersects(&b2))
	b2 = Box{Vector3f{0.1, 0.1, -1}, Vector3f{0.9, 0.9, 2}}
	equals(t, true, b.Intersects(&b2))
}

func TestInitializesRoot(t *testing.T) {
	o := CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})

	//o.Add(99, Vector3f{10, 0, 0})

	equals(t, true, o.root.point == nil)
	equals(t, Vector3f{0, 0, 0}, o.root.box.min)
	equals(t, Vector3f{1, 1, 1}, o.root.box.max)
	equals(t, false, o.root.hasChildren)
	equals(t, 0, len(o.root.children))
}

func TestInsertsContainedElements(t *testing.T) {
	o := CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})

	equals(t, true, o.Add(99, Vector3f{1.00000000001, 1, 1}) == nil)
	equals(t, false, o.root.hasChildren)
	equals(t, true, o.root.children == nil)
	equals(t, true, o.root.point == nil)

	equals(t, true, o.Add(99, Vector3f{-0.0000000001, 0, 0}) == nil)
	equals(t, false, o.root.hasChildren)
	equals(t, true, o.root.children == nil)
	equals(t, true, o.root.point == nil)

	equals(t, false, o.Add(88, Vector3f{0.5, 0, 0}) == nil)
	equals(t, false, o.root.hasChildren)
	equals(t, true, o.root.children == nil)
	equals(t, false, o.root.point == nil)
}

func TestEqualPointsSubdivide(t *testing.T) {
	o := CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})

	o.Add(1, Vector3f{0, 0, 0})
	o.Add(1, Vector3f{0, 0, 0})
	equals(t, false, o.root.hasChildren)
	equals(t, true, o.root.children == nil)
	equals(t, true, *o.root.point == Vector3f{0, 0, 0})
	o.Add(1, Vector3f{1, 1, 1})
	equals(t, true, o.root.hasChildren)
	equals(t, false, o.root.children == nil)
	equals(t, true, o.root.point == nil)
}

func TestRetrievesElementsIn(t *testing.T) {
	o := CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})

	o.Add(11, Vector3f{0, 0, 0})
	// contains point
	equals(t, 1, len(o.ElementsIn(Box{Vector3f{-1, -1, -1}, Vector3f{0.1, 0.1, 0.1}})))
	// 0 size at point
	equals(t, 1, len(o.ElementsIn(Box{Vector3f{0, 0, 0}, Vector3f{0, 0, 0}})))
	// contains box
	equals(t, 1, len(o.ElementsIn(Box{Vector3f{-1, -1, -1}, Vector3f{2, 2, 2}})))

	// coincident point
	o.Add(12, Vector3f{0, 0, 0})
	// contains point
	equals(t, 2, len(o.ElementsIn(Box{Vector3f{-1, -1, -1}, Vector3f{0.1, 0.1, 0.1}})))
	// 0 size at point
	equals(t, 2, len(o.ElementsIn(Box{Vector3f{0, 0, 0}, Vector3f{0, 0, 0}})))
	// contains box
	equals(t, 2, len(o.ElementsIn(Box{Vector3f{-1, -1, -1}, Vector3f{2, 2, 2}})))

	// non-coincident point
	o.Add(2, Vector3f{1, 1, 1})
	equals(t, true, o.root.hasChildren)
	equals(t, false, o.root.children == nil)
	equals(t, 8, len(o.root.children))
	equals(t, true, o.root.point == nil)

	// contains point
	equals(t, 2, len(o.ElementsIn(Box{Vector3f{-1, -1, -1}, Vector3f{0.1, 0.1, 0.1}})))
	// 0 size at point
	equals(t, 2, len(o.ElementsIn(Box{Vector3f{0, 0, 0}, Vector3f{0, 0, 0}})))
	// contains box
	equals(t, 3, len(o.ElementsIn(Box{Vector3f{-1, -1, -1}, Vector3f{2, 2, 2}})))

	// fresh octree
	o = CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})
	equals(t, false, o.root.hasChildren)

	o.Add(11, Vector3f{0.4, 0.4, 0.4})
	equals(t, 1, len(o.ElementsIn(Box{Vector3f{-1, -1, -1}, Vector3f{2, 2, 2}})))
	equals(t, 1, len(o.ElementsIn(Box{Vector3f{0.4, 0.4, 0.4}, Vector3f{0.6, 0.6, 0.6}})))
	equals(t, 1, len(o.ElementsIn(Box{Vector3f{-1, 0.4, 0.4}, Vector3f{1, 0.6, 0.6}})))

	o.Add(12, Vector3f{0.68, 0.69, 0.7})
	equals(t, 1, len(o.ElementsIn(Box{Vector3f{-1, 0.4, 0.4}, Vector3f{1, 0.6, 0.6}})))
	equals(t, 2, len(o.ElementsIn(Box{Vector3f{-1, -1, -1}, Vector3f{1, 1, 1}})))

	// add coincident point in octree
	o.Add(13, Vector3f{0.68, 0.69, 0.7})
	equals(t, 3, len(o.ElementsIn(Box{Vector3f{-1, -1, -1}, Vector3f{1, 1, 1}})))
	equals(t, 2, len(o.ElementsIn(Box{Vector3f{0.68, 0.69, 0.7}, Vector3f{0.68, 0.69, 0.7}})))
	equals(t, 1, len(o.ElementsIn(Box{Vector3f{0.35, 0.35, 0.35}, Vector3f{0.45, 0.45, 0.45}})))

	o.Add(14, Vector3f{0.1, 0.9, 0.1})

	// values
	equals(t, 11, o.ElementsIn(Box{Vector3f{0.35, 0.35, 0.35}, Vector3f{0.45, 0.45, 0.45}})[0])
	equals(t, 12, o.ElementsIn(Box{Vector3f{0.65, 0.65, 0.65}, Vector3f{0.75, 0.75, 0.75}})[0])
	equals(t, 13, o.ElementsIn(Box{Vector3f{0.65, 0.65, 0.65}, Vector3f{0.75, 0.75, 0.75}})[1])
}

func TestRetrievesElementsAt(t *testing.T) {
	o := CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})

	o.Add(11, Vector3f{0.1, 0.1, 0.1})
	// finds element at point
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))

	// coincident point with different value
	o.Add(12, Vector3f{0.1, 0.1, 0.1})

	// finds elements at point
	equals(t, 2, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[1]))

	// finds elements at point after subdivision
	o.Add(13, Vector3f{0.7, 0.7, 0.7})
	equals(t, 2, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[1]))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.7, 0.7, 0.7})))
	equals(t, 13, (o.ElementsAt(Vector3f{0.7, 0.7, 0.7})[0]))

	// finds elements at point after multiple subdivisions
	o.Add(14, Vector3f{0.1, 0.1, 0.2})
	equals(t, 2, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[1]))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.7, 0.7, 0.7})))
	equals(t, 13, (o.ElementsAt(Vector3f{0.7, 0.7, 0.7})[0]))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.2})))
	equals(t, 14, (o.ElementsAt(Vector3f{0.1, 0.1, 0.2})[0]))
}

func TestRemovesElements(t *testing.T) {
	o := CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})

	// removes element
	o.Add(11, Vector3f{0.1, 0.1, 0.1})
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, true, o.Remove(11))
	equals(t, 0, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))

	// remove correct element
	o.Add(11, Vector3f{0.1, 0.1, 0.1})
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	o.Add(12, Vector3f{0.1, 0.1, 0.1})
	equals(t, 2, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[1]))
	equals(t, true, o.Remove(11))
	equals(t, false, o.Remove(11))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, true, o.Remove(12))
	equals(t, false, o.Remove(12))
	equals(t, 0, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))

	o.Add(11, Vector3f{0.1, 0.1, 0.1})
	o.Add(12, Vector3f{0.1, 0.1, 0.1})
	o.Add(13, Vector3f{0.7, 0.7, 0.7})
	equals(t, 2, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[1]))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.7, 0.7, 0.7})))
	equals(t, 13, (o.ElementsAt(Vector3f{0.7, 0.7, 0.7})[0]))
	equals(t, true, o.Remove(11))
	equals(t, false, o.Remove(11))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.7, 0.7, 0.7})))
	equals(t, 13, (o.ElementsAt(Vector3f{0.7, 0.7, 0.7})[0]))
	equals(t, true, o.Remove(12))
	equals(t, false, o.Remove(12))
	equals(t, 0, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.7, 0.7, 0.7})))
	equals(t, 13, (o.ElementsAt(Vector3f{0.7, 0.7, 0.7})[0]))
	equals(t, true, o.Remove(13))
	equals(t, false, o.Remove(13))
	equals(t, 0, len(o.ElementsAt(Vector3f{0.7, 0.7, 0.7})))
}

func TestRemovesElementsUsing(t *testing.T) {
	o := CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})

	// removes element using node ref
	node11 := o.Add(11, Vector3f{0.1, 0.1, 0.1})
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, true, o.RemoveUsing(11, node11))
	equals(t, 0, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))

	// removes element after subdivision using node ref
	node11 = o.Add(11, Vector3f{0.1, 0.1, 0.1})
	node12 := o.Add(12, Vector3f{0.1, 0.1, 0.1})
	node13 := o.Add(13, Vector3f{0.7, 0.7, 0.7})
	node13b := o.Add(13, Vector3f{0.1, 0.1, 0.2})
	equals(t, 2, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[1]))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.7, 0.7, 0.7})))
	equals(t, true, o.RemoveUsing(13, node13))
	equals(t, false, o.RemoveUsing(13, node13))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.2})))
	equals(t, 13, (o.ElementsAt(Vector3f{0.1, 0.1, 0.2})[0]))
	equals(t, true, o.RemoveUsing(13, node13b))
	equals(t, false, o.RemoveUsing(13, node13b))
	equals(t, 0, len(o.ElementsAt(Vector3f{0.7, 0.7, 0.7})))
	equals(t, 2, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 11, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[1]))
	equals(t, true, o.RemoveUsing(11, node11))
	equals(t, false, o.RemoveUsing(11, node11))
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	equals(t, 12, (o.ElementsAt(Vector3f{0.1, 0.1, 0.1})[0]))
	equals(t, true, o.RemoveUsing(12, node12))
	equals(t, false, o.RemoveUsing(12, node12))
	equals(t, 0, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
}

func TestClearTree(t *testing.T) {
	o := CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})
	equals(t, 0, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
	o.Add(11, Vector3f{0.1, 0.1, 0.1})
	equals(t, 1, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))

	o.Clear()
	equals(t, 0, len(o.ElementsAt(Vector3f{0.1, 0.1, 0.1})))
}
