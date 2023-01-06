package logging

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDecorateShouldHaveSetRequestIdContext(t *testing.T) {
	checkHandlerFunction := func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			t.Error("request id is not set in context")
		}
	}

	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	Decorate(checkHandlerFunction).ServeHTTP(rr, req)
}

func TestDecorateHeaderShouldHaveSetRequestIdInHeaderAndContext(t *testing.T) {
	checkHandlerFunction := func(w http.ResponseWriter, r *http.Request) {
		id, ok := r.Context().Value(requestIDKey).(string)
		if !ok {
			t.Error("request id is not set in context")
		}

		hId := w.Header().Get(RequestIdHeader)
		if hId != id {
			t.Errorf("request id in header %s does not match request id in context %s", hId, id)
		}
	}

	rr := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	DecorateHeader(checkHandlerFunction).ServeHTTP(rr, req)
}
