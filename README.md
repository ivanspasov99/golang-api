# Summary


## Graph Algorithm
- The implementation is using maps which does no guarantee order of the topological sorting. It can be 
implemented with arrays or can be improved with following logic
```go
// key will be used as map key
type key struct {
	name string
}

// graph will look like
type graph struct {
	vertices map[key]*Vertex
	edges    map[key]*Edge
}

// comparison function
func (k key) Less(other key) bool {
	// Return true if k should be sorted before the other.
	// Return false otherwise.
}
```