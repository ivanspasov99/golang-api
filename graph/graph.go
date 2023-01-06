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

func NewGraph() *graph {
	g := graph{}
	g.vertices = make(map[string]*Vertex)
	g.edges = make(map[string]*Edge)
	return &g
}

type Vertex struct {
	Name string
}

type Edge struct {
	From, To *Vertex
}

type graph struct {
	vertices map[string]*Vertex
	edges    map[string]*Edge
}

// TopologicalSort is doing topological sort and returns GraphCycleErr if cycle appears
func (g *graph) TopologicalSort() ([]string, error) {
	var sortedTasks []string
	visited := make(map[string]bool)
	processing := make(map[string]bool)

	for _, v := range g.vertices {
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
func (g *graph) processTask(v *Vertex, sortedTasks *[]string, visited map[string]bool, processing map[string]bool) error {
	// processing keeps track of currently processed vertexes and is used to identify cycles
	processing[v.Name] = true
	for _, edge := range g.edges {
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
func (g *graph) Vertex(name string) (*Vertex, error) {
	if _, ok := g.vertices[name]; !ok {
		return nil, VertexNotFoundErr
	}
	return g.vertices[name], nil
}

func (g *graph) AddVertex(name string) {
	v := Vertex{Name: name}
	g.vertices[name] = &v
}

// AddEdge add edge and returns VertexNotFoundErr
func (g *graph) AddEdge(from, to *Vertex) error {
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
	g.edges[edgeName] = &edge
	return nil
}

func (g *graph) validateVertexExistence(v *Vertex) error {
	if _, ok := g.vertices[v.Name]; !ok {
		return fmt.Errorf("%w, Name: %s ", VertexNotFoundErr, v.Name)
	}
	return nil
}
