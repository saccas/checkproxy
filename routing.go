package main

import (
	r "github.com/unprofession-al/routing"
)

func (s Server) Routes() r.Route {
	return r.Route{
		R: r.Routes{
			"checks/{name}": {
				H: r.Handlers{
					"GET":  {F: s.Report, Q: []*r.QueryParam{}},
					"POST": {F: s.Store, Q: []*r.QueryParam{}},
				},
			},
		},
	}
}
