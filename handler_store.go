package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

func (s Server) Store(res http.ResponseWriter, req *http.Request) {
	// Check auth header
	if !s.Auth(authWrite, req.Header) {
		s.respond(res, req, http.StatusUnauthorized, "access denied")
		return
	}

	// Get URL Path
	vars := mux.Vars(req)
	checkName, ok := vars["name"]
	if !ok {
		s.respond(res, req, http.StatusNotFound, fmt.Sprintf("check name not provided"))
		return
	}

	// Get Request Body
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		s.respond(res, req, http.StatusInternalServerError, fmt.Sprintf("could not read request body: %s", err.Error()))
		return
	}

	// Get Params
	params := &StoreRequestQueryParams{}
	err = params.Parse(req)
	if err != nil {
		s.respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	err = params.Validate()
	if err != nil {
		s.respond(res, req, http.StatusInternalServerError, err.Error())
		return
	}

	// Persist
	err = s.states.Set(checkName, params.Status, params.ValidityDuration, body)
	if err != nil {
		s.respond(res, req, http.StatusInternalServerError, fmt.Sprintf("check '%s' could not be saved: %s", checkName, err.Error()))
		return
	}

	s.respond(res, req, http.StatusOK, "saved")
}

type StoreRequestQueryParams struct {
	Status           int    `schema:"status,required"`
	ValidityDuration string `schema:"validity_duration"`
}

func (r *StoreRequestQueryParams) Default() {
	r.Status = 501
	r.ValidityDuration = "10m"
}

func (r *StoreRequestQueryParams) Parse(req *http.Request) error {
	r.Default()
	return schema.NewDecoder().Decode(r, req.URL.Query())
}

func (r *StoreRequestQueryParams) Validate() error {
	return nil
}
