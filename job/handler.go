package job

import (
	"github.com/ivanspasov99/golang-api/graph"
)

type Job struct {
}

type Task struct {
	Name     string
	Required []string
}

type Graph interface {
	TopologicalSort() ([]string, error)
	Vertex(name string) (*graph.Vertex, error)
	AddVertex(name string)
	AddEdge(from, to *graph.Vertex) error
}

func populateGraph(tasks []Task, g Graph) error {
	// Add a vertex for each task
	for _, t := range tasks {
		g.AddVertex(t.Name)
	}

	// Add an edge for each required
	for _, t := range tasks {
		for _, r := range t.Required {
			from, err := g.Vertex(t.Name)
			if err != nil {
				return err
			}
			to, err := g.Vertex(r)
			if err != nil {
				return err
			}
			if err := g.AddEdge(from, to); err != nil {
				return err
			}
		}
	}
	return nil
}
