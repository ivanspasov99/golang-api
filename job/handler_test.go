package job

import (
	"github.com/ivanspasov99/golang-api/graph"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func initErrorGraph(edge error, vertex error, topologicalError error) *MockGraph {
	return &MockGraph{
		edge:             edge,
		vertex:           vertex,
		topologicalError: topologicalError,
	}
}

type MockGraph struct {
	edge             error
	vertex           error
	topologicalError error
}

func (mg MockGraph) TopologicalSort() ([]string, error) {
	return nil, mg.topologicalError
}

func (mg MockGraph) Vertex(name string) (*graph.Vertex, error) {
	return nil, mg.vertex
}

func (mg MockGraph) AddVertex(name string) {}

func (mg MockGraph) AddEdge(from, to *graph.Vertex) error {
	return mg.edge
}

var testGenerateGraph = []struct {
	name          string
	tasks         []Task
	input         Graph
	expected      *graph.DirectedGraph
	hasError      bool
	expectedError error
}{
	{
		"Test should finish with vertex error",
		[]Task{
			{Name: "task2"},
			{Name: "task1", Required: []string{"task2"}},
		},
		initErrorGraph(nil, graph.VertexNotFoundErr, nil),
		nil,
		true,
		graph.VertexNotFoundErr,
	},
	{
		"Test should finish with edge error",
		[]Task{
			{Name: "task2"},
			{Name: "task1", Required: []string{"task2"}},
		},
		initErrorGraph(graph.VertexIsNotDefinedErr, nil, nil),
		nil,
		true,
		graph.VertexIsNotDefinedErr,
	},
	{
		"Test with single task",
		[]Task{
			{Name: "task1"},
		},
		graph.NewGraph(),
		&graph.DirectedGraph{
			Vertices: map[string]*graph.Vertex{
				"task1": {Name: "task1"},
			},
			Edges: map[string]*graph.Edge{},
		},
		false,
		nil,
	},
	{
		"Test with cycle should not return error",
		[]Task{
			{Name: "task1", Required: []string{"task1"}},
		},
		graph.NewGraph(),
		&graph.DirectedGraph{
			Vertices: map[string]*graph.Vertex{
				"task1": {Name: "task1"},
			},
			Edges: map[string]*graph.Edge{
				"task1-task1": {
					From: &graph.Vertex{Name: "task1"},
					To:   &graph.Vertex{Name: "task1"},
				},
			},
		},
		false,
		nil,
	},
	{
		"Test with multiple vertexes and edges",
		[]Task{
			{Name: "task1", Required: []string{"task2"}},
			{Name: "task2", Required: []string{"task3"}},
			{Name: "task3"},
		},
		graph.NewGraph(),
		&graph.DirectedGraph{
			Vertices: map[string]*graph.Vertex{
				"task1": {Name: "task1"},
				"task2": {Name: "task2"},
				"task3": {Name: "task3"},
			},
			Edges: map[string]*graph.Edge{
				"task1-task2": {
					From: &graph.Vertex{Name: "task1"},
					To:   &graph.Vertex{Name: "task2"},
				},
				"task2-task3": {
					From: &graph.Vertex{Name: "task2"},
					To:   &graph.Vertex{Name: "task3"},
				},
			},
		},
		false,
		nil,
	},
}

func TestGenerateGraph(t *testing.T) {
	for _, tt := range testGenerateGraph {
		t.Run(tt.name, func(t *testing.T) {
			err := populateGraph(tt.tasks, tt.input)
			if tt.hasError {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				return
			}
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(tt.input, tt.expected))
		})
	}
}
