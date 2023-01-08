package job

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ResponseWriter func type is an adapter (interface like function) to allow the use of ordinary functions as Job response writers.
type ResponseWriter func(http.ResponseWriter, []Command) error

const bash = "bash"

func writeBash(w http.ResponseWriter, commands []Command) error {
	arr := make([]string, len(commands)+1)

	// we could identify where bash is installed
	// or use community dependency for generating bash script
	bashHeader := "#!/usr/bin/env bash"

	arr[0] = bashHeader

	for i, command := range commands {
		arr[i+1] = command.Command
	}

	s := strings.Join(arr, "\n")

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(s))
	if err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, commands []Command) error {
	jsonResp, err := json.Marshal(commands)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResp)
	if err != nil {
		return err
	}
	return nil
}

func getJobModeWriter(r *http.Request) ResponseWriter {
	mode := r.URL.Query().Get("mode")
	switch strings.ToLower(mode) {
	case bash:
		return writeBash
	default:
		return writeJSON
	}
}
