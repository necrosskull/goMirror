package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"goMirror/cfg"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	config, err := cfg.LoadConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Error loading cfg")
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("method", r.Method).Str("url", r.URL.String()).Msg("Incoming request")

		query := strings.TrimPrefix(r.URL.Path, "/")
		targetURL := fmt.Sprintf("%s/%s", config.Url, query)

		req, err := http.NewRequest(r.Method, targetURL, r.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error creating request")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		req.Header = make(http.Header)
		for key, values := range r.Header {
			req.Header[key] = values
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("Error making request to API")
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Error().Err(err).Msg("Error closing response body")
				return
			}
		}(resp.Body)

		for key, values := range resp.Header {
			w.Header()[key] = values
		}

		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error copying response body")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info().Int("status", resp.StatusCode).Msg("Outgoing response with")

	})

	err = http.ListenAndServe(":8085", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
		return
	}
}
