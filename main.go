package main

import (
	"github.com/ivanspasov99/golang-api/job"
	"github.com/ivanspasov99/golang-api/logging"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	http.HandleFunc("/", logging.DecorateHeader(job.HandleJob))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Msg(err.Error())
	}
}