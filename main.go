package main

import (
	"fmt"
	"goMirror/cfg"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	config, err := cfg.LoadConfig()

	if err != nil {
		log.Fatal().Err(err).Msg("Error loading cfg")
		return
	}

	router := gin.Default()

	router.Any("/*path", func(c *gin.Context) {
		query := strings.TrimPrefix(c.Request.URL.Path, "/")

		// Получаем параметры из query
		queryParams := c.Request.URL.Query()
		queryString := queryParams.Encode()

		// Если есть параметры, добавляем их к целевому URL
		if queryString != "" {
			query = fmt.Sprintf("%s?%s", query, queryString)
		}

		targetURL := fmt.Sprintf("%s/%s", config.Url, query)

		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error creating request")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		req.Header = make(http.Header)
		for key, values := range c.Request.Header {
			req.Header[key] = values
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error().Err(err).Msg("Error making request to API")
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Error().Err(err).Msg("Error closing response body")
			}
		}()

		for key, values := range resp.Header {
			c.Header(key, values[0])
		}

		c.Status(resp.StatusCode)

		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error copying response body")
			return
		}
	})
	log.Info().Msg("Server started on http://127.0.0.1:" + config.Port)

	err = router.Run(":" + config.Port)
	if err != nil {
		log.Fatal().Err(err).Msg("Error starting server")
		return
	}
}
