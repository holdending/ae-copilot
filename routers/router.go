package routers

import (
	"bytes"
	"io"
	"net/http"

	"github.com/LiveRamp/ae-copilot/controllers"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var heartbeatController = &controllers.HeartbeatController{}

var routes = []Route{
	{"HeartbeatGet", http.MethodGet, "/heartbeat", heartbeatController.Get},
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != "/v1/heartbeat" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				logs.Error("Error reading body: %v", err)
				http.Error(w, "can't read body", http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(body))
			logs.Debug("Request URI: %v %v, request body: %v", r.Method, r.RequestURI, string(body))
		}
		next.ServeHTTP(w, r)
	})
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		logs.Debug(route.Name)
		router.
			Methods(route.Method).
			PathPrefix("/v1").
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	router.Use(loggingMiddleware)

	return router
}
