# Go Octree

This is an implementation of an [Octree](https://en.wikipedia.org/wiki/Octree) in [Go](https://golang.org/). It is based on [this](https://github.com/raywenderlich/swift-algorithm-club/tree/master/Octree) implementation in Swift and was created to help teach myself Go.

### Examples

```
// Create
oct := CreateOctree(Vector3f{0, 0, 0}, Vector3f{1, 1, 1})

// Add element at point
oct.Add(1, Vector3f{0.1, 0.2, 0.3})
oct.Add(2, Vector3f{0.2, 0.3, 0.4})
oct.Add(3, Vector3f{0.3, 0.4, 0.5})
node4 := oct.Add(4, Vector3f{0.3, 0.4, 0.5}) // save for removal later

// Retrieval at point
oct.ElementsAt(Vector3f{0.1, 0.2, 0.3}) // [1]
oct.ElementsAt(Vector3f{0.2, 0.3, 0.4}) // [2]
oct.ElementsAt(Vector3f{0.3, 0.4, 0.5}) // [3 4]

// Retrieval in box
oct.ElementsIn(Box{Vector3f{0.1, 0.2, 0.3}, Vector3f{0.2, 0.3, 0.4}}) // [1 2]

// Remove first of element in tree (slower)
oct.Remove(1) // true

// Remove first of element within node (faster)
oct.RemoveUsing(4, node4) // true

// Clear contents
oct.Clear() // true
```

#### License

[MIT](https://opensource.org/licenses/MIT)
