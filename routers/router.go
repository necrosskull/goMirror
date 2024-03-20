package routers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

var (
	targetURL    string
	reverseProxy *httputil.ReverseProxy
)

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		modifyRequest(req)
	}

	proxy.ModifyResponse = modifyResponse()
	proxy.ErrorHandler = errorHandler()
	return proxy, nil
}

func modifyRequest(req *http.Request) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing target URL")
		return
	}
	req.Host = parsedURL.Hostname()
	req.Header.Set("Host", parsedURL.Hostname())
}

func errorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		log.Error().Err(err).Msg("Error in reverse proxy")
	}
}

func modifyResponse() func(*http.Response) error {
	return func(resp *http.Response) error {
		return nil
	}
}

func Router(r *gin.Engine) {
	r.Any("/*path", proxy)
	r.NoRoute(proxy)
}

func Init(targetURL string) {
	proxy, err := NewProxy(targetURL)
	log.Info().Str("targetURL", targetURL).Msg("Creating reverse proxy")
	if err != nil {
		log.Error().Err(err).Msg("Error creating reverse proxy")
		return
	}
	reverseProxy = proxy
}

func proxy(c *gin.Context) {
	// log.Info().
	// 	Str("CF-Connecting-IP", c.Request.Header.Get("CF-Connecting-IP")).
	// 	Str("ua", c.Request.UserAgent()).
	// 	Str("method", c.Request.Method).
	// 	Str("path", c.Request.URL.Path).
	// 	Msg("proxy request")
	reverseProxy.ServeHTTP(c.Writer, c.Request)
}
