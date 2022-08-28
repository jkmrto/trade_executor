package http

import (
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Config is the HTTP config DTO
type Config struct {
	Address string `json:"-"`
}

// EndpointHandlerMethod contains the information regarding the endpoint, method and handler
type EndpointHandlerMethod struct {
	Endpoint    string
	Method      string
	HandlerFunc http.HandlerFunc
}

// NewHandler creates the chi router applying the specified middlewares
func newHandler(handlers []EndpointHandlerMethod) http.Handler {
	r := chi.NewRouter()

	for _, handler := range handlers {
		r.With(middleware.DefaultLogger).MethodFunc(handler.Method, handler.Endpoint, handler.HandlerFunc)
	}

	return r
}

// ListenAndServe creates the server and listens and then serves it.
// TODO: Handle context canceled here to close the server
func ListenAndServe(conf Config, handlers []EndpointHandlerMethod) error {
	handler := newHandler(handlers)
	server := &http.Server{Handler: handler}

	ln, err := net.Listen("tcp", conf.Address)
	if err != nil {
		return err
	}

	err = server.Serve(ln)
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
