package job

import (
	"encoding/json"
	"github.com/ivanspasov99/golang-api/logging"
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

		logging.Println(r.Context(), zerolog.ErrorLevel, err.Error())
		// Deal with error here - the idea of middleware is important

		// depending on the error could be generated different status code
		w.WriteHeader(http.StatusInternalServerError)

		// write meaningful error to the user
		// write the error just for simplicity
		type ErrorResponse struct {
			Message string `json:"Message"`
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
