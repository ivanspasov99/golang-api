package job

import (
	"encoding/json"
	"github.com/ivanspasov99/golang-api/pkg/graph"
	"github.com/ivanspasov99/golang-api/pkg/logging"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"net/http"
)

// HTTPTypeHandler defines handler type which will require error handler type
// Done with the idea of middleware pattern (separation of concern, chain of responsibility)
type HTTPTypeHandler func(w http.ResponseWriter, r *http.Request) error

// HandleError is function (middleware) which process errors return by job.Handle
func HandleError(h HTTPTypeHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		// Deal with error here - the idea of middleware is important
		logging.Println(r.Context(), zerolog.ErrorLevel, err.Error())

		type ErrorResponse struct {
			Message string `json:"Message"`
		}

		// depending on the error could be generated different status code, different responses, server reaction as alerting etc.
		w.Header().Set("Content-Type", "application/json")
		switch err {
		case graph.GraphCycleErr:
			w.WriteHeader(http.StatusBadRequest)
			err = errors.Errorf("Please evaluate tasks. Processing feedback: %s", err.Error())
		case graph.VertexNotFoundErr:
			w.WriteHeader(http.StatusBadRequest)
			err = errors.Errorf("Please evaluate required tasks as one of the defined one is not existing. Processing feedback: %s", err.Error())
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

		eR := ErrorResponse{Message: err.Error()}

		b, err := json.Marshal(eR)
		if err != nil {
			logging.Println(r.Context(), zerolog.ErrorLevel, err.Error())
			// ignore error just for simplicity
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(b)

		// Integrate Sentry for example to notify us by slack for error
		// err := sentry.Init(sentry.ClientOptions{
		//		Dsn: "https://examplePublicKey@o0.ingest.sentry.io/0",
		//		// Enable printing of SDK debug messages.
		//		// Useful when getting started or trying to figure something out.
		//		Debug: true,
		//	})
		//	if err != nil {
		//		log.Fatalf("sentry.Init: %s", err)
		//	}
		//	// Flush buffered events before the program terminates.
		//	// Set the timeout to the maximum duration the program can afford to wait.
		//	defer sentry.Flush(2 * time.Second)
	})
}
