package job

import (
	"encoding/json"
	"github.com/ivanspasov99/golang-api/graph"
	"io"
	"log"
	"net/http"
)

type Job struct {
	Tasks []Task `json:"tasks"`
}

// TODO should with command
type Task struct {
	Name     string   `json:"name"`
	Required []string `json:"requires"`
}

type Order struct {
	Commands []string `json:"commands"`
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

// HandleJob processes Job which tasks are being sorted in required order
// Internally it is using graph.DirectedGraph which is doing sorting in linear complexity
// A Job is a collection of tasks, where each Task has a name and a shell command. Tasks may
// depend on other tasks and require that those are executed beforehand.
func HandleJob(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}

	j := Job{}
	// handle error
	json.Unmarshal(b, &j)

	// pass the graph for two reason
	// g would be big object as it is possible to have thousands of tasks
	// therefore using stack instead of the heap is smarter (if g is declared in populateGraph we are going to use the heap)
	// golang idiom - pass abstractions return concretions (easier testing also as interface is mockable)
	g := graph.NewGraph()
	if err := populateGraph(j.Tasks, g); err != nil {
		log.Fatalf("Populate graph. Err: %s", err)
		return
	}

	arr, err := g.TopologicalSort()
	if err != nil {
		log.Fatalf("Topological sort error. Err: %s", err)
		return
	}

	o := Order{Commands: arr}
	jsonResp, err := json.Marshal(o)
	if err != nil {
		// return error
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Fatalf("Error happened write. Err: %s", err)
		return
	}
	return
}
