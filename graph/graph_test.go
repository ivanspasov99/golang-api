package graph

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testTopologicalSort = []struct {
	name                string
	Vertices            []*Vertex
	Edges               []*Edge
	hasError            bool
	expectedError       error
	expectedSortedArray []string
}{
	{
		"Test with four Vertices and three Edges which are connected",
		[]*Vertex{{Name: "v1"}, {Name: "v2"}, {Name: "v3"}, {Name: "v4"}},
		[]*Edge{
			{From: &Vertex{Name: "v1"}, To: &Vertex{Name: "v2"}},
			{From: &Vertex{Name: "v2"}, To: &Vertex{Name: "v3"}},
			{From: &Vertex{Name: "v3"}, To: &Vertex{Name: "v4"}},
		},
		false,
		nil,
		[]string{"v4", "v3", "v2", "v1"},
	},
	{
		"Test with cycle",
		[]*Vertex{{Name: "v1"}, {Name: "v2"}},
		[]*Edge{
			{From: &Vertex{Name: "v1"}, To: &Vertex{Name: "v2"}},
			{From: &Vertex{Name: "v2"}, To: &Vertex{Name: "v1"}},
		},
		true,
		GraphCycleErr,
		nil,
	},
}

func TestTopologicalSort(t *testing.T) {
	for _, tt := range testTopologicalSort {
		t.Run(tt.name, func(t *testing.T) {
			g := initNewTestingGraph(t, tt.Vertices, tt.Edges)
			arr, err := g.TopologicalSort()
			if tt.hasError {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedSortedArray, arr)
		})
	}
}

func initNewTestingGraph(t *testing.T, Vertices []*Vertex, Edges []*Edge) *DirectedGraph {
	g := NewGraph(2)
	for _, v := range Vertices {
		g.AddVertex(v.Name)
	}

	for _, e := range Edges {
		if err := g.AddEdge(e.From, e.To); err != nil {
			t.Fatal("Adding edge failed with", err)
			return nil
		}
	}
	return g
}

var testProcessTask = []struct {
	name                string
	Vertices            []*Vertex
	Edges               []*Edge
	hasError            bool
	expectedError       error
	expectedSortedArray []string
}{
	{
		"Test with four Vertices and three Edges which are connected starting from v1",
		[]*Vertex{{Name: "v1"}, {Name: "v2"}, {Name: "v3"}, {Name: "v4"}},
		[]*Edge{
			{From: &Vertex{Name: "v1"}, To: &Vertex{Name: "v2"}},
			{From: &Vertex{Name: "v2"}, To: &Vertex{Name: "v3"}},
			{From: &Vertex{Name: "v3"}, To: &Vertex{Name: "v4"}},
		},
		false,
		nil,
		[]string{"v4", "v3", "v2", "v1"},
	},
	{
		"Test with cycle starting from v1",
		[]*Vertex{{Name: "v1"}, {Name: "v2"}, {Name: "v3"}},
		[]*Edge{
			{From: &Vertex{Name: "v1"}, To: &Vertex{Name: "v2"}},
			{From: &Vertex{Name: "v2"}, To: &Vertex{Name: "v3"}},
			{From: &Vertex{Name: "v3"}, To: &Vertex{Name: "v1"}},
		},
		true,
		GraphCycleErr,
		nil,
	},
}

func TestProcessTask(t *testing.T) {
	for _, tt := range testProcessTask {
		t.Run(tt.name, func(t *testing.T) {
			g := initNewTestingGraph(t, tt.Vertices, tt.Edges)
			sortedTasks := make([]string, 0)
			err := g.processTask(&Vertex{Name: "v1"}, &sortedTasks, make(map[string]bool), make(map[string]bool))
			if tt.hasError {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedSortedArray, sortedTasks)
		})
	}
}

func TestVertex(t *testing.T) {
	// Create a DAG with three Vertices.
	g := NewGraph(2)
	g.AddVertex("v1")
	g.AddVertex("v2")

	vertex, err := g.Vertex("v1")
	assert.Nil(t, err)
	assert.Equal(t, "v1", vertex.Name)
	vertex, err = g.Vertex("v2")
	assert.Nil(t, err)
	assert.Equal(t, "v2", vertex.Name)

	_, err = g.Vertex("v3")
	assert.NotNil(t, err)
	assert.True(t, errors.Is(VertexNotFoundErr, err))
}

func TestAddVertex(t *testing.T) {
	g := NewGraph(2)
	assert.Len(t, g.Vertices, 0)

	g.AddVertex("v1")
	assert.Len(t, g.Vertices, 1)
	vertex, ok := g.Vertices["v1"]
	assert.True(t, ok)
	assert.Equal(t, "v1", vertex.Name)

	g.AddVertex("v2")
	assert.Len(t, g.Vertices, 2)
	vertex, ok = g.Vertices["v2"]
	assert.True(t, ok)
	assert.Equal(t, "v2", vertex.Name)
}

var testAddEdge = []struct {
	name          string
	from          *Vertex
	to            *Vertex
	hasError      bool
	expectedError error
}{
	{"Test add edge with nil from vertex should return error", nil, &Vertex{Name: "to"}, true, VertexIsNotDefinedErr},
	{"Test add edge with nil to vertex should return error", &Vertex{Name: "from"}, nil, true, VertexIsNotDefinedErr},
	{"Test add edge with existing Vertices should add Vertices", &Vertex{Name: "from"}, &Vertex{Name: "to"}, false, nil},
	{"Test add edge with non existing Vertices should return", &Vertex{Name: "non"}, &Vertex{Name: "exist"}, true, VertexNotFoundErr},
}

func TestAddEdge(t *testing.T) {
	g := NewGraph(2)
	g.AddVertex("from")
	g.AddVertex("to")

	for _, tt := range testAddEdge {
		t.Run(tt.name, func(t *testing.T) {
			err := g.AddEdge(tt.from, tt.to)
			if tt.hasError {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, tt.expectedError))
			} else {
				assert.Nil(t, err)
				edge, ok := g.Edges["from-to"]
				assert.True(t, ok)
				assert.Equal(t, "from", edge.From.Name)
				assert.Equal(t, "to", edge.To.Name)
			}
		})
	}
}
