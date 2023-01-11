package job

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type ErrorResponseWriter struct{}

func (e ErrorResponseWriter) Header() http.Header {
	h := make(http.Header, 2)
	return h
}

func (e ErrorResponseWriter) Write(bytes []byte) (int, error) {
	return 0, fmt.Errorf("write error")
}

func (e ErrorResponseWriter) WriteHeader(statusCode int) {
}

var testWriteBash = []struct {
	name           string
	responseWriter http.ResponseWriter
	commands       []Command
	expectedBash   string
	hasError       bool
}{
	{
		"Test with commands should fail due to write error",
		ErrorResponseWriter{},
		[]Command{
			{Name: "c1", Script: "echo hello world"},
		},
		"",
		true,
	},
	{
		"Test with no commands should return only bash header",
		httptest.NewRecorder(),
		[]Command{},
		"#!/usr/bin/env bash",
		false,
	},
	{
		"Test with commands should return bash script ready for execution",
		httptest.NewRecorder(),
		[]Command{
			{Name: "c1", Script: "echo hello"},
			{Name: "c2", Script: "echo world"},
		},
		"#!/usr/bin/env bash\necho hello\necho world",
		false,
	},
}

func TestWriteBash(t *testing.T) {
	for _, tt := range testWriteBash {
		t.Run(tt.name, func(t *testing.T) {
			err := writeBash(tt.responseWriter, tt.commands)
			if tt.hasError {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			rr := tt.responseWriter.(*httptest.ResponseRecorder)
			assert.Equal(t, rr.Body.String(), tt.expectedBash)

		})
	}
}

var testWriteJSON = []struct {
	name           string
	responseWriter http.ResponseWriter
	commands       []Command
	expectedJSON   string
	hasError       bool
}{
	{
		"Test with commands should fail due to write error",
		ErrorResponseWriter{},
		[]Command{
			{Name: "c1", Script: "echo hello"},
		},
		"",
		true,
	},
	{
		"Test with commands should return json result",
		httptest.NewRecorder(),
		[]Command{
			{Name: "c1", Script: "echo hello"},
			{Name: "c2", Script: "echo world"},
		},
		`[{"name":"c1","command":"echo hello"},{"name":"c2","command":"echo world"}]`,
		false,
	},
	{
		"Test with empty commands should return empty array json",
		httptest.NewRecorder(),
		[]Command{},
		"[]",
		false,
	},
}

func TestWriteJSON(t *testing.T) {
	for _, tt := range testWriteJSON {
		t.Run(tt.name, func(t *testing.T) {
			err := writeJSON(tt.responseWriter, tt.commands)
			if tt.hasError {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)

			rr := tt.responseWriter.(*httptest.ResponseRecorder)
			assert.Equal(t, tt.expectedJSON, rr.Body.String())

		})
	}
}
