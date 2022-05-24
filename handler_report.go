package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s Server) Report(res http.ResponseWriter, req *http.Request) {
	if !s.Auth("r", req.Header) {
		s.respond(res, req, http.StatusUnauthorized, "access denied")
		return
	}

	vars := mux.Vars(req)
	checkName, ok := vars["name"]
	if !ok {
		s.respond(res, req, http.StatusNotFound, "check name not provided")
		return
	}

	check, err := s.states.Get(checkName)
	if err != nil {
		s.respond(res, req, http.StatusInternalServerError, fmt.Sprintf("%s", err.Error()))
		return
	}

	if check.TimedOut() {
		s.respond(res, req, http.StatusGatewayTimeout, fmt.Sprintf("The check has not been updated in the last %f seconds", check.ValidityDuration.Seconds()))
		return
	}

	s.raw(res, check.State, check.Body)
}
