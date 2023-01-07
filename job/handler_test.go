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
		graph.NewGraph(1),
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
		graph.NewGraph(1),
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
		graph.NewGraph(3),
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

var testGenerateCommandOrder = []struct {
	name                  string
	sortedTasks           []string
	requestTasks          []Task
	inputCommandBuffer    []Command
	expectedCommandBuffer []Command
	hasError              bool
	expectedError         error
}{
	{
		"Test with three sorted tasks should return correct command order",
		[]string{"t3", "t1", "t2"},
		[]Task{
			{Name: "t1", Command: "c1"},
			{Name: "t2", Command: "c2"},
			{Name: "t3", Command: "c3"},
		},
		make([]Command, 3),
		[]Command{
			{Name: "t3", Command: "c3"},
			{Name: "t1", Command: "c1"},
			{Name: "t2", Command: "c2"},
		},
		false,
		nil,
	},
	{
		"Test with empty sorted tasks should return empty command buffer",
		[]string{},
		[]Task{},
		make([]Command, 0),
		[]Command{},
		false,
		nil,
	},
	{
		"Test with empty sorted task not equal to the commands buffer should return specific error",
		[]string{},
		[]Task{},
		make([]Command, 2),
		[]Command{},
		true,
		commandBufferSizeErr,
	},
	{
		"Test with missing required task in sorted tasks should return specific error",
		[]string{},
		[]Task{
			{Name: "t1", Command: "c1"},
		},
		make([]Command, 0),
		[]Command{},
		true,
		requestTaskDoesNotExistErr,
	},
}

func TestGenerateCommandOrder(t *testing.T) {
	for _, tt := range testGenerateCommandOrder {
		t.Run(tt.name, func(t *testing.T) {
			err := generateCommandOrder(tt.sortedTasks, tt.requestTasks, tt.inputCommandBuffer)
			if tt.hasError {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				return
			}
			assert.Nil(t, err)
			assert.True(t, reflect.DeepEqual(tt.inputCommandBuffer, tt.expectedCommandBuffer))
		})
	}
}
