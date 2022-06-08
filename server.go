package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/apex/gateway"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/kr/pretty"
	"gopkg.in/yaml.v2"
)

// Server holds all dependencies of the webserver. All related functions such as the
// HTTP handlers dangle off the server struct to allow easy access the dependencies.
type Server struct {
	listener string
	config   *Config
	handler  http.Handler

	states *CheckStates
}

// NewServer creates a populates a Server struct with its dependencies and returns the
// resulting server.
func NewServer(listener string, c *Config) Server {
	s := Server{
		listener: listener,
		config:   c,

		states: NewCheckStates(c.PersistanceBase),
	}

	r := mux.NewRouter().StrictSlash(true)
	routes := s.Routes()
	routes.Populate(r, c.PathPrefix)
	s.handler = alice.New(s.LoggerMiddleware).Then(r)
	return s
}

// run launches the actual web server. If no listener is provided (as flag) the server is
// lauched as a AWL Lambda.
func (s Server) run(mode string) error {
	switch strings.ToLower(mode) {
	case "local":
		if s.listener == "" {
			return fmt.Errorf("No listener defined")
		}
		fmt.Printf("Running locally at '%s'...\n", s.listener)
		return http.ListenAndServe(s.listener, s.handler)
	case "azurefunc":
		port, ok := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
		if !ok {
			return fmt.Errorf("Environment FUNCTIONS_CUSTOMHANDLER_PORT not defined")
		}
		listener := fmt.Sprintf(":%s", port)
		fmt.Printf("Running as Azure Function at '%s'...\n", listener)
		return http.ListenAndServe(listener, s.handler)
	case "awslambda":
		fmt.Printf("Running as AWS Lambda...\n")
		return gateway.ListenAndServe(s.listener, s.handler)
	default:
		return fmt.Errorf("Unknown mode '%s'", mode)
	}
}

// respond allows to easily return some arbitrary data while respecting the `Accept`
// Header to some extend.
func (s Server) respond(res http.ResponseWriter, req *http.Request, code int, data interface{}) {
	var err error
	var errMesg []byte
	var out []byte

	f := req.Header.Get("Accept")
	if f == "text/yaml" {
		res.Header().Set("Content-Type", "text/yaml; charset=utf-8")
		out, err = yaml.Marshal(data)
		errMesg = []byte("--- error: failed while rendering data to yaml")
	} else if f == "text/gostruct" {
		res.Header().Set("Content-Type", "text/gostruct; charset=utf-8")
		out = []byte(pretty.Sprint(data))
	} else {
		res.Header().Set("Content-Type", "application/json; charset=utf-8")
		out, err = json.MarshalIndent(data, "", "    ")
		errMesg = []byte("{ 'error': 'failed while rendering data to json' }")
	}

	if err != nil {
		out = errMesg
		code = http.StatusInternalServerError
	}
	res.WriteHeader(code)
	res.Write(out)
}

// raw is a helper function stat takes a response and returns the plain bytes.
func (s Server) raw(res http.ResponseWriter, code int, data []byte) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(code)
	res.Write(data)
}
