package graph

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	VertexNotFoundErr     = errors.New("vertex not found")
	GraphCycleErr         = errors.New("there is cycle in the graph")
	VertexIsNotDefinedErr = errors.New("vertex is not defined")
)

// NewGraph should be used to initialize the internal structures
func NewGraph() *DirectedGraph {
	g := DirectedGraph{}
	g.Vertices = make(map[string]*Vertex)
	g.Edges = make(map[string]*Edge)
	return &g
}

type DirectedGraph struct {
	Vertices map[string]*Vertex
	Edges    map[string]*Edge
}

type Vertex struct {
	Name string
}

type Edge struct {
	From, To *Vertex
}

// TopologicalSort is doing topological sort and returns GraphCycleErr if cycle appears
func (g *DirectedGraph) TopologicalSort() ([]string, error) {
	var sortedTasks []string
	visited := make(map[string]bool)
	processing := make(map[string]bool)

	for _, v := range g.Vertices {
		if !visited[v.Name] {
			err := g.processTask(v, &sortedTasks, visited, processing)
			if err != nil {
				return nil, err
			}
		}
	}
	return sortedTasks, nil
}

// ProcessTask is recursive function for building doing bfs and ordering vertices
// returns GraphCycleErr if cycle appears
func (g *DirectedGraph) processTask(v *Vertex, sortedTasks *[]string, visited map[string]bool, processing map[string]bool) error {
	// processing keeps track of currently processed vertexes and is used to identify cycles
	processing[v.Name] = true
	for _, edge := range g.Edges {
		if edge.From.Name == v.Name {
			if b, _ := processing[edge.To.Name]; b {
				return GraphCycleErr
			}

			// If m is not visited, then visit m.
			if !visited[edge.To.Name] {
				if err := g.processTask(edge.To, sortedTasks, visited, processing); err != nil {
					return err
				}
			}
		}
	}
	// Mark node as visited and remove from temporary state.
	visited[v.Name] = true
	processing[v.Name] = false
	// Add n to the end of L.
	*sortedTasks = append(*sortedTasks, v.Name)
	return nil
}

// Vertex retrieves a vertex by name and returns VertexNotFoundErr
func (g *DirectedGraph) Vertex(name string) (*Vertex, error) {
	if _, ok := g.Vertices[name]; !ok {
		return nil, VertexNotFoundErr
	}
	return g.Vertices[name], nil
}

func (g *DirectedGraph) AddVertex(name string) {
	v := Vertex{Name: name}
	g.Vertices[name] = &v
}

// AddEdge add edge and returns VertexNotFoundErr
func (g *DirectedGraph) AddEdge(from, to *Vertex) error {
	if from == nil || to == nil {
		return VertexIsNotDefinedErr
	}

	if err := g.validateVertexExistence(from); err != nil {
		return err
	}

	if err := g.validateVertexExistence(to); err != nil {
		return err
	}

	edge := Edge{From: from, To: to}
	edgeName := fmt.Sprintf("%s-%s", from.Name, to.Name)
	g.Edges[edgeName] = &edge
	return nil
}

func (g *DirectedGraph) validateVertexExistence(v *Vertex) error {
	if _, ok := g.Vertices[v.Name]; !ok {
		return fmt.Errorf("%w, Name: %s ", VertexNotFoundErr, v.Name)
	}
	return nil
}
