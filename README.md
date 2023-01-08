# Summary

## Job Handler
Job Handler handles every job in separate goroutine, it is highly possible real big data scenario so 
performance is crucial. Therefore, it is taken into account and time complexity is linear - Graph Implementation O(n + e).

Encoding/Decoding special symbols use-cases are not taken into account

## Testing
Testing is created using Table Driven Testing over BDT (Behavior Driven Testing). Output could be improved when test fails, as it would 
bring big value in debugging faster. 

## Logging Package
Package encapsulate productive json requirement logging which is required by a lot of analysing log tools

Logging package could be extended with dynamic logging and log level state which represent the option to change the level of logging (debug, warn, info, error)
This help in generating fewer logs when not needed and set more logs when problem arise for debugging purposes
This could be implemented through a configmap in the which is deployed in k8s cluster for examples separately from
the application, then you can consume/read it as `env` variable in the code  

## Graph Algorithm
**It is better to use already implemented packages which are community adopted and tested**, but I have decided to refresh my skills a little bit

**It is best to be implemented using generics as now it is very limited to one type/struct etc**

**Time Complexity - O(n + e)**

The implementation is using maps which does no guarantee order of the topological sorting. It can be 
implemented with arrays or can be improved with following custom map key logic
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

### Security
Security not part of the task