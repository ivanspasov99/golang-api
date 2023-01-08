package main

import (
	"github.com/ivanspasov99/golang-api/pkg/job"
	"github.com/ivanspasov99/golang-api/pkg/logging"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	http.HandleFunc("/job", logging.DecorateHeader(job.HandleError(job.Handle)))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Msg(err.Error())
	}
}
