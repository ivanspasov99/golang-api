package job

import (
	"encoding/json"
	"fmt"
	"github.com/ivanspasov99/golang-api/graph"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

var (
	requestTaskDoesNotExistErr = errors.New("request task does not exist in the sorted ones")

	commandBufferSizeErr = errors.New("sorted tasks are more than the passed buffer size")
)

type Job struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Name     string   `json:"name"`
	Command  string   `json:"command"`
	Required []string `json:"requires"`
}

type Command struct {
	Name   string `json:"name"`
	Script string `json:"command"`
}

type Graph interface {
	TopologicalSort() ([]string, error)
	Vertex(name string) (*graph.Vertex, error)
	AddVertex(name string)
	AddEdge(from, to *graph.Vertex) error
}

// Handle processes Job which tasks are being sorted in required order
// Internally it is using graph.DirectedGraph which is doing sorting in linear complexity
// A Job is a collection of tasks, where each Task has a name and a shell command. Tasks may
// depend on other tasks and require that those are executed beforehand.
// returns
func Handle(w http.ResponseWriter, r *http.Request) error {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	j := Job{}
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	g := graph.NewGraph(len(j.Tasks))
	if err := populateGraph(j.Tasks, g); err != nil {
		return err
	}

	sortedArr, err := g.TopologicalSort()
	if err != nil {
		return err
	}

	commandBuffer := make([]Command, len(sortedArr))
	if err := generateCommandOrder(sortedArr, j.Tasks, commandBuffer); err != nil {
		return err
	}

	writeResponse := getJobModeWriter(r)
	if err := writeResponse(w, commandBuffer); err != nil {
		return err
	}
	return nil
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

// generateCommandOrder populates commandBuffer with ordered commands based on sorted tasks
func generateCommandOrder(sortedTasks []string, requestTasks []Task, commandBuffer []Command) error {
	if len(sortedTasks) != len(commandBuffer) {
		return commandBufferSizeErr
	}

	// use map for constant access
	tmp := make(map[string]int)
	for i, v := range sortedTasks {
		tmp[v] = i
	}

	for _, t := range requestTasks {
		v, ok := tmp[t.Name]
		if !ok {
			return fmt.Errorf("%w, task: %s", requestTaskDoesNotExistErr, t.Name)
		}
		commandBuffer[v] = Command{Name: t.Name, Script: t.Command}
	}
	return nil
}
